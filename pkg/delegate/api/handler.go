package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tendant/simple-idm/pkg/client"
	delegate "github.com/tendant/simple-idm/pkg/delegate"
	"github.com/tendant/simple-idm/pkg/login/api"
	"github.com/tendant/simple-idm/pkg/mapper"
	tg "github.com/tendant/simple-idm/pkg/tokengenerator"
)

// Constants for token cookie names
const (
	ACCESS_TOKEN_NAME  = api.ACCESS_TOKEN_NAME
	REFRESH_TOKEN_NAME = api.REFRESH_TOKEN_NAME
)

// Handler implements the ServerInterface for delegation API
type Handle struct {
	service            *delegate.Service
	tokenService       tg.TokenService
	tokenCookieService tg.TokenCookieService
}

// NewHandler creates a new delegation API handler
func NewHandler(service *delegate.Service, tokenService tg.TokenService, tokenCookieService tg.TokenCookieService) *Handle {
	return &Handle{
		service:            service,
		tokenService:       tokenService,
		tokenCookieService: tokenCookieService,
	}
}

// CreateDelegate handles the POST /delegate endpoint
// It creates an impersonation session allowing a delegatee to access a delegator's account
func (h *Handle) CreateDelegate(w http.ResponseWriter, r *http.Request) *Response {
	// Get the current user from context (this would be set by your auth middleware)
	authUser, ok := r.Context().Value(client.AuthUserKey).(*client.AuthUser)
	if !ok {
		slog.Error("Failed to get authenticated user from context")
		return CreateDelegateJSON401Response(ErrorResponse{
			Error: "Unauthorized",
			Code:  stringPtr("unauthorized"),
		})
	}

	// Get the login UUID from authUser (it's already a uuid.UUID type)
	delegateeUuidStr := authUser.UserId
	delegateeUuid, err := uuid.Parse(delegateeUuidStr)
	if err != nil {
		slog.Error("Failed to parse delegatee UUID", "error", err)
		return CreateDelegateJSON400Response(ErrorResponse{
			Error: "Invalid delegatee UUID",
			Code:  stringPtr("invalid_uuid"),
		})
	}

	// Parse request body
	var reqBody CreateDelegateJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		return CreateDelegateJSON400Response(ErrorResponse{
			Error: "Invalid request body",
			Code:  stringPtr("invalid_request"),
		})
	}

	// Validate delegator_user_uuid
	delegatorUserUUID, err := uuid.Parse(reqBody.DelegatorUserUUID)
	if err != nil {
		slog.Error("Invalid delegator_user_uuid", "error", err)
		return CreateDelegateJSON400Response(ErrorResponse{
			Error: "Invalid delegator user UUID",
			Code:  stringPtr("invalid_uuid"),
		})
	}

	// Get delegated users for the current user
	delegators, err := h.service.FindDelegators(r.Context(), delegateeUuid)
	if err != nil {
		slog.Error("Failed to get delegators", "error", err, "delegatee_user_id", delegateeUuid)
		return CreateDelegateJSON403Response(ErrorResponse{
			Error: "Failed to get delegated users",
			Code:  stringPtr("server_error"),
		})
	}
	slog.Info("Delegators retrieved", "delegators", delegators)

	// Check if the requested delegator is in the list of delegated users
	var foundDelegator bool
	var selectedUser mapper.User

	for _, user := range delegators {
		if user.UserId == delegatorUserUUID.String() {
			foundDelegator = true
			selectedUser = user
			slog.Info("Delegator found", "delegator_uuid", delegatorUserUUID)
			break
		}
	}

	if !foundDelegator {
		slog.Error("User not authorized to delegate the requested delegator", "delegator_uuid", delegatorUserUUID)
		return CreateDelegateJSON403Response(ErrorResponse{
			Error: "Not authorized to delegate this user",
			Code:  stringPtr("forbidden"),
		})
	}

	// generate extra_claims and add delegate_user_id into extra_claims
	_, extraClaims := h.service.ToTokenClaims(selectedUser)
	extraClaims["delegate_user_id"] = authUser.UserId

	// Generate tokens using the token service
	tokens, err := h.tokenService.GenerateTokens(selectedUser.UserId, nil, extraClaims)
	if err != nil {
		slog.Error("Failed to generate tokens", "error", err)
		return &Response{
			body: "Failed to generate tokens",
			Code: http.StatusInternalServerError,
		}
	}

	// Set cookies for the tokens
	err = h.tokenCookieService.SetTokensCookie(w, tokens)
	if err != nil {
		slog.Error("Failed to set tokens in cookies", "error", err)
		return &Response{
			body: "Failed to set tokens in cookies",
			Code: http.StatusInternalServerError,
		}
	}

	slog.Info("Impersonation successful", "delegator_uuid", delegatorUserUUID)

	// Return the success response
	resp := SuccessResponse{}
	resp["message"] = "success"
	return CreateDelegateJSON200Response(resp)
}

// CreateDelegateBack handles the POST /delegate/back endpoint
// It ends the current delegation session and returns to the original user context
func (h *Handle) CreateDelegateBack(w http.ResponseWriter, r *http.Request) *Response {
	// Get the current user from context (this would be set by your auth middleware)
	authUser, ok := r.Context().Value(client.AuthUserKey).(*client.AuthUser)
	if !ok {
		slog.Error("Failed to get authenticated user from context")
		return CreateDelegateJSON401Response(ErrorResponse{
			Error: "Unauthorized",
			Code:  stringPtr("unauthorized"),
		})
	}

	// Check if the current user is in an impersonation session
	// Get the access token from the cookie
	accessTokenCookie, err := r.Cookie(ACCESS_TOKEN_NAME)
	if err != nil {
		slog.Error("Failed to get access token cookie", "error", err)
		return CreateDelegateJSON401Response(ErrorResponse{
			Error: "No access token found",
			Code:  stringPtr("unauthorized"),
		})
	}

	// Parse the token to verify impersonation
	token, err := h.tokenService.ParseToken(accessTokenCookie.Value)
	if err != nil {
		slog.Error("Failed to parse access token", "error", err)
		return CreateDelegateJSON401Response(ErrorResponse{
			Error: "Invalid access token",
			Code:  stringPtr("unauthorized"),
		})
	}

	// Check if the token contains impersonation claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		slog.Error("Failed to get claims from token")
		return CreateDelegateJSON400Response(ErrorResponse{
			Error: "Not in an impersonation session",
			Code:  stringPtr("not_impersonating"),
		})
	}

	// Extract extra claims from the token
	extraClaimsRaw, ok := claims["extra_claims"]
	if !ok {
		slog.Error("No extra claims found in token")
		return CreateDelegateJSON400Response(ErrorResponse{
			Error: "Not in an impersonation session",
			Code:  stringPtr("not_impersonating"),
		})
	}

	// Convert extra claims to map
	extraClaimsMap, ok := extraClaimsRaw.(map[string]interface{})
	if !ok {
		slog.Error("Extra claims is not a map")
		return CreateDelegateJSON400Response(ErrorResponse{
			Error: "Invalid token format",
			Code:  stringPtr("invalid_token"),
		})
	}

	// Check for original user information in the claims
	delegateUserID, ok := extraClaimsMap["delegate_user_id"].(string)
	if !ok {
		slog.Error("No delegate user ID found in token claims")
		return CreateDelegateJSON400Response(ErrorResponse{
			Error: "Not in an impersonation session",
			Code:  stringPtr("not_impersonating"),
		})
	}
	slog.Info("Delegate user ID found in token claims", "delegate_user_id", delegateUserID)

	if delegateUserID == authUser.UserId {
		slog.Error("Delegate user ID matches current user ID")
		return CreateDelegateJSON400Response(ErrorResponse{
			Error: "Not in an impersonation session",
			Code:  stringPtr("not_impersonating"),
		})
	}

	slog.Info("Impersonation session found", "delegate_user_id", delegateUserID, "user_id", authUser.UserId)

	// Parse the original user ID
	delegateUserUUID, err := uuid.Parse(delegateUserID)
	if err != nil {
		slog.Error("Invalid delegate user ID in token claims", "error", err)
		return CreateDelegateJSON400Response(ErrorResponse{
			Error: "Invalid impersonation data",
			Code:  stringPtr("invalid_impersonation"),
		})
	}

	// Create a mappedUser for token generation
	originalUser, err := h.service.GetOriginalUser(r.Context(), delegateUserUUID)
	if err != nil {
		slog.Error("Failed to get original user", "error", err)
		return &Response{
			body: "Failed to get original user",
			Code: http.StatusInternalServerError,
		}
	}

	slog.Info("Original user retrieved", "user_id", originalUser.UserId)

	_, extraClaims := h.service.ToTokenClaims(originalUser)
	// Generate tokens using the token service
	tokens, err := h.tokenService.GenerateTokens(originalUser.UserId, nil, extraClaims)
	if err != nil {
		slog.Error("Failed to generate tokens", "error", err)
		return &Response{
			body: "Failed to generate tokens",
			Code: http.StatusInternalServerError,
		}
	}

	// Set cookies for the tokens
	err = h.tokenCookieService.SetTokensCookie(w, tokens)
	if err != nil {
		slog.Error("Failed to set tokens in cookies", "error", err)
		return &Response{
			body: "Failed to set tokens in cookies",
			Code: http.StatusInternalServerError,
		}
	}
	slog.Info("Delegation back succeed")

	// Return the success response
	resp := SuccessResponse{}
	resp["message"] = "success"
	return CreateDelegateJSON200Response(resp)
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}
