package api

type Error struct {
	Code int `json:"code"`
	Err string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err: err,
	}
}