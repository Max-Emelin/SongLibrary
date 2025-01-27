package model

type Lyrics struct {
	ID          int    `json:"id" db:"id"`
	SongID      int    `json:"song_id" db:"song_id" binding:"required"`
	VerseNumber int    `json:"verse_number" db:"verse_number" binding:"required"`
	Text        string `json:"text" db:"text" binding:"required"`
}
