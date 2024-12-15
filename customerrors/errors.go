package customerrors

import "net/http"

type APIError struct {
	Status  int    `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}

var (
	Success = APIError{Status: http.StatusOK, Code: 0, Message: "Success"}

	// Common errors - 5000 - 5999
	ErrInvalidUUID = APIError{Status: http.StatusBadRequest, Code: 5000, Message: "Invalid UUID"}

	// User errors - 1000 - 1999
	ErrUserExists      = APIError{Status: http.StatusBadRequest, Code: 1001, Message: "Username already exists"}
	ErrEmailExists     = APIError{Status: http.StatusBadRequest, Code: 1002, Message: "Email already exists"}
	ErrHashPassword    = APIError{Status: http.StatusInternalServerError, Code: 1003, Message: "Error hashing password"}
	ErrCreateUser      = APIError{Status: http.StatusInternalServerError, Code: 1004, Message: "Error creating user"}
	ErrUsernameInvalid = APIError{Status: http.StatusBadRequest, Code: 1005, Message: "Invalid username"}
	ErrEmailInvalid    = APIError{Status: http.StatusBadRequest, Code: 1006, Message: "Invalid email"}
	ErrListUsers       = APIError{Status: http.StatusInternalServerError, Code: 1007, Message: "Error listing users"}
	ErrUserNotFound    = APIError{Status: http.StatusNotFound, Code: 1008, Message: "User not found"}
	ErrDeleteUser      = APIError{Status: http.StatusInternalServerError, Code: 1009, Message: "Error deleting user"}
	ErrUpdateUser      = APIError{Status: http.StatusInternalServerError, Code: 1010, Message: "Error updating user"}

	// Share Errors - 2000 - 2999
	ErrGettingShare = APIError{Status: http.StatusInternalServerError, Code: 2000, Message: "Error while getting a share from DB"}
	ErrShareExists  = APIError{Status: http.StatusConflict, Code: 2001, Message: "This share already has been shared."}
	ErrInvalidInput = APIError{Status: http.StatusBadRequest, Code: 2002, Message: "Invalid input. A share needs shared_by, shared_with and valid_until"}
	ErrDeleteShare  = APIError{Status: http.StatusInternalServerError, Code: 2003, Message: "Error deleting share"}

)
