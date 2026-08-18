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

type User struct {
	UserID         string `json:"user_id"`
	ProvidedUserID string `json:"provided_user_id"`
	Login          string `json:"login"`
	Name           string `json:"name"`
}

type Session struct {
	SessionID  string `json:"session_id"`
	UserID     string `json:"user_id"`
	Token      string `json:"token"`
	ExpiredInd string `json:"expired_ind"`
	ExpiryDt   string `json:"expiry_dt"`
}
