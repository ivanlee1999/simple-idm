// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package logindb

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const findLoginByUsername = `-- name: FindLoginByUsername :one
SELECT l.uuid, l.username, l.password, l.created_at, l.updated_at
FROM login l
WHERE l.username = $1
AND l.deleted_at IS NULL
`

type FindLoginByUsernameRow struct {
	Uuid      uuid.UUID      `json:"uuid"`
	Username  sql.NullString `json:"username"`
	Password  []byte         `json:"password"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (q *Queries) FindLoginByUsername(ctx context.Context, username sql.NullString) (FindLoginByUsernameRow, error) {
	row := q.db.QueryRow(ctx, findLoginByUsername, username)
	var i FindLoginByUsernameRow
	err := row.Scan(
		&i.Uuid,
		&i.Username,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findUser = `-- name: FindUser :one
SELECT uuid, username, password
FROM login
WHERE username = $1
AND deleted_at IS NULL
`

type FindUserRow struct {
	Uuid     uuid.UUID      `json:"uuid"`
	Username sql.NullString `json:"username"`
	Password []byte         `json:"password"`
}

func (q *Queries) FindUser(ctx context.Context, username sql.NullString) (FindUserRow, error) {
	row := q.db.QueryRow(ctx, findUser, username)
	var i FindUserRow
	err := row.Scan(&i.Uuid, &i.Username, &i.Password)
	return i, err
}

const findUserInfoWithRoles = `-- name: FindUserInfoWithRoles :one
SELECT u.email, u.username, u.name, COALESCE(array_agg(r.name), '{}') AS roles
FROM users u
LEFT JOIN user_roles ur ON u.uuid = ur.user_uuid
LEFT JOIN roles r ON ur.role_uuid = r.uuid
WHERE u.uuid = $1
AND u.deleted_at IS NULL
GROUP BY u.email, u.username, u.name
`

type FindUserInfoWithRolesRow struct {
	Email    string         `json:"email"`
	Username sql.NullString `json:"username"`
	Name     sql.NullString `json:"name"`
	Roles    interface{}    `json:"roles"`
}

func (q *Queries) FindUserInfoWithRoles(ctx context.Context, argUuid uuid.UUID) (FindUserInfoWithRolesRow, error) {
	row := q.db.QueryRow(ctx, findUserInfoWithRoles, argUuid)
	var i FindUserInfoWithRolesRow
	err := row.Scan(
		&i.Email,
		&i.Username,
		&i.Name,
		&i.Roles,
	)
	return i, err
}

const findUserRolesByUserUuid = `-- name: FindUserRolesByUserUuid :many
SELECT r.name
FROM user_roles ur
LEFT JOIN roles r ON ur.role_uuid = r.uuid
WHERE ur.user_uuid = $1
`

func (q *Queries) FindUserRolesByUserUuid(ctx context.Context, userUuid uuid.UUID) ([]sql.NullString, error) {
	rows, err := q.db.Query(ctx, findUserRolesByUserUuid, userUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []sql.NullString
	for rows.Next() {
		var name sql.NullString
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findUsernameByEmail = `-- name: FindUsernameByEmail :one
SELECT u.username
FROM users u
WHERE u.email = $1
AND u.deleted_at IS NULL
`

func (q *Queries) FindUsernameByEmail(ctx context.Context, email string) (sql.NullString, error) {
	row := q.db.QueryRow(ctx, findUsernameByEmail, email)
	var username sql.NullString
	err := row.Scan(&username)
	return username, err
}

const get2FAByLoginUuid = `-- name: Get2FAByLoginUuid :one
SELECT u.two_factor_secret
FROM users u
JOIN login l ON l.user_uuid = u.uuid
WHERE l.uuid = $1
AND u.deleted_at IS NULL
`

func (q *Queries) Get2FAByLoginUuid(ctx context.Context, argUuid uuid.UUID) (pgtype.Text, error) {
	row := q.db.QueryRow(ctx, get2FAByLoginUuid, argUuid)
	var two_factor_secret pgtype.Text
	err := row.Scan(&two_factor_secret)
	return two_factor_secret, err
}

const get2FASecret = `-- name: Get2FASecret :one
SELECT two_factor_secret
FROM users
WHERE users.uuid = $1
AND deleted_at IS NULL
`

func (q *Queries) Get2FASecret(ctx context.Context, argUuid uuid.UUID) (pgtype.Text, error) {
	row := q.db.QueryRow(ctx, get2FASecret, argUuid)
	var two_factor_secret pgtype.Text
	err := row.Scan(&two_factor_secret)
	return two_factor_secret, err
}

const getLoginByUUID = `-- name: GetLoginByUUID :one
SELECT l.uuid as login_uuid, l.username, l.password, l.created_at, l.updated_at,
       l.two_factor_enabled, l.two_factor_secret, l.two_factor_backup_codes
FROM login l
WHERE l.uuid = $1
AND l.deleted_at IS NULL
`

type GetLoginByUUIDRow struct {
	LoginUuid            uuid.UUID      `json:"login_uuid"`
	Username             sql.NullString `json:"username"`
	Password             []byte         `json:"password"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	TwoFactorEnabled     pgtype.Bool    `json:"two_factor_enabled"`
	TwoFactorSecret      pgtype.Text    `json:"two_factor_secret"`
	TwoFactorBackupCodes []string       `json:"two_factor_backup_codes"`
}

func (q *Queries) GetLoginByUUID(ctx context.Context, argUuid uuid.UUID) (GetLoginByUUIDRow, error) {
	row := q.db.QueryRow(ctx, getLoginByUUID, argUuid)
	var i GetLoginByUUIDRow
	err := row.Scan(
		&i.LoginUuid,
		&i.Username,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.TwoFactorEnabled,
		&i.TwoFactorSecret,
		&i.TwoFactorBackupCodes,
	)
	return i, err
}

const getUsersByLoginUuid = `-- name: GetUsersByLoginUuid :many
SELECT u.uuid, u.username, u.name, u.email, u.created_at, u.last_modified_at,
       COALESCE(array_agg(r.name) FILTER (WHERE r.name IS NOT NULL), '{}') as roles
FROM users u
LEFT JOIN user_roles ur ON u.uuid = ur.user_uuid
LEFT JOIN roles r ON ur.role_uuid = r.uuid
WHERE u.login_uuid = $1
AND u.deleted_at IS NULL
GROUP BY u.uuid, u.username, u.name, u.email, u.created_at, u.last_modified_at
`

type GetUsersByLoginUuidRow struct {
	Uuid           uuid.UUID      `json:"uuid"`
	Username       sql.NullString `json:"username"`
	Name           sql.NullString `json:"name"`
	Email          string         `json:"email"`
	CreatedAt      time.Time      `json:"created_at"`
	LastModifiedAt time.Time      `json:"last_modified_at"`
	Roles          interface{}    `json:"roles"`
}

func (q *Queries) GetUsersByLoginUuid(ctx context.Context, loginUuid uuid.NullUUID) ([]GetUsersByLoginUuidRow, error) {
	rows, err := q.db.Query(ctx, getUsersByLoginUuid, loginUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUsersByLoginUuidRow
	for rows.Next() {
		var i GetUsersByLoginUuidRow
		if err := rows.Scan(
			&i.Uuid,
			&i.Username,
			&i.Name,
			&i.Email,
			&i.CreatedAt,
			&i.LastModifiedAt,
			&i.Roles,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const initPasswordByUsername = `-- name: InitPasswordByUsername :one
SELECT uuid
FROM login
WHERE username = $1
`

func (q *Queries) InitPasswordByUsername(ctx context.Context, username sql.NullString) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, initPasswordByUsername, username)
	var uuid uuid.UUID
	err := row.Scan(&uuid)
	return uuid, err
}

const initPasswordResetToken = `-- name: InitPasswordResetToken :exec
INSERT INTO password_reset_tokens (user_uuid, token, expire_at)
VALUES ($1, $2, $3)
`

type InitPasswordResetTokenParams struct {
	UserUuid uuid.UUID          `json:"user_uuid"`
	Token    string             `json:"token"`
	ExpireAt pgtype.Timestamptz `json:"expire_at"`
}

func (q *Queries) InitPasswordResetToken(ctx context.Context, arg InitPasswordResetTokenParams) error {
	_, err := q.db.Exec(ctx, initPasswordResetToken, arg.UserUuid, arg.Token, arg.ExpireAt)
	return err
}

const markBackupCodeUsed = `-- name: MarkBackupCodeUsed :exec
UPDATE login l
SET two_factor_backup_codes = array_remove(two_factor_backup_codes, $1::text)
WHERE l.uuid = $2
AND l.deleted_at IS NULL
`

type MarkBackupCodeUsedParams struct {
	Code string    `json:"code"`
	Uuid uuid.UUID `json:"uuid"`
}

func (q *Queries) MarkBackupCodeUsed(ctx context.Context, arg MarkBackupCodeUsedParams) error {
	_, err := q.db.Exec(ctx, markBackupCodeUsed, arg.Code, arg.Uuid)
	return err
}

const markPasswordResetTokenUsed = `-- name: MarkPasswordResetTokenUsed :exec
UPDATE password_reset_tokens
SET used_at = NOW()
WHERE token = $1
`

func (q *Queries) MarkPasswordResetTokenUsed(ctx context.Context, token string) error {
	_, err := q.db.Exec(ctx, markPasswordResetTokenUsed, token)
	return err
}

const resetPassword = `-- name: ResetPassword :exec
UPDATE login
SET password = $1, 
    last_modified_at = NOW()
WHERE username = $2
`

type ResetPasswordParams struct {
	Password []byte         `json:"password"`
	Username sql.NullString `json:"username"`
}

func (q *Queries) ResetPassword(ctx context.Context, arg ResetPasswordParams) error {
	_, err := q.db.Exec(ctx, resetPassword, arg.Password, arg.Username)
	return err
}

const resetPasswordByUuid = `-- name: ResetPasswordByUuid :exec
UPDATE login
SET password = $1,
    updated_at = NOW()
WHERE login.uuid = $2
`

type ResetPasswordByUuidParams struct {
	Password []byte    `json:"password"`
	Uuid     uuid.UUID `json:"uuid"`
}

func (q *Queries) ResetPasswordByUuid(ctx context.Context, arg ResetPasswordByUuidParams) error {
	_, err := q.db.Exec(ctx, resetPasswordByUuid, arg.Password, arg.Uuid)
	return err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE login
SET password = $1,
    updated_at = NOW()
WHERE uuid = $2
`

type UpdateUserPasswordParams struct {
	Password []byte    `json:"password"`
	Uuid     uuid.UUID `json:"uuid"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.Exec(ctx, updateUserPassword, arg.Password, arg.Uuid)
	return err
}

const validateBackupCode = `-- name: ValidateBackupCode :one
SELECT EXISTS (
  SELECT 1
  FROM login l
  WHERE l.uuid = $1
  AND $2::text = ANY(l.two_factor_backup_codes)
  AND l.deleted_at IS NULL
) AS is_valid
`

type ValidateBackupCodeParams struct {
	Uuid uuid.UUID `json:"uuid"`
	Code string    `json:"code"`
}

func (q *Queries) ValidateBackupCode(ctx context.Context, arg ValidateBackupCodeParams) (bool, error) {
	row := q.db.QueryRow(ctx, validateBackupCode, arg.Uuid, arg.Code)
	var is_valid bool
	err := row.Scan(&is_valid)
	return is_valid, err
}

const validatePasswordResetToken = `-- name: ValidatePasswordResetToken :one
SELECT prt.uuid as uuid, prt.user_uuid as user_uuid
FROM password_reset_tokens prt
JOIN users u ON u.uuid = prt.user_uuid 
WHERE prt.token = $1
  AND prt.expire_at > NOW()
  AND prt.used_at IS NULL
LIMIT 1
`

type ValidatePasswordResetTokenRow struct {
	Uuid     uuid.UUID `json:"uuid"`
	UserUuid uuid.UUID `json:"user_uuid"`
}

func (q *Queries) ValidatePasswordResetToken(ctx context.Context, token string) (ValidatePasswordResetTokenRow, error) {
	row := q.db.QueryRow(ctx, validatePasswordResetToken, token)
	var i ValidatePasswordResetTokenRow
	err := row.Scan(&i.Uuid, &i.UserUuid)
	return i, err
}
