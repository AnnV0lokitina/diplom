package handler

type JSONUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
