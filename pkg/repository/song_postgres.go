package repository

import (
	"SongLibrary/pkg/model"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type SongPostgres struct {
	db *sqlx.DB
}

func NewSongPostgres(db *sqlx.DB) *SongPostgres {
	return &SongPostgres{db: db}
}

func (r *SongPostgres) Create(song model.Song) (int, error) {
	var id int

	createSongQuery := fmt.Sprintf(`INSERT INTO %s (group, song_name) 
									VALUES ($1, $2) 
									RETURNING id`,
		songsTable)
	err := r.db.QueryRow(createSongQuery, song.Group, song.SongName).Scan(&id)

	return id, err
}

func (r *SongPostgres) GetLyrics(songId int, limit, offset int) ([]model.Lyrics, error) {
	var lyrics []model.Lyrics

	getLyricsQuery := fmt.Sprintf(`SELECT * 
								FROM %s 
								WHERE song_id = $1 
								ORDER BY verse_number 
								LIMIT $2 
								OFFSET $3`,
		lyricsTable)
	err := r.db.Select(&lyrics, getLyricsQuery, songId, limit, offset)

	return lyrics, err
}

func (r *SongPostgres) GetAllSongsWithFilter(filter model.SongFilter, limit, offset int) ([]model.Song, error) {
	var songs []model.Song

	getAllSongsWithFilterQuery := fmt.Sprintf(`SELECT *
											FROM %s
											WHERE ($1 IS NULL OR group ILIKE $1) 
												AND ($2 IS NULL OR song_name ILIKE $2) 
											LIMIT $3 
											OFFSET $4`,
		songsTable)
	err := r.db.Select(&songs, getAllSongsWithFilterQuery, filter.Group, filter.Song, limit, offset)

	return songs, err
}

func (r *SongPostgres) GetById(songId int) (model.Song, error) {
	var song model.Song

	getByIdQuery := fmt.Sprintf(`SELECT * 
						FROM  %s 
						WHERE id = $1`,
		songsTable)
	err := r.db.Get(&song, getByIdQuery, songId)

	return song, err
}

func (r *SongPostgres) Delete(songId int) error {
	deleteQuery := fmt.Sprintf(`DELETE FROM %s s
								USING %s l
								WHERE s.id = l.song_id
									AND s.id = $1`,
		songsTable, lyricsTable)
	_, err := r.db.Exec(deleteQuery, songId)

	return err
}

func (r *SongPostgres) Update(songId int, input model.UpdateSongInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Group != nil {
		setValues = append(setValues, fmt.Sprintf("group=$%d", argId))
		args = append(args, *input.Group)
		argId++
	}

	if input.SongName != nil {
		setValues = append(setValues, fmt.Sprintf("song_name=$%d", argId))
		args = append(args, *input.SongName)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s  
						SET %s 
						WHERE id = ul.list_id`,
		songsTable, setQuery)

	logrus.Debugf("updated query: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
