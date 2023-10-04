CREATE TABLE IF NOT EXISTS shortlink (
    slink TEXT PRIMARY KEY,
    llink TEXT NOT NULL,
    CHECK (llink <> '')
);
