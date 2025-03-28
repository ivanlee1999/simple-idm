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

// Error defines model for Error.
type Error struct {
	// Error code
	Code string `json:"code"`

	// Error message
	Message string `json:"message"`
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

// TwoFactorMethods defines model for TwoFactorMethods.
type TwoFactorMethods struct {
	Count   int               `json:"count"`
	Methods []TwoFactorMethod `json:"methods"`
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

// ChangePasswordJSONBody defines parameters for ChangePassword.
type ChangePasswordJSONBody struct {
	// User's current password
	CurrentPassword string `json:"current_password"`

	// User's new password
	NewPassword string `json:"new_password"`
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

// ChangePasswordJSONRequestBody defines body for ChangePassword for application/json ContentType.
type ChangePasswordJSONRequestBody ChangePasswordJSONBody

// Bind implements render.Binder.
func (ChangePasswordJSONRequestBody) Bind(*http.Request) error {
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
	// Change user password
	// (PUT /password)
	ChangePassword(w http.ResponseWriter, r *http.Request) *Response
	// Get password policy
	// (GET /password/policy)
	GetPasswordPolicy(w http.ResponseWriter, r *http.Request) *Response
	// Change username
	// (PUT /username)
	ChangeUsername(w http.ResponseWriter, r *http.Request) *Response
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
		r.Put("/password", wrapper.ChangePassword)
		r.Get("/password/policy", wrapper.GetPasswordPolicy)
		r.Put("/username", wrapper.ChangeUsername)
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
	"P9WG4E7TarVa/R4AAP//0Ye0uPseAAA=",
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
