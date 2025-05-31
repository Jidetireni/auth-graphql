-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash,
    phone_number,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, NOW(), NOW()
)
RETURNING id, email, phone_number, created_at, updated_at;

-- name: GetUserByEmail :one
-- Selects a user by their email address, excluding soft-deleted users.
SELECT id, email, password_hash, phone_number, created_at, updated_at, deleted_at
FROM users
WHERE email = $1 AND deleted_at IS NULL;

-- name: GetUserByID :one
-- Selects a user by their ID, excluding soft-deleted users.
SELECT id, email, password_hash, phone_number, created_at, updated_at, deleted_at
FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetAllUsers :many
-- Selects all non-soft-deleted users.
SELECT id, email, phone_number, created_at, updated_at, deleted_at
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateUser :one
-- Updates user details. Note: password updates should be handled carefully.
-- This example allows updating email and phone_number.
UPDATE users
SET
    email = $2,
    phone_number = $3,
    updated_at = NOW()
WHERE
    id = $1 AND deleted_at IS NULL
RETURNING id, email, phone_number, created_at, updated_at, deleted_at;

-- name: SoftDeleteUser :exec
-- Marks a user as deleted by setting the deleted_at timestamp.
UPDATE users
SET
    deleted_at = NOW(),
    updated_at = NOW()
WHERE
    id = $1 AND deleted_at IS NULL;

-- name: HardDeleteUser :exec
-- Permanently deletes a user from the database. Use with caution.
DELETE FROM users
WHERE id = $1;