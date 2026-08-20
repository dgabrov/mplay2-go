package data

type Media struct {
	Id          string `json:"id"`
	UserId      string `json:"userId"`
	Description string `json:"description"`
	ContentType string `json:"contentType"`
	Size        int64  `json:"size"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
}

type PlayList struct {
	PlaylistId  string `json:"playlistId"`
	UserId      string `json:"userId"`
	Description string `json:"description"`
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
