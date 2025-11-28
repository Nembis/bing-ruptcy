-- name: CreateUser :one

INSERT INTO users(id, user_name, email, created_at, updated_at, magic_link, magic_link_created_at)
VALUES (
	gen_random_uuid(),
	$1,
	$2,
	NOW(),
	NOW(),
	$3,
	NOW()
	)
RETURNING *;

-- name: GetUser :one

SELECT *
FROM users
WHERE magic_link = $1;

-- name: InsertMagicLink :one

UPDATE users
SET updated_at = now(), magic_link = $2, magic_link_created_at = now()
WHERE email = $1
RETURNING *; 

