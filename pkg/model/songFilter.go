package model

type SongFilter struct {
	Group  *string `json:"group" form:"group"`
	Song   *string `json:"song" form:"song"`
	Limit  int     `json:"limit" form:"limit"`
	Offset int     `json:"offset" form:"offset"`
}
