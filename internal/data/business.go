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
