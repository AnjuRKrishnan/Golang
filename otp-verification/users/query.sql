-- query.sql
-- name: Create :execresult
INSERT INTO users (name, phone_number) VALUES ($1, $2);

-- name: Get :one
SELECT * FROM users WHERE phone_number = $1;

-- name: Update :exec
UPDATE users SET otp = $1, otp_expiration_time = $2 WHERE id = $3 ;
