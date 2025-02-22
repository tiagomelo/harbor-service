CREATE TABLE IF NOT EXISTS harbors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    unloc TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    city TEXT NOT NULL,
    country TEXT NOT NULL,
    alias TEXT,     
    regions TEXT, 
    coordinates TEXT,
    province TEXT,
    timezone TEXT,
    unlocs TEXT,
    code TEXT
);