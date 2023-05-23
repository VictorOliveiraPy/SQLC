

-- name: GetCategories :one
SELECT * FROM categories
WHERE id =  $1;

-- name: ListCategories :many
SELECT * FROM categories;


-- name: CreateCategory :exec
INSERT INTO categories (id, name, description) VALUES ($1, $2, $3) RETURNING *;


-- name: UpdateCategory :one
UPDATE categories
  set name = $2,
  description = $3
WHERE id = $1
RETURNING *;


-- name: DeleteCategory :one
DELETE FROM categories WHERE id = $1 RETURNING *;