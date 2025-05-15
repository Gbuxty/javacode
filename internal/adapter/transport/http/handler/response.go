package handler 

type errResponse struct {
	Error string `json:"error"`
}

type successResponse struct{
	Success bool `json:"success"`
}