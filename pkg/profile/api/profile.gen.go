// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/discord-gophers/goapi-gen version v0.3.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// DeliveryOption defines model for DeliveryOption.
type DeliveryOption struct {
	DisplayValue string `json:"display_value,omitempty"`
	HashedValue  string `json:"hashed_value,omitempty"`
	UserID       string `json:"user_id,omitempty"`
}

// Error defines model for Error.
type Error struct {
	// Error code
	Code string `json:"code"`

	// Error message
	Message string `json:"message"`
}

// Login defines model for Login.
type Login struct {
	// Token for 2FA verification if required
	LoginToken *string `json:"loginToken,omitempty"`
	Message    string  `json:"message"`

	// Whether 2FA verification is required
	Requires2fA *bool  `json:"requires2FA,omitempty"`
	Status      string `json:"status"`
	User        User   `json:"user"`

	// List of users associated with the login. Usually contains one user, but may contain multiple if same username is shared.
	Users []User `json:"users,omitempty"`
}

// LoginOption defines model for LoginOption.
type LoginOption struct {
	// Whether this is the current login
	Current bool `json:"current,omitempty"`

	// ID of the login
	ID string `json:"id,omitempty"`

	// Username of the login
	Username string `json:"username,omitempty"`
}

// LoginSelectionRequiredResponse defines model for LoginSelectionRequiredResponse.
type LoginSelectionRequiredResponse struct {
	LoginOptions []LoginOption `json:"login_options"`
	Message      string        `json:"message"`
	Status       string        `json:"status"`
}

// MultiUsersResponse defines model for MultiUsersResponse.
type MultiUsersResponse struct {
	Users []User `json:"users,omitempty"`
}

// PasswordPolicyResponse defines model for PasswordPolicyResponse.
type PasswordPolicyResponse struct {
	// Whether common passwords are disallowed
	DisallowCommonPwds *bool `json:"disallow_common_pwds,omitempty"`

	// Number of days until password expires
	ExpirationDays *int `json:"expiration_days,omitempty"`

	// Number of previous passwords to check against
	HistoryCheckCount *int `json:"history_check_count,omitempty"`

	// Maximum number of repeated characters allowed
	MaxRepeatedChars *int `json:"max_repeated_chars,omitempty"`

	// Minimum length of the password
	MinLength *int `json:"min_length,omitempty"`

	// Whether the password requires a digit
	RequireDigit *bool `json:"require_digit,omitempty"`

	// Whether the password requires a lowercase letter
	RequireLowercase *bool `json:"require_lowercase,omitempty"`

	// Whether the password requires a special character
	RequireSpecialChar *bool `json:"require_special_char,omitempty"`

	// Whether the password requires an uppercase letter
	RequireUppercase *bool `json:"require_uppercase,omitempty"`
}

// Structure added for integration compatibility purposes
type SingleUserResponse struct {
	User User `json:"user,omitempty"`
}

// SuccessResponse defines model for SuccessResponse.
type SuccessResponse struct {
	Result string `json:"result,omitempty"`
}

// TwoFactorMethod defines model for TwoFactorMethod.
type TwoFactorMethod struct {
	Enabled     bool   `json:"enabled"`
	TwoFactorID string `json:"two_factor_id,omitempty"`
	Type        string `json:"type"`
}

// TwoFactorMethodSelection defines model for TwoFactorMethodSelection.
type TwoFactorMethodSelection struct {
	DeliveryOptions []DeliveryOption `json:"delivery_options,omitempty"`
	Type            string           `json:"type,omitempty"`
}

// TwoFactorMethods defines model for TwoFactorMethods.
type TwoFactorMethods struct {
	Count   int               `json:"count"`
	Methods []TwoFactorMethod `json:"methods"`
}

// TwoFactorRequiredResponse defines model for TwoFactorRequiredResponse.
type TwoFactorRequiredResponse struct {
	Message string `json:"message,omitempty"`
	Status  string `json:"status,omitempty"`

	// Temporary token to use for 2FA verification
	TempToken        string                     `json:"temp_token,omitempty"`
	TwoFactorMethods []TwoFactorMethodSelection `json:"two_factor_methods,omitempty"`
}

// User defines model for User.
type User struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

// Delete2faJSONBody defines parameters for Delete2fa.
type Delete2faJSONBody struct {
	TwofaID   *string                    `json:"twofa_id,omitempty"`
	TwofaType Delete2faJSONBodyTwofaType `json:"twofa_type"`
}

// Delete2faJSONBodyTwofaType defines parameters for Delete2fa.
type Delete2faJSONBodyTwofaType string

// Post2faDisableJSONBody defines parameters for Post2faDisable.
type Post2faDisableJSONBody struct {
	TwofaType Post2faDisableJSONBodyTwofaType `json:"twofa_type"`
}

// Post2faDisableJSONBodyTwofaType defines parameters for Post2faDisable.
type Post2faDisableJSONBodyTwofaType string

// Post2faEnableJSONBody defines parameters for Post2faEnable.
type Post2faEnableJSONBody struct {
	TwofaType Post2faEnableJSONBodyTwofaType `json:"twofa_type"`
}

// Post2faEnableJSONBodyTwofaType defines parameters for Post2faEnable.
type Post2faEnableJSONBodyTwofaType string

// Post2faSetupJSONBody defines parameters for Post2faSetup.
type Post2faSetupJSONBody struct {
	TwofaType Post2faSetupJSONBodyTwofaType `json:"twofa_type"`
}

// Post2faSetupJSONBodyTwofaType defines parameters for Post2faSetup.
type Post2faSetupJSONBodyTwofaType string

// AssociateLoginJSONBody defines parameters for AssociateLogin.
type AssociateLoginJSONBody struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// CompleteLoginAssociationJSONBody defines parameters for CompleteLoginAssociation.
type CompleteLoginAssociationJSONBody struct {
	// ID of the login the user selected
	LoginID string `json:"login_id"`
}

// ChangePasswordJSONBody defines parameters for ChangePassword.
type ChangePasswordJSONBody struct {
	// User's current password
	CurrentPassword string `json:"current_password"`

	// User's new password
	NewPassword string `json:"new_password"`
}

// PostUserSwitchJSONBody defines parameters for PostUserSwitch.
type PostUserSwitchJSONBody struct {
	// ID of the user to switch to
	UserID string `json:"user_id"`
}

// ChangeUsernameJSONBody defines parameters for ChangeUsername.
type ChangeUsernameJSONBody struct {
	// User's current password for verification
	CurrentPassword string `json:"currentPassword"`

	// New username to set
	NewUsername string `json:"newUsername"`
}

// Delete2faJSONRequestBody defines body for Delete2fa for application/json ContentType.
type Delete2faJSONRequestBody Delete2faJSONBody

// Bind implements render.Binder.
func (Delete2faJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// Post2faDisableJSONRequestBody defines body for Post2faDisable for application/json ContentType.
type Post2faDisableJSONRequestBody Post2faDisableJSONBody

// Bind implements render.Binder.
func (Post2faDisableJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// Post2faEnableJSONRequestBody defines body for Post2faEnable for application/json ContentType.
type Post2faEnableJSONRequestBody Post2faEnableJSONBody

// Bind implements render.Binder.
func (Post2faEnableJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// Post2faSetupJSONRequestBody defines body for Post2faSetup for application/json ContentType.
type Post2faSetupJSONRequestBody Post2faSetupJSONBody

// Bind implements render.Binder.
func (Post2faSetupJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// AssociateLoginJSONRequestBody defines body for AssociateLogin for application/json ContentType.
type AssociateLoginJSONRequestBody AssociateLoginJSONBody

// Bind implements render.Binder.
func (AssociateLoginJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// CompleteLoginAssociationJSONRequestBody defines body for CompleteLoginAssociation for application/json ContentType.
type CompleteLoginAssociationJSONRequestBody CompleteLoginAssociationJSONBody

// Bind implements render.Binder.
func (CompleteLoginAssociationJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// ChangePasswordJSONRequestBody defines body for ChangePassword for application/json ContentType.
type ChangePasswordJSONRequestBody ChangePasswordJSONBody

// Bind implements render.Binder.
func (ChangePasswordJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PostUserSwitchJSONRequestBody defines body for PostUserSwitch for application/json ContentType.
type PostUserSwitchJSONRequestBody PostUserSwitchJSONBody

// Bind implements render.Binder.
func (PostUserSwitchJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// ChangeUsernameJSONRequestBody defines body for ChangeUsername for application/json ContentType.
type ChangeUsernameJSONRequestBody ChangeUsernameJSONBody

// Bind implements render.Binder.
func (ChangeUsernameJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// Response is a common response struct for all the API calls.
// A Response object may be instantiated via functions for specific operation responses.
// It may also be instantiated directly, for the purpose of responding with a single status code.
type Response struct {
	body        interface{}
	Code        int
	contentType string
}

// Render implements the render.Renderer interface. It sets the Content-Type header
// and status code based on the response definition.
func (resp *Response) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", resp.contentType)
	render.Status(r, resp.Code)
	return nil
}

// Status is a builder method to override the default status code for a response.
func (resp *Response) Status(code int) *Response {
	resp.Code = code
	return resp
}

// ContentType is a builder method to override the default content type for a response.
func (resp *Response) ContentType(contentType string) *Response {
	resp.contentType = contentType
	return resp
}

// MarshalJSON implements the json.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.body)
}

// MarshalXML implements the xml.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.Encode(resp.body)
}

// Get2faMethodsJSON200Response is a constructor method for a Get2faMethods response.
// A *Response is returned with the configured status code and content type from the spec.
func Get2faMethodsJSON200Response(body TwoFactorMethods) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// Get2faMethodsJSON404Response is a constructor method for a Get2faMethods response.
// A *Response is returned with the configured status code and content type from the spec.
func Get2faMethodsJSON404Response(body struct {
	Message *string `json:"message,omitempty"`
}) *Response {
	return &Response{
		body:        body,
		Code:        404,
		contentType: "application/json",
	}
}

// Delete2faJSON200Response is a constructor method for a Delete2fa response.
// A *Response is returned with the configured status code and content type from the spec.
func Delete2faJSON200Response(body SuccessResponse) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// Post2faDisableJSON200Response is a constructor method for a Post2faDisable response.
// A *Response is returned with the configured status code and content type from the spec.
func Post2faDisableJSON200Response(body SuccessResponse) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// Post2faEnableJSON200Response is a constructor method for a Post2faEnable response.
// A *Response is returned with the configured status code and content type from the spec.
func Post2faEnableJSON200Response(body SuccessResponse) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// Post2faSetupJSON201Response is a constructor method for a Post2faSetup response.
// A *Response is returned with the configured status code and content type from the spec.
func Post2faSetupJSON201Response(body SuccessResponse) *Response {
	return &Response{
		body:        body,
		Code:        201,
		contentType: "application/json",
	}
}

// AssociateLoginJSON200Response is a constructor method for a AssociateLogin response.
// A *Response is returned with the configured status code and content type from the spec.
func AssociateLoginJSON200Response(body SuccessResponse) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// AssociateLoginJSON202Response is a constructor method for a AssociateLogin response.
// A *Response is returned with the configured status code and content type from the spec.
func AssociateLoginJSON202Response(body interface{}) *Response {
	return &Response{
		body:        body,
		Code:        202,
		contentType: "application/json",
	}
}

// CompleteLoginAssociationJSON200Response is a constructor method for a CompleteLoginAssociation response.
// A *Response is returned with the configured status code and content type from the spec.
func CompleteLoginAssociationJSON200Response(body SuccessResponse) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// ChangePasswordJSON400Response is a constructor method for a ChangePassword response.
// A *Response is returned with the configured status code and content type from the spec.
func ChangePasswordJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// ChangePasswordJSON401Response is a constructor method for a ChangePassword response.
// A *Response is returned with the configured status code and content type from the spec.
func ChangePasswordJSON401Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        401,
		contentType: "application/json",
	}
}

// ChangePasswordJSON403Response is a constructor method for a ChangePassword response.
// A *Response is returned with the configured status code and content type from the spec.
func ChangePasswordJSON403Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        403,
		contentType: "application/json",
	}
}

// ChangePasswordJSON500Response is a constructor method for a ChangePassword response.
// A *Response is returned with the configured status code and content type from the spec.
func ChangePasswordJSON500Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        500,
		contentType: "application/json",
	}
}

// GetPasswordPolicyJSON200Response is a constructor method for a GetPasswordPolicy response.
// A *Response is returned with the configured status code and content type from the spec.
func GetPasswordPolicyJSON200Response(body PasswordPolicyResponse) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// PostUserSwitchJSON200Response is a constructor method for a PostUserSwitch response.
// A *Response is returned with the configured status code and content type from the spec.
func PostUserSwitchJSON200Response(body interface{}) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// PostUserSwitchJSON400Response is a constructor method for a PostUserSwitch response.
// A *Response is returned with the configured status code and content type from the spec.
func PostUserSwitchJSON400Response(body struct {
	Message *string `json:"message,omitempty"`
}) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// PostUserSwitchJSON403Response is a constructor method for a PostUserSwitch response.
// A *Response is returned with the configured status code and content type from the spec.
func PostUserSwitchJSON403Response(body struct {
	Message *string `json:"message,omitempty"`
}) *Response {
	return &Response{
		body:        body,
		Code:        403,
		contentType: "application/json",
	}
}

// ChangeUsernameJSON400Response is a constructor method for a ChangeUsername response.
// A *Response is returned with the configured status code and content type from the spec.
func ChangeUsernameJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// ChangeUsernameJSON401Response is a constructor method for a ChangeUsername response.
// A *Response is returned with the configured status code and content type from the spec.
func ChangeUsernameJSON401Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        401,
		contentType: "application/json",
	}
}

// ChangeUsernameJSON403Response is a constructor method for a ChangeUsername response.
// A *Response is returned with the configured status code and content type from the spec.
func ChangeUsernameJSON403Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        403,
		contentType: "application/json",
	}
}

// ChangeUsernameJSON409Response is a constructor method for a ChangeUsername response.
// A *Response is returned with the configured status code and content type from the spec.
func ChangeUsernameJSON409Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        409,
		contentType: "application/json",
	}
}

// ChangeUsernameJSON500Response is a constructor method for a ChangeUsername response.
// A *Response is returned with the configured status code and content type from the spec.
func ChangeUsernameJSON500Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        500,
		contentType: "application/json",
	}
}

// FindUsersWithLoginJSON200Response is a constructor method for a FindUsersWithLogin response.
// A *Response is returned with the configured status code and content type from the spec.
func FindUsersWithLoginJSON200Response(body interface{}) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get login 2FA methods
	// (GET /2fa)
	Get2faMethods(w http.ResponseWriter, r *http.Request) *Response
	// Delete a 2FA method
	// (POST /2fa/delete)
	Delete2fa(w http.ResponseWriter, r *http.Request) *Response
	// Disable an existing 2FA method
	// (POST /2fa/disable)
	Post2faDisable(w http.ResponseWriter, r *http.Request) *Response
	// Enable an existing 2FA method
	// (POST /2fa/enable)
	Post2faEnable(w http.ResponseWriter, r *http.Request) *Response
	// Create a new 2FA method
	// (POST /2fa/setup)
	Post2faSetup(w http.ResponseWriter, r *http.Request) *Response
	// Associate a login
	// (POST /login/associate)
	AssociateLogin(w http.ResponseWriter, r *http.Request) *Response
	// Complete login association after user selection
	// (POST /login/associate/complete)
	CompleteLoginAssociation(w http.ResponseWriter, r *http.Request) *Response
	// Change user password
	// (PUT /password)
	ChangePassword(w http.ResponseWriter, r *http.Request) *Response
	// Get password policy
	// (GET /password/policy)
	GetPasswordPolicy(w http.ResponseWriter, r *http.Request) *Response
	// Switch to a different user when multiple users are available for the same login
	// (POST /user/switch)
	PostUserSwitch(w http.ResponseWriter, r *http.Request) *Response
	// Change username
	// (PUT /username)
	ChangeUsername(w http.ResponseWriter, r *http.Request) *Response
	// Get a list of users associated with the current login
	// (GET /users)
	FindUsersWithLogin(w http.ResponseWriter, r *http.Request) *Response
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler          ServerInterface
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// Get2faMethods operation middleware
func (siw *ServerInterfaceWrapper) Get2faMethods(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.Get2faMethods(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// Delete2fa operation middleware
func (siw *ServerInterfaceWrapper) Delete2fa(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.Delete2fa(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// Post2faDisable operation middleware
func (siw *ServerInterfaceWrapper) Post2faDisable(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.Post2faDisable(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// Post2faEnable operation middleware
func (siw *ServerInterfaceWrapper) Post2faEnable(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.Post2faEnable(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// Post2faSetup operation middleware
func (siw *ServerInterfaceWrapper) Post2faSetup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.Post2faSetup(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// AssociateLogin operation middleware
func (siw *ServerInterfaceWrapper) AssociateLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.AssociateLogin(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// CompleteLoginAssociation operation middleware
func (siw *ServerInterfaceWrapper) CompleteLoginAssociation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.CompleteLoginAssociation(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// ChangePassword operation middleware
func (siw *ServerInterfaceWrapper) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.ChangePassword(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// GetPasswordPolicy operation middleware
func (siw *ServerInterfaceWrapper) GetPasswordPolicy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetPasswordPolicy(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostUserSwitch operation middleware
func (siw *ServerInterfaceWrapper) PostUserSwitch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostUserSwitch(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// ChangeUsername operation middleware
func (siw *ServerInterfaceWrapper) ChangeUsername(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.ChangeUsername(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// FindUsersWithLogin operation middleware
func (siw *ServerInterfaceWrapper) FindUsersWithLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.FindUsersWithLogin(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter %s: %v", err.paramName, err.err)
}

func (err UnescapedCookieParamError) Unwrap() error { return err.err }

type UnmarshalingParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err UnmarshalingParamError) Error() string {
	return fmt.Sprintf("error unmarshaling parameter %s as JSON: %v", err.paramName, err.err)
}

func (err UnmarshalingParamError) Unwrap() error { return err.err }

type RequiredParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err RequiredParamError) Error() string {
	if err.err == nil {
		return fmt.Sprintf("query parameter %s is required, but not found", err.paramName)
	} else {
		return fmt.Sprintf("query parameter %s is required, but errored: %s", err.paramName, err.err)
	}
}

func (err RequiredParamError) Unwrap() error { return err.err }

type RequiredHeaderError struct {
	paramName string
}

// Error implements error.
func (err RequiredHeaderError) Error() string {
	return fmt.Sprintf("header parameter %s is required, but not found", err.paramName)
}

type InvalidParamFormatError struct {
	err       error
	paramName string
}

// Error implements error.
func (err InvalidParamFormatError) Error() string {
	return fmt.Sprintf("invalid format for parameter %s: %v", err.paramName, err.err)
}

func (err InvalidParamFormatError) Unwrap() error { return err.err }

type TooManyValuesForParamError struct {
	NumValues int
	paramName string
}

// Error implements error.
func (err TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("expected one value for %s, got %d", err.paramName, err.NumValues)
}

// ParameterName is an interface that is implemented by error types that are
// relevant to a specific parameter.
type ParameterError interface {
	error
	// ParamName is the name of the parameter that the error is referring to.
	ParamName() string
}

func (err UnescapedCookieParamError) ParamName() string  { return err.paramName }
func (err UnmarshalingParamError) ParamName() string     { return err.paramName }
func (err RequiredParamError) ParamName() string         { return err.paramName }
func (err RequiredHeaderError) ParamName() string        { return err.paramName }
func (err InvalidParamFormatError) ParamName() string    { return err.paramName }
func (err TooManyValuesForParamError) ParamName() string { return err.paramName }

type ServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

type ServerOption func(*ServerOptions)

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface, opts ...ServerOption) http.Handler {
	options := &ServerOptions{
		BaseURL:    "/",
		BaseRouter: chi.NewRouter(),
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
	}

	for _, f := range opts {
		f(options)
	}

	r := options.BaseRouter
	wrapper := ServerInterfaceWrapper{
		Handler:          si,
		ErrorHandlerFunc: options.ErrorHandlerFunc,
	}

	r.Route(options.BaseURL, func(r chi.Router) {
		r.Get("/2fa", wrapper.Get2faMethods)
		r.Post("/2fa/delete", wrapper.Delete2fa)
		r.Post("/2fa/disable", wrapper.Post2faDisable)
		r.Post("/2fa/enable", wrapper.Post2faEnable)
		r.Post("/2fa/setup", wrapper.Post2faSetup)
		r.Post("/login/associate", wrapper.AssociateLogin)
		r.Post("/login/associate/complete", wrapper.CompleteLoginAssociation)
		r.Put("/password", wrapper.ChangePassword)
		r.Get("/password/policy", wrapper.GetPasswordPolicy)
		r.Post("/user/switch", wrapper.PostUserSwitch)
		r.Put("/username", wrapper.ChangeUsername)
		r.Get("/users", wrapper.FindUsersWithLogin)
	})
	return r
}

func WithRouter(r chi.Router) ServerOption {
	return func(s *ServerOptions) {
		s.BaseRouter = r
	}
}

func WithServerBaseURL(url string) ServerOption {
	return func(s *ServerOptions) {
		s.BaseURL = url
	}
}

func WithErrorHandler(handler func(w http.ResponseWriter, r *http.Request, err error)) ServerOption {
	return func(s *ServerOptions) {
		s.ErrorHandlerFunc = handler
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"Wq3+EwAA///XRnjBNDEAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
