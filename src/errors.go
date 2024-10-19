package src

import "net/http"

type Err struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Obj     error
}

func (e Err) Error() string {
	return e.Message
}

var (
	ErrAuthorize = Err{
		Message: "unauthorize error",
		Status:  http.StatusUnauthorized,
	}
	ErrValidation = Err{
		Message: "validation error",
		Status:  http.StatusBadRequest,
	}
	ErrNotFound = Err{
		Message: "resource not found",
		Status:  http.StatusNotFound,
	}
	ErrInvalidWithdrawAmount = Err{
		Message: "insufficient balance",
		Status:  http.StatusUnprocessableEntity,
	}
	ErrInvalidQty = Err{
		Message: "invalid qty",
		Status:  http.StatusUnprocessableEntity,
	}
	ErrBadRequest = Err{
		Message: "bad request",
		Status:  http.StatusBadRequest,
	}
)
