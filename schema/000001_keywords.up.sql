CREATE TABLE IF NOT EXISTS users (
    Id INTEGER PRIMARY KEY AUTOINCREMENT,
    Email TEXT,
	Name     TEXT,
	password TEXT,
	Role     INTEGER
);

CREATE TABLE IF NOT EXISTS keyword (
    number INTEGER PRIMARY KEY,
    url TEXT,
    keywords TEXT
);

CREATE TABLE IF NOT EXISTS words_index (
    word TEXT PRIMARY KEY,
    numbers TEXT
);