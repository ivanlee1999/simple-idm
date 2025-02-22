package login

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/tendant/simple-idm/pkg/login/db"
	"github.com/tendant/simple-idm/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
)

type LoginService struct {
	queries *db.Queries
}

func New(queries *db.Queries) *LoginService {
	return &LoginService{
		queries: queries,
	}
}

type LoginParams struct {
	Email    string
	Username string
}

type IdmUser struct {
	UserUuid string   `json:"user_uuid,omitempty"`
	Role     []string `json:"role,omitempty"`
}

func (s LoginService) Login(ctx context.Context, params LoginParams) ([]db.FindUserByUsernameRow, error) {
	user, err := s.queries.FindUserByUsername(ctx, utils.ToNullString(params.Username))
	return user, err
}

type RegisterParam struct {
	Email    string
	Name     string
	Password string
}

// HashPassword hashes the plain-text password using bcrypt.
func HashPassword(password string) (string, error) {
	if password == "" { // Null or empty password check
		return "", fmt.Errorf("password cannot be empty")
	}

	// Generate the bcrypt hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash compares the plain-text password with the stored hashed password.
func CheckPasswordHash(password, hashedPassword string) (bool, error) {
	if password == "" || hashedPassword == "" { // Null or empty checks
		return false, fmt.Errorf("password and hashed password cannot be empty")
	}

	// Compare the password with the hashed password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s LoginService) Create(ctx context.Context, params RegisterParam) (db.User, error) {
	slog.Debug("Registering user use params:", "params", params)
	registerRequest := db.RegisterUserParams{}
	copier.Copy(&registerRequest, params)
	user, err := s.queries.RegisterUser(ctx, registerRequest)
	if err != nil {
		slog.Error("Failed to register user", "params", params, "err", err)
		return db.User{}, err
	}
	return user, err
}

func (s LoginService) EmailVerify(ctx context.Context, param string) error {
	slog.Debug("Verifying user use params:", "params", param)
	err := s.queries.EmailVerify(ctx, param)
	if err != nil {
		slog.Error("Failed to verify user", "params", param, "err", err)
		return err
	}
	return nil
}

func (s LoginService) ResetPasswordUsers(ctx context.Context, params PasswordReset) error {
	resetPasswordParams := db.ResetPasswordParams{}
	slog.Debug("resetPasswordParams", "params", params)
	copier.Copy(&resetPasswordParams, params)
	err := s.queries.ResetPassword(ctx, resetPasswordParams)
	return err
}

func (s LoginService) FindUserRoles(ctx context.Context, uuid uuid.UUID) ([]sql.NullString, error) {
	slog.Debug("FindUserRoles", "params", uuid)
	roles, err := s.queries.FindUserRolesByUserUuid(ctx, uuid)
	return roles, err
}

func (s LoginService) GetMe(ctx context.Context, userUuid uuid.UUID) (db.FindUserInfoWithRolesRow, error) {
	slog.Debug("GetMe", "userUuid", userUuid)
	userInfo, err := s.queries.FindUserInfoWithRoles(ctx, userUuid)
	if err != nil {
		slog.Error("Failed getting userinfo with roles", "err", err)
		return db.FindUserInfoWithRolesRow{}, err
	}
	return userInfo, err
}
