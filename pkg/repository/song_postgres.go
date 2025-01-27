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

	createSongQuery := fmt.Sprintf(`INSERT INTO %s ("group", song_name, release_date, link)  
									VALUES ($1, $2, NULL, NULL) 
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

func (r *SongPostgres) GetAllSongsWithFilter(filter model.SongFilter) ([]model.Song, error) {
	var songs []model.Song
	queryParams := []interface{}{filter.Limit, filter.Offset}
	filterSql := ""
	argId := 3

	if filter.Group != nil {
		fmt.Printf("\ng\n")
		filterSql += fmt.Sprintf(" AND \"group\" ILIKE $%d", argId)
		argId++
		queryParams = append(queryParams, *filter.Group)
	}

	if filter.Song != nil {
		fmt.Printf("\ns\n")
		filterSql += fmt.Sprintf(" AND song_name ILIKE $%d", argId)
		argId++
		queryParams = append(queryParams, *filter.Song)
	}

	getAllSongsWithFilterQuery := fmt.Sprintf(`SELECT *
												FROM %s
												WHERE (1=1) %s
												LIMIT $1 
												OFFSET $2`,
		songsTable, filterSql)

	logrus.Infof("Final Query: %s", getAllSongsWithFilterQuery)
	logrus.Infof("Final Params: %v", queryParams)

	err := r.db.Select(&songs, getAllSongsWithFilterQuery, queryParams...)

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

	query := fmt.Sprintf(`UPDATE %s  
						SET %s 
						WHERE id = $%d`,
		songsTable, setQuery, argId)
	args = append(args, songId)

	logrus.Debugf("updated query: %s", query)
	logrus.Debugf("args: %v", args)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *SongPostgres) UpdateSongWithAPIInfo(updateSongApiData model.UpdateSongApiData) error {
	tx, err := r.db.Begin()
	if err != nil {
		return nil
	}

	updateSongQuery := fmt.Sprintf(`UPDATE %s  
							SET release_date = $1, link = $2  
							WHERE id = $3`,
		songsTable)
	_, err = tx.Exec(updateSongQuery, updateSongApiData.ReleaseDate, updateSongApiData.Link, updateSongApiData.SongId)
	if err != nil {
		tx.Rollback()
		return err
	}

	deleteLyricsQuery := fmt.Sprintf(`DELETE 
									FROM %s 
									WHERE song_id = $1`,
		lyricsTable)
	_, err = tx.Exec(deleteLyricsQuery, updateSongApiData.SongId)
	if err != nil {
		tx.Rollback()
		return err
	}

	createLyricsQuery := fmt.Sprintf(`INSERT INTO %s (song_id, verse_number, text) 
									VALUES ($1, $2, $3)`,
		lyricsTable)
	for i, verse := range updateSongApiData.Lyrics {
		_, err = tx.Exec(createLyricsQuery, updateSongApiData.SongId, i+1, verse)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
