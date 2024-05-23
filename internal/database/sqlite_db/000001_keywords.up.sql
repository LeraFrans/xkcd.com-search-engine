CREATE TABLE IF NOT EXISTS keyword (
    number INTEGER PRIMARY KEY,
    url TEXT,
    keywords TEXT
);

CREATE TABLE IF NOT EXISTS words_index (
    word TEXT PRIMARY KEY,
    numbers TEXT
);