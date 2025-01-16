-- sqlite
CREATE TABLE IF NOT EXISTS Files (
  id INTEGER PRIMARY KEY,
  user_id INTEGER,
  name TEXT,
  created_date TEXT,
  FOREIGN KEY(user_id) REFERENCES users(id)
);