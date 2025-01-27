package model

type SongFilter struct {
	Group  string `json:"group"`
	Song   string `json:"song"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}
