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
	logrus.Debug("NewSongPostgres - initializing SongPostgres repository")
	return &SongPostgres{db: db}
}

func (r *SongPostgres) Create(song model.Song) (int, error) {
	logrus.Debugf("CreateSong - Group: %s, SongName: %s", song.Group, song.SongName)

	var id int
	createSongQuery := fmt.Sprintf(`INSERT INTO %s ("group", song_name, release_date, link)  
									VALUES ($1, $2, NULL, NULL) 
									RETURNING id`, songsTable)
	err := r.db.QueryRow(createSongQuery, song.Group, song.SongName).Scan(&id)
	if err != nil {
		logrus.Errorf("Error while creating song: %s", err)
		return 0, err
	}

	return id, nil
}

func (r *SongPostgres) GetLyrics(songId int, limit, offset int) ([]model.Lyrics, error) {
	logrus.Debugf("GetLyrics - SongID: %d, Limit: %d, Offset: %d", songId, limit, offset)

	var lyrics []model.Lyrics
	getLyricsQuery := fmt.Sprintf(`SELECT * FROM %s WHERE song_id = $1 ORDER BY verse_number LIMIT $2 OFFSET $3`, lyricsTable)
	err := r.db.Select(&lyrics, getLyricsQuery, songId, limit, offset)
	if err != nil {
		logrus.Errorf("Error fetching lyrics for song %d: %s", songId, err)
	}
	return lyrics, err
}

func (r *SongPostgres) GetAllSongsWithFilter(filter model.SongFilter) ([]model.Song, error) {
	logrus.Debugf("GetAllSongsWithFilter - Filter: %+v", filter)

	var songs []model.Song
	queryParams := []interface{}{filter.Limit, filter.Offset}
	filterSql := ""
	argId := 3

	if filter.Group != nil {
		filterSql += fmt.Sprintf(" AND \"group\" ILIKE $%d", argId)
		argId++
		queryParams = append(queryParams, *filter.Group)
	}

	if filter.Song != nil {
		filterSql += fmt.Sprintf(" AND song_name ILIKE $%d", argId)
		argId++
		queryParams = append(queryParams, *filter.Song)
	}

	getAllSongsWithFilterQuery := fmt.Sprintf(`SELECT * FROM %s WHERE (1=1) %s LIMIT $1 OFFSET $2`, songsTable, filterSql)

	err := r.db.Select(&songs, getAllSongsWithFilterQuery, queryParams...)
	if err != nil {
		logrus.Errorf("Error fetching songs with filter: %v", err)
	}

	return songs, err
}

func (r *SongPostgres) GetById(songId int) (model.Song, error) {
	logrus.Debugf("GetById - SongID: %d", songId)

	var song model.Song
	getByIdQuery := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, songsTable)
	err := r.db.Get(&song, getByIdQuery, songId)

	if err != nil {
		logrus.Errorf("Error fetching song with ID %d: %s", songId, err)
	}
	return song, err
}

func (r *SongPostgres) Delete(songId int) error {
	logrus.Debugf("DeleteSong - SongID: %d", songId)

	deleteQuery := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, songsTable)
	_, err := r.db.Exec(deleteQuery, songId)
	if err != nil {
		logrus.Errorf("Error deleting song with ID %d: %s", songId, err)
	}

	return err
}

func (r *SongPostgres) Update(songId int, input model.UpdateSongInput) error {
	logrus.Debugf("UpdateSong - SongID: %d, Input: %+v", songId, input)

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Group != nil {
		setValues = append(setValues, fmt.Sprintf(`"group"=$%d`, argId))
		args = append(args, *input.Group)
		argId++
	}

	if input.SongName != nil {
		setValues = append(setValues, fmt.Sprintf("song_name=$%d", argId))
		args = append(args, *input.SongName)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $%d`, songsTable, setQuery, argId)
	args = append(args, songId)

	logrus.Debugf("Update query: %s", query)
	logrus.Debugf("Args: %v", args)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		logrus.Errorf("Error updating song with ID %d: %s", songId, err)
	}
	return err
}

func (r *SongPostgres) UpdateSongWithAPIInfo(updateSongApiData model.UpdateSongApiData) error {
	logrus.Debugf("UpdateSongWithAPIInfo - SongID: %d", updateSongApiData.SongId)

	tx, err := r.db.Begin()
	if err != nil {
		logrus.Errorf("Error starting transaction: %s", err)
		return err
	}

	updateSongQuery := fmt.Sprintf(`UPDATE %s SET release_date = $1, link = $2 WHERE id = $3`, songsTable)
	_, err = tx.Exec(updateSongQuery, updateSongApiData.ReleaseDate, updateSongApiData.Link, updateSongApiData.SongId)
	if err != nil {
		tx.Rollback()
		logrus.Errorf("Error updating song with ID %d: %s", updateSongApiData.SongId, err)
		return err
	}

	deleteLyricsQuery := fmt.Sprintf(`DELETE FROM %s WHERE song_id = $1`, lyricsTable)
	_, err = tx.Exec(deleteLyricsQuery, updateSongApiData.SongId)
	if err != nil {
		tx.Rollback()
		logrus.Errorf("Error deleting lyrics for song with ID %d: %s", updateSongApiData.SongId, err)
		return err
	}

	createLyricsQuery := fmt.Sprintf(`INSERT INTO %s (song_id, verse_number, text) VALUES ($1, $2, $3)`, lyricsTable)
	for i, verse := range updateSongApiData.Lyrics {
		_, err = tx.Exec(createLyricsQuery, updateSongApiData.SongId, i+1, verse)
		if err != nil {
			tx.Rollback()
			logrus.Errorf("Error inserting lyrics for song with ID %d: %s", updateSongApiData.SongId, err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		logrus.Errorf("Error committing transaction for song with ID %d: %s", updateSongApiData.SongId, err)
	}

	return err
}
