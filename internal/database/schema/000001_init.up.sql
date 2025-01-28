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

INSERT INTO songs ("group", "song_name", "release_date", "link")
VALUES ('group','song_name', NULL, '');

INSERT INTO lyrics ("song_id", "verse_number", "text")
VALUES ('1','1','Ooh baby, don''t you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?');

INSERT INTO lyrics ("song_id", "verse_number", "text")
VALUES ('1','2','Ooh\nYou set my soul alight\nOoh\nYou set my soul alight');