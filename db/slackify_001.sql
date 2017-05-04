CREATE TABLE IF NOT EXISTS user (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username TEXT NOT NULL,
    previous_status TEXT,
    previous_emoji TEXT,
    spotify_emoji TEXT,
    state INT,
    enabled BOOL,
    slack_access_token TEXT,
    spotify_refresh_token TEXT
);
