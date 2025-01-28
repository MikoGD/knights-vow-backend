SELECT f.id, f.name, f.created_date, u.id, u.username
FROM files f
JOIN users u ON u.id = f.user_id
WHERE f.id = ?