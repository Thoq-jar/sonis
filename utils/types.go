package utils

type Track struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Track  int    `json:"track"`
	Album  string `json:"album"`
	Artist string `json:"artist"`
	Path   string
}

type Album struct {
	Name   string   `json:"name"`
	Artist string   `json:"artist"`
	Tracks []*Track `json:"tracks"`
}

type Artist struct {
	Name   string   `json:"name"`
	Albums []*Album `json:"albums"`
}

type Library struct {
	Artists []*Artist `json:"artists"`
}
