CREATE TABLE IF NOT EXISTS shortlink (
    slink STRING PRIMARY KEY,
    llink STRING NOT NULL,
    CHECK (llink <> '')
);
