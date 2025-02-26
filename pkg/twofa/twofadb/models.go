// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package twofadb

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type BackupCode struct {
	ID        uuid.UUID          `json:"id"`
	UserID    uuid.UUID          `json:"user_id"`
	Code      string             `json:"code"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UsedAt    pgtype.Timestamptz `json:"used_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
}

type GooseDbVersion struct {
	ID        int32     `json:"id"`
	VersionID int64     `json:"version_id"`
	IsApplied bool      `json:"is_applied"`
	Tstamp    time.Time `json:"tstamp"`
}

type Login struct {
	ID                   uuid.UUID      `json:"id"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            sql.NullTime   `json:"deleted_at"`
	CreatedBy            sql.NullString `json:"created_by"`
	Password             []byte         `json:"password"`
	Username             sql.NullString `json:"username"`
	TwoFactorSecret      pgtype.Text    `json:"two_factor_secret"`
	TwoFactorEnabled     pgtype.Bool    `json:"two_factor_enabled"`
	TwoFactorBackupCodes []string       `json:"two_factor_backup_codes"`
}

type Login2fa struct {
	ID                   uuid.UUID      `json:"id"`
	LoginID              uuid.UUID      `json:"login_id"`
	TwoFactorSecret      pgtype.Text    `json:"two_factor_secret"`
	TwoFactorEnabled     pgtype.Bool    `json:"two_factor_enabled"`
	TwoFactorType        sql.NullString `json:"two_factor_type"`
	TwoFactorBackupCodes []string       `json:"two_factor_backup_codes"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            sql.NullTime   `json:"updated_at"`
	DeletedAt            sql.NullTime   `json:"deleted_at"`
}

type LoginPasswordResetToken struct {
	ID        uuid.UUID          `json:"id"`
	Token     string             `json:"token"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	ExpireAt  pgtype.Timestamptz `json:"expire_at"`
	UsedAt    pgtype.Timestamptz `json:"used_at"`
	LoginID   uuid.UUID          `json:"login_id"`
}

type Role struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
}

type User struct {
	ID                   uuid.UUID      `json:"id"`
	CreatedAt            time.Time      `json:"created_at"`
	LastModifiedAt       time.Time      `json:"last_modified_at"`
	DeletedAt            sql.NullTime   `json:"deleted_at"`
	CreatedBy            sql.NullString `json:"created_by"`
	Email                string         `json:"email"`
	Name                 sql.NullString `json:"name"`
	Password             []byte         `json:"password"`
	VerifiedAt           sql.NullTime   `json:"verified_at"`
	Username             sql.NullString `json:"username"`
	TwoFactorSecret      pgtype.Text    `json:"two_factor_secret"`
	TwoFactorEnabled     pgtype.Bool    `json:"two_factor_enabled"`
	TwoFactorBackupCodes []string       `json:"two_factor_backup_codes"`
	LoginID              uuid.NullUUID  `json:"login_id"`
}

type UserRole struct {
	UserID uuid.UUID `json:"user_id"`
	RoleID uuid.UUID `json:"role_id"`
}
