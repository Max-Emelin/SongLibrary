CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    "group" TEXT NOT NULL,
    song_name TEXT NOT NULL,
    release_date DATE,
    link TEXT
);

CREATE TABLE lyrics (
    id SERIAL PRIMARY KEY,
    song_id INT REFERENCES songs(id) ON DELETE CASCADE,
    verse_number INT NOT NULL,
    text TEXT NOT NULL
);

INSERT INTO songs (id, "group", song_name, release_date, link)
VALUES (0,'group','song_name', NULL, '');