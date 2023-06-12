package response

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  *Error      `json:"error"`
}

type Error struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details []ErrorDetails `json:"details"`
}

type ErrorDetails struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e Response) Error() string {
	return "Not found error"
}
