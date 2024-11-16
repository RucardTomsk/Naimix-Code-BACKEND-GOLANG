package api

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"net/http"
)

func ResponseFromServiceError(serviceError base.ServiceError) base.ResponseFailure {
	return base.ResponseFailure{
		Status:  http.StatusText(serviceError.Code),
		Blame:   serviceError.Blame,
		Message: serviceError.Message,
	}
}

func GeneralParsingError() base.ResponseFailure {
	return base.ResponseFailure{
		Status:  http.StatusText(http.StatusBadRequest),
		Blame:   base.BlameUser,
		Message: "failed to parse request parameters",
	}
}

func GeneralSortError() base.ResponseFailure {
	return base.ResponseFailure{
		Status:  http.StatusText(http.StatusBadRequest),
		Blame:   base.BlameUser,
		Message: "bad sort parameters",
	}
}

func GeneralPaginationError() base.ResponseFailure {
	return base.ResponseFailure{
		Status:  http.StatusText(http.StatusBadRequest),
		Blame:   base.BlameUser,
		Message: "bad pagination parameters",
	}
}

func GeneralFilterError() base.ResponseFailure {
	return base.ResponseFailure{
		Status:  http.StatusText(http.StatusBadRequest),
		Blame:   base.BlameUser,
		Message: "bad filter parameters",
	}
}
func GeneralUnexpectedError() base.ResponseFailure {
	return base.ResponseFailure{
		Status:  http.StatusText(http.StatusInternalServerError),
		Blame:   base.BlameUnknown,
		Message: "internal error",
	}
}

func ResponseUnauthorizedError() base.ResponseFailure {
	return base.ResponseFailure{
		Status:  http.StatusText(http.StatusUnauthorized),
		Blame:   base.BlameUnknown,
		Message: "unauthorized",
	}
}
