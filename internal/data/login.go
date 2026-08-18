package data

type LoginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Id     string   `json:"id"`
	Login  string   `json:"login"`
	Name   string   `json:"name"`
	Rights []string `json:"rights"`
}
