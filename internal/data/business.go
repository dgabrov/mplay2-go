package data

type Media struct {
	Id          string
	UserId      string
	Description string
	ContentType string
	Size        int64
	Width       int
	Height      int
}

type PlayList struct {
	PlaylistId  string
	UserId      string
	Description string
}

type DescriptionUpdate struct {
	Id          string
	Description string
}

type DeleteMediaRequest struct {
	Ids []string `json:"ids"`
}

type DeleteMediaResponse struct {
	Deleted int `json:"deleted"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

type RemoveMediaFromPlaylistRequest struct {
	PlaylistId string   `json:"playlistId"`
	Ids        []string `json:"ids"`
}
