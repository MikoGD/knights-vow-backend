-- sqlite3
SELECT f.id, f.user_id, f.name, f.created_date, u.username
FROM Files f
JOIN Users u ON f.user_id = u.id;