-- name: GetRoleByName :one
SELECt *
FROM roles
WHERE name = $1;