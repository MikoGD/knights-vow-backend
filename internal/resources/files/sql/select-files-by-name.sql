--- sqlite3
SELECT f.id, f.name, f.created_date, f.user_id, u.username
FROM files f
JOIN users u ON f.user_id = u.ID
WHERE f.name LIKE ?
UNION
SELECT f.id, f.name, f.created_date, f.user_id, u.username
FROM files f
JOIN users u ON f.user_id = u.id
WHERE f.name LIKE ?;