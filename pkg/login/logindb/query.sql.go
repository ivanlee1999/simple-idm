// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package logindb

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const addPasswordToHistory = `-- name: AddPasswordToHistory :exec
INSERT INTO login_password_history (login_id, password_hash, password_version)
VALUES ($1, $2, $3)
`

type AddPasswordToHistoryParams struct {
	LoginID         uuid.UUID `json:"login_id"`
	PasswordHash    []byte    `json:"password_hash"`
	PasswordVersion int32     `json:"password_version"`
}

func (q *Queries) AddPasswordToHistory(ctx context.Context, arg AddPasswordToHistoryParams) error {
	_, err := q.db.Exec(ctx, addPasswordToHistory, arg.LoginID, arg.PasswordHash, arg.PasswordVersion)
	return err
}

const findEmailByEmail = `-- name: FindEmailByEmail :one
SELECT u.email
FROM users u
WHERE u.email = $1
AND u.deleted_at IS NULL
`

func (q *Queries) FindEmailByEmail(ctx context.Context, email string) (string, error) {
	row := q.db.QueryRow(ctx, findEmailByEmail, email)
	err := row.Scan(&email)
	return email, err
}

const findLoginByUsername = `-- name: FindLoginByUsername :one
SELECT l.id, l.username, l.password, l.password_version, l.created_at, l.updated_at
FROM login l
WHERE l.username = $1
AND l.deleted_at IS NULL
`

type FindLoginByUsernameRow struct {
	ID              uuid.UUID      `json:"id"`
	Username        sql.NullString `json:"username"`
	Password        []byte         `json:"password"`
	PasswordVersion pgtype.Int4    `json:"password_version"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (q *Queries) FindLoginByUsername(ctx context.Context, username sql.NullString) (FindLoginByUsernameRow, error) {
	row := q.db.QueryRow(ctx, findLoginByUsername, username)
	var i FindLoginByUsernameRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.PasswordVersion,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findUser = `-- name: FindUser :one
SELECT id, username, password, password_version
FROM login
WHERE username = $1
AND deleted_at IS NULL
`

type FindUserRow struct {
	ID              uuid.UUID      `json:"id"`
	Username        sql.NullString `json:"username"`
	Password        []byte         `json:"password"`
	PasswordVersion pgtype.Int4    `json:"password_version"`
}

func (q *Queries) FindUser(ctx context.Context, username sql.NullString) (FindUserRow, error) {
	row := q.db.QueryRow(ctx, findUser, username)
	var i FindUserRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.PasswordVersion,
	)
	return i, err
}

const findUserInfoWithRoles = `-- name: FindUserInfoWithRoles :one
SELECT u.email, u.name, COALESCE(array_agg(r.name), '{}') AS roles
FROM users u
LEFT JOIN user_roles ur ON u.id = ur.user_id
LEFT JOIN roles r ON ur.role_id = r.id
WHERE u.id = $1
AND u.deleted_at IS NULL
GROUP BY u.email, u.name
`

type FindUserInfoWithRolesRow struct {
	Email string         `json:"email"`
	Name  sql.NullString `json:"name"`
	Roles interface{}    `json:"roles"`
}

func (q *Queries) FindUserInfoWithRoles(ctx context.Context, id uuid.UUID) (FindUserInfoWithRolesRow, error) {
	row := q.db.QueryRow(ctx, findUserInfoWithRoles, id)
	var i FindUserInfoWithRolesRow
	err := row.Scan(&i.Email, &i.Name, &i.Roles)
	return i, err
}

const findUserRolesByUserId = `-- name: FindUserRolesByUserId :many
SELECT r.name
FROM user_roles ur
LEFT JOIN roles r ON ur.role_id = r.id
WHERE ur.user_id = $1
`

func (q *Queries) FindUserRolesByUserId(ctx context.Context, userID uuid.UUID) ([]sql.NullString, error) {
	rows, err := q.db.Query(ctx, findUserRolesByUserId, userID)
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
SELECT l.username
FROM login l
JOIN users u ON u.login_id = l.id
WHERE u.email = $1
AND u.deleted_at IS NULL
LIMIT 1
`

func (q *Queries) FindUsernameByEmail(ctx context.Context, email string) (sql.NullString, error) {
	row := q.db.QueryRow(ctx, findUsernameByEmail, email)
	var username sql.NullString
	err := row.Scan(&username)
	return username, err
}

const get2FAByLoginId = `-- name: Get2FAByLoginId :one
SELECT NULL::text as two_factor_secret
FROM users u
JOIN login l ON l.id = l.id -- Self-join as a placeholder
WHERE l.id = $1
AND u.deleted_at IS NULL
LIMIT 1
`

// This query is no longer valid as two_factor_secret has been removed
// Keeping the query name for compatibility but returning NULL
func (q *Queries) Get2FAByLoginId(ctx context.Context, id uuid.UUID) (pgtype.Text, error) {
	row := q.db.QueryRow(ctx, get2FAByLoginId, id)
	var two_factor_secret pgtype.Text
	err := row.Scan(&two_factor_secret)
	return two_factor_secret, err
}

const get2FASecret = `-- name: Get2FASecret :one
SELECT NULL::text as two_factor_secret
FROM users
WHERE users.id = $1
AND deleted_at IS NULL
`

// This query is no longer valid as two_factor_secret has been removed
// Keeping the query name for compatibility but returning NULL
func (q *Queries) Get2FASecret(ctx context.Context, id uuid.UUID) (pgtype.Text, error) {
	row := q.db.QueryRow(ctx, get2FASecret, id)
	var two_factor_secret pgtype.Text
	err := row.Scan(&two_factor_secret)
	return two_factor_secret, err
}

const getLoginById = `-- name: GetLoginById :one
SELECT l.id as login_id, l.username, l.password, l.created_at, l.updated_at
FROM login l
WHERE l.id = $1
AND l.deleted_at IS NULL
`

type GetLoginByIdRow struct {
	LoginID   uuid.UUID      `json:"login_id"`
	Username  sql.NullString `json:"username"`
	Password  []byte         `json:"password"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (q *Queries) GetLoginById(ctx context.Context, id uuid.UUID) (GetLoginByIdRow, error) {
	row := q.db.QueryRow(ctx, getLoginById, id)
	var i GetLoginByIdRow
	err := row.Scan(
		&i.LoginID,
		&i.Username,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getLoginByUserId = `-- name: GetLoginByUserId :one
SELECT l.id as login_id, l.username, l.password, l.created_at, l.updated_at
FROM login l
JOIN users u ON l.id = u.login_id
WHERE u.id = $1
AND l.deleted_at IS NULL
`

type GetLoginByUserIdRow struct {
	LoginID   uuid.UUID      `json:"login_id"`
	Username  sql.NullString `json:"username"`
	Password  []byte         `json:"password"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (q *Queries) GetLoginByUserId(ctx context.Context, id uuid.UUID) (GetLoginByUserIdRow, error) {
	row := q.db.QueryRow(ctx, getLoginByUserId, id)
	var i GetLoginByUserIdRow
	err := row.Scan(
		&i.LoginID,
		&i.Username,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPasswordHistory = `-- name: GetPasswordHistory :many
SELECT id, login_id, password_hash, password_version, created_at
FROM login_password_history
WHERE login_id = $1
ORDER BY created_at DESC
LIMIT $2
`

type GetPasswordHistoryParams struct {
	LoginID uuid.UUID `json:"login_id"`
	Limit   int32     `json:"limit"`
}

type GetPasswordHistoryRow struct {
	ID              uuid.UUID `json:"id"`
	LoginID         uuid.UUID `json:"login_id"`
	PasswordHash    []byte    `json:"password_hash"`
	PasswordVersion int32     `json:"password_version"`
	CreatedAt       time.Time `json:"created_at"`
}

func (q *Queries) GetPasswordHistory(ctx context.Context, arg GetPasswordHistoryParams) ([]GetPasswordHistoryRow, error) {
	rows, err := q.db.Query(ctx, getPasswordHistory, arg.LoginID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPasswordHistoryRow
	for rows.Next() {
		var i GetPasswordHistoryRow
		if err := rows.Scan(
			&i.ID,
			&i.LoginID,
			&i.PasswordHash,
			&i.PasswordVersion,
			&i.CreatedAt,
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

const getPasswordVersion = `-- name: GetPasswordVersion :one
SELECT password_version
FROM login
WHERE id = $1
AND deleted_at IS NULL
`

func (q *Queries) GetPasswordVersion(ctx context.Context, id uuid.UUID) (pgtype.Int4, error) {
	row := q.db.QueryRow(ctx, getPasswordVersion, id)
	var password_version pgtype.Int4
	err := row.Scan(&password_version)
	return password_version, err
}

const getUsersByLoginId = `-- name: GetUsersByLoginId :many
SELECT u.id, u.name, u.email, u.created_at, u.last_modified_at,
       COALESCE(array_agg(r.name) FILTER (WHERE r.name IS NOT NULL), '{}') as roles
FROM users u
LEFT JOIN user_roles ur ON u.id = ur.user_id
LEFT JOIN roles r ON ur.role_id = r.id
WHERE u.login_id = $1
AND u.deleted_at IS NULL
GROUP BY u.id, u.name, u.email, u.created_at, u.last_modified_at
`

type GetUsersByLoginIdRow struct {
	ID             uuid.UUID      `json:"id"`
	Name           sql.NullString `json:"name"`
	Email          string         `json:"email"`
	CreatedAt      time.Time      `json:"created_at"`
	LastModifiedAt time.Time      `json:"last_modified_at"`
	Roles          interface{}    `json:"roles"`
}

func (q *Queries) GetUsersByLoginId(ctx context.Context, loginID uuid.NullUUID) ([]GetUsersByLoginIdRow, error) {
	rows, err := q.db.Query(ctx, getUsersByLoginId, loginID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUsersByLoginIdRow
	for rows.Next() {
		var i GetUsersByLoginIdRow
		if err := rows.Scan(
			&i.ID,
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
SELECT id
FROM login
WHERE username = $1
`

func (q *Queries) InitPasswordByUsername(ctx context.Context, username sql.NullString) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, initPasswordByUsername, username)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const initPasswordResetToken = `-- name: InitPasswordResetToken :exec
INSERT INTO login_password_reset_tokens (login_id, token, expire_at)
VALUES ($1, $2, $3)
`

type InitPasswordResetTokenParams struct {
	LoginID  uuid.UUID          `json:"login_id"`
	Token    string             `json:"token"`
	ExpireAt pgtype.Timestamptz `json:"expire_at"`
}

func (q *Queries) InitPasswordResetToken(ctx context.Context, arg InitPasswordResetTokenParams) error {
	_, err := q.db.Exec(ctx, initPasswordResetToken, arg.LoginID, arg.Token, arg.ExpireAt)
	return err
}

const markBackupCodeUsed = `-- name: MarkBackupCodeUsed :exec
SELECT 1
`

// This query is no longer valid as two_factor_backup_codes has been removed
// Keeping the query name for compatibility but doing nothing
func (q *Queries) MarkBackupCodeUsed(ctx context.Context) error {
	_, err := q.db.Exec(ctx, markBackupCodeUsed)
	return err
}

const markPasswordResetTokenUsed = `-- name: MarkPasswordResetTokenUsed :exec
UPDATE login_password_reset_tokens
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

const resetPasswordById = `-- name: ResetPasswordById :exec
UPDATE login
SET password = $1,
    updated_at = NOW()
WHERE login.id = $2
`

type ResetPasswordByIdParams struct {
	Password []byte    `json:"password"`
	ID       uuid.UUID `json:"id"`
}

func (q *Queries) ResetPasswordById(ctx context.Context, arg ResetPasswordByIdParams) error {
	_, err := q.db.Exec(ctx, resetPasswordById, arg.Password, arg.ID)
	return err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE login
SET password = $1,
    updated_at = NOW()
WHERE id = $2
`

type UpdateUserPasswordParams struct {
	Password []byte    `json:"password"`
	ID       uuid.UUID `json:"id"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.Exec(ctx, updateUserPassword, arg.Password, arg.ID)
	return err
}

const updateUserPasswordAndVersion = `-- name: UpdateUserPasswordAndVersion :exec
UPDATE login
SET password = $1,
    password_version = $3,
    updated_at = NOW()
WHERE id = $2
`

type UpdateUserPasswordAndVersionParams struct {
	Password        []byte      `json:"password"`
	ID              uuid.UUID   `json:"id"`
	PasswordVersion pgtype.Int4 `json:"password_version"`
}

func (q *Queries) UpdateUserPasswordAndVersion(ctx context.Context, arg UpdateUserPasswordAndVersionParams) error {
	_, err := q.db.Exec(ctx, updateUserPasswordAndVersion, arg.Password, arg.ID, arg.PasswordVersion)
	return err
}

const validateBackupCode = `-- name: ValidateBackupCode :one
SELECT false AS is_valid
`

// This query is no longer valid as two_factor_backup_codes has been removed
// Keeping the query name for compatibility but returning false
func (q *Queries) ValidateBackupCode(ctx context.Context) (bool, error) {
	row := q.db.QueryRow(ctx, validateBackupCode)
	var is_valid bool
	err := row.Scan(&is_valid)
	return is_valid, err
}

const validatePasswordResetToken = `-- name: ValidatePasswordResetToken :one
SELECT prt.id as id, prt.login_id as login_id
FROM login_password_reset_tokens prt
JOIN login l ON l.id = prt.login_id 
WHERE prt.token = $1
  AND prt.expire_at > NOW()
  AND prt.used_at IS NULL
LIMIT 1
`

type ValidatePasswordResetTokenRow struct {
	ID      uuid.UUID `json:"id"`
	LoginID uuid.UUID `json:"login_id"`
}

func (q *Queries) ValidatePasswordResetToken(ctx context.Context, token string) (ValidatePasswordResetTokenRow, error) {
	row := q.db.QueryRow(ctx, validatePasswordResetToken, token)
	var i ValidatePasswordResetTokenRow
	err := row.Scan(&i.ID, &i.LoginID)
	return i, err
}
