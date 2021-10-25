-- -- name: CreateEntry :exec
-- INSERT INTO entry (
--   accountId,
--   amount
-- ) VALUES (
--   $1, $2 
-- )
-- RETURNING *;

-- -- name: GetEntries :many
-- SELECT * FROM entry
-- WHERE accountId = $1
-- ORDER by id;

-- -- -- name: UpdateEntry :exec
-- -- UPDATE authors SET bio = $2
-- -- WHERE id = $1;


