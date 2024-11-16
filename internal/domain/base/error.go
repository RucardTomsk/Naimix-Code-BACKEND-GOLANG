package base

import (
	"fmt"
	"net/http"
)

// ServiceError is a general optional error that can be
// returned by any type of service. NOT SERIALIZABLE.
type ServiceError struct {
	Message string
	Blame   Blame
	Code    int
	Err     error
}

// NewPostgresWriteError returns ServiceError with general write error message.
func NewPostgresWriteError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlamePostgres,
		Code:    http.StatusInternalServerError,
		Message: "failed to write data to database",
	}
}

func NewSendMessageError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlameMail,
		Code:    http.StatusInternalServerError,
		Message: "failed send message",
	}
}

// NewPostgresReadError returns ServiceError with general read error message.
func NewPostgresReadError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlamePostgres,
		Code:    http.StatusInternalServerError,
		Message: "failed to read data from database",
	}
}

func NewPostgresDuplicatedKeyError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlamePostgres,
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
}

// NewNotFoundError returns ServiceError with general not found error message.
func NewNotFoundError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlameUser,
		Code:    http.StatusNotFound,
		Message: "not found",
	}
}

func NewUnauthorizedError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlameUser,
		Code:    http.StatusUnauthorized,
		Message: "unauthorized user",
	}
}

func NewParseEnumError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlameUser,
		Code:    http.StatusInternalServerError,
		Message: "failed to parse enum",
	}
}

func NewHttpServerConnectError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlameServer,
		Code:    http.StatusInternalServerError,
		Message: "failed connect http server",
	}
}

func NewHttpServerRequestError(err error, code int) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlameServer,
		Code:    code,
		Message: "failed to request",
	}
}

func NewPathError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlameServer,
		Code:    http.StatusInternalServerError,
		Message: "failed path request",
	}
}

func NewReadByteError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlameServer,
		Code:    http.StatusInternalServerError,
		Message: "failed read byte",
	}
}

func NewJsonUnmarshalError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlameServer,
		Code:    http.StatusInternalServerError,
		Message: "failed unmarshal json",
	}
}

func NewJsonMarshalError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlameServer,
		Code:    http.StatusInternalServerError,
		Message: "failed marshal json",
	}
}

func NewLoginError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlameUser,
		Code:    http.StatusUnauthorized,
		Message: "failed checks email or password",
	}
}

func NewNotSessionError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlameUser,
		Code:    http.StatusBadRequest,
		Message: "not found session",
	}
}

func NewExpiredDate(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlameUser,
		Code:    http.StatusUnauthorized,
		Message: "refresh token expired",
	}
}

func NewCreateJWTError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlameServer,
		Code:    http.StatusInternalServerError,
		Message: "failed build jwt",
	}
}
func (e *ServiceError) Error() string {
	return fmt.Sprintf("[%d] %v (blame: %s)", e.Code, e.Err, e.Blame)
}
