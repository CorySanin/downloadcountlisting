CREATE TABLE downloads (
    Path VARCHAR(128),
    Filename VARCHAR(128),
    Timestamp DATETIME
);

CREATE INDEX IF NOT EXISTS download_index ON downloads (Path, Filename)
