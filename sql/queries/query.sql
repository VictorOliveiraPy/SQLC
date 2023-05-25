

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



-- name: CreateCourse :exec
INSERT INTO courses (id, name, description, price) VALUES ($1, $2, $3, $4) RETURNING *;


-- name: ListCourses :many
SELECT c.*, ca.name as category_name 
FROM courses c JOIN categories ca ON c.category_id = ca.id;