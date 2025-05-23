package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/tendant/simple-idm/pkg/client"
	"github.com/tendant/simple-idm/pkg/device"
)

// DeviceHandler handles HTTP requests for device management
type DeviceHandler struct {
	deviceService *device.DeviceService
}

// NewDeviceHandler creates a new device handler
func NewDeviceHandler(deviceService *device.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		deviceService: deviceService,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// DeviceWithLogin represents a device with its linked login information
type DeviceWithLogin struct {
	device.Device
	LinkedLogins []LoginInfo `json:"linked_logins,omitempty"`
	ExpiresAt    string      `json:"expires_at,omitempty"` // When the device-login link expires
}

// LoginInfo represents basic login information
type LoginInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// ListDevicesResponse represents the response body for listing devices
type ListDevicesResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Devices []DeviceWithLogin `json:"devices"`
}

// GetDevicesByLogin handles fetching devices linked to a specific login
func (h *DeviceHandler) GetDevicesByLogin(w http.ResponseWriter, r *http.Request) {
	// Get login ID from URL parameter
	loginIDStr := chi.URLParam(r, "login_id")
	if loginIDStr == "" {
		renderErrorResponse(w, r, http.StatusBadRequest, "Missing required parameter", "login_id is required")
		return
	}

	// Parse login ID
	loginID, err := uuid.Parse(loginIDStr)
	if err != nil {
		slog.Error("Failed to parse login ID", "error", err)
		renderErrorResponse(w, r, http.StatusBadRequest, "Invalid login ID", err.Error())
		return
	}

	// Get authenticated user from context
	authUser, ok := r.Context().Value(client.AuthUserKey).(*client.AuthUser)
	if !ok || authUser == nil {
		renderErrorResponse(w, r, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Check if user has permission to view devices for this login
	// Either the user is an admin or they're viewing their own login's devices
	if !client.IsAdmin(authUser) && authUser.LoginID != loginID {
		renderErrorResponse(w, r, http.StatusForbidden, "Permission denied", "You don't have permission to view devices for this login")
		return
	}

	// Get devices for the login
	devices, err := h.deviceService.FindDevicesByLogin(r.Context(), loginID)
	if err != nil {
		slog.Error("Failed to get devices for login", "error", err)
		renderErrorResponse(w, r, http.StatusInternalServerError, "Failed to get devices for login", err.Error())
		return
	}

	// Convert devices to DeviceWithLogin
	devicesWithLogin := make([]DeviceWithLogin, 0, len(devices))
	for _, d := range devices {
		// Get the login device link to get expiration information
		loginDevice, err := h.deviceService.FindLoginDeviceByFingerprintAndLoginID(r.Context(), d.Fingerprint, loginID)
		if err != nil {
			slog.Error("Failed to get login device link", "fingerprint", d.Fingerprint, "loginID", loginID, "error", err)
			// Continue with other devices even if we can't get link info for this one
			deviceWithLogin := DeviceWithLogin{
				Device: d,
				LinkedLogins: []LoginInfo{
					{
						ID:       loginID.String(),
						Username: "N/A", // We don't have the username here, but we know it's linked to this login
					},
				},
			}
			devicesWithLogin = append(devicesWithLogin, deviceWithLogin)
			continue
		}

		deviceWithLogin := DeviceWithLogin{
			Device: d,
			LinkedLogins: []LoginInfo{
				{
					ID:       loginID.String(),
					Username: "N/A", // We don't have the username here, but we know it's linked to this login
				},
			},
			ExpiresAt: loginDevice.ExpiresAt.Format(http.TimeFormat),
		}
		devicesWithLogin = append(devicesWithLogin, deviceWithLogin)
	}

	// Return success response
	response := ListDevicesResponse{
		Status:  "success",
		Message: "Devices retrieved successfully",
		Devices: devicesWithLogin,
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// UnlinkDeviceFromLoginRequest represents the request body for unlinking a device from a login
type UnlinkDeviceFromLoginRequest struct {
	Fingerprint string `json:"fingerprint"`
}

// UnlinkDeviceFromLoginResponse represents the response body for unlinking a device from a login
type UnlinkDeviceFromLoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// UnlinkDeviceFromLogin handles unlinking a device from a login
func (h *DeviceHandler) UnlinkDeviceFromLogin(w http.ResponseWriter, r *http.Request) {
	// Get login ID from URL parameter
	loginIDStr := chi.URLParam(r, "login_id")
	if loginIDStr == "" {
		renderErrorResponse(w, r, http.StatusBadRequest, "Missing required parameter", "login_id is required")
		return
	}

	// Parse login ID
	loginID, err := uuid.Parse(loginIDStr)
	if err != nil {
		slog.Error("Failed to parse login ID", "error", err)
		renderErrorResponse(w, r, http.StatusBadRequest, "Invalid login ID", err.Error())
		return
	}

	// Get authenticated user from context
	authUser, ok := r.Context().Value(client.AuthUserKey).(*client.AuthUser)
	if !ok || authUser == nil {
		renderErrorResponse(w, r, http.StatusUnauthorized, "Authentication required", "")
		return
	}

	// Check if user has permission to unlink devices for this login
	// Either the user is an admin or they're unlinking from their own login
	if !client.IsAdmin(authUser) && authUser.LoginID != loginID {
		renderErrorResponse(w, r, http.StatusForbidden, "Permission denied", "You don't have permission to unlink devices for this login")
		return
	}

	// Parse request body
	var req UnlinkDeviceFromLoginRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		renderErrorResponse(w, r, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if req.Fingerprint == "" {
		renderErrorResponse(w, r, http.StatusBadRequest, "Missing required field", "fingerprint is required")
		return
	}

	// Unlink the device from the login
	err = h.deviceService.UnlinkLoginFromDevice(r.Context(), loginID, req.Fingerprint)
	if err != nil {
		slog.Error("Failed to unlink device from login", "error", err)
		renderErrorResponse(w, r, http.StatusInternalServerError, "Failed to unlink device from login", err.Error())
		return
	}

	// Return success response
	response := UnlinkDeviceFromLoginResponse{
		Status:  "success",
		Message: "Device unlinked successfully",
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// Handler returns a http.Handler for the device API
func Handler(h *DeviceHandler) http.Handler {
	r := chi.NewRouter()

	r.Get("/login/{login_id}", h.GetDevicesByLogin)
	r.Post("/login/{login_id}/unlink", h.UnlinkDeviceFromLogin)

	return r
}

// renderErrorResponse renders an error response with the given status code and message
func renderErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, message, errorDetail string) {
	response := ErrorResponse{
		Status:  "error",
		Message: message,
	}

	if errorDetail != "" {
		response.Error = errorDetail
	}

	render.Status(r, statusCode)
	render.JSON(w, r, response)
}
