package myerror

type NotFoundError struct {
}

func (e NotFoundError) Error() string {
	return "Not found error"
}

type InternalServerError struct {
}

func (e InternalServerError) Error() string {
	return "Internal server error"
}
