package httperrors

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	MsgUserAlreadyExists  = "User already exists"
	ErrBadRequest         = "Bad request"
	ErrEmailAlreadyExists = "User with given email already exists"
	ErrNoSuchUser         = "User not found"
	ErrWrongCredentials   = "Wrong Credentials"
	ErrNotFound           = "Not Found"
	ErrUnauthorized       = "Unauthorized"
	ErrForbidden          = "Forbidden"
	ErrBadQueryParams     = "Invalid query params"
)

var (
	BadRequest            = errors.New("bad request")
	WrongCredentials      = errors.New("wrong Credentials")
	NotFound              = errors.New("not Found")
	Unauthorized          = errors.New("unauthorized")
	Forbidden             = errors.New("forbidden")
	PermissionDenied      = errors.New("permission Denied")
	ExpiredCSRFError      = errors.New("expired CSRF token")
	WrongCSRFToken        = errors.New("wrong CSRF token")
	CSRFNotPresented      = errors.New("CSRF not presented")
	NotRequiredFields     = errors.New("no such required fields")
	BadQueryParams        = errors.New("invalid query params")
	InternalServerError   = errors.New("internal Server Error")
	RequestTimeoutError   = errors.New("request Timeout")
	ExistsEmailError      = errors.New("user with given email already exists")
	InvalidJWTToken       = errors.New("invalid JWT token")
	InvalidJWTClaims      = errors.New("invalid JWT claims")
	NotAllowedImageHeader = errors.New("not allowed image header")
	NoCookie              = errors.New("not found cookie header")
)

// Rest error interface
type RestErr interface {
	Status() int
	Error() string
	Causes() interface{}
}

// Rest error struct
type RestError struct {
	ErrStatus int         `json:"status,omitempty"`
	ErrError  string      `json:"error,omitempty"`
	ErrCauses interface{} `json:"-"`
}

// Error  Error() interface method
func (e RestError) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", e.ErrStatus, e.ErrError, e.ErrCauses)
}

// Error status
func (e RestError) Status() int {
	return e.ErrStatus
}

// RestError Causes
func (e RestError) Causes() interface{} {
	return e.ErrCauses
}

// New Rest Error
func NewRestError(status int, err string, causes interface{}) RestErr {
	return RestError{
		ErrStatus: status,
		ErrError:  err,
		ErrCauses: causes,
	}
}

// New Rest Error With Message
func NewRestErrorWithMessage(status int, err string, causes interface{}) RestErr {
	return RestError{
		ErrStatus: status,
		ErrError:  err,
		ErrCauses: causes,
	}
}

// New Rest Error From Bytes
func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var apiErr RestError
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return apiErr, nil
}

// New Bad Request Error
func NewBadRequestError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusBadRequest,
		ErrError:  BadRequest.Error(),
		ErrCauses: causes,
	}
}

// New Not Found Error
func NewNotFoundError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusNotFound,
		ErrError:  NotFound.Error(),
		ErrCauses: causes,
	}
}

// New Unauthorized Error
func NewUnauthorizedError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusUnauthorized,
		ErrError:  Unauthorized.Error(),
		ErrCauses: causes,
	}
}

// New Forbidden Error
func NewForbiddenError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusForbidden,
		ErrError:  Forbidden.Error(),
		ErrCauses: causes,
	}
}

// New Internal Server Error
func NewInternalServerError(causes interface{}) RestErr {
	result := RestError{
		ErrStatus: http.StatusInternalServerError,
		ErrError:  InternalServerError.Error(),
		ErrCauses: causes,
	}
	return result
}

func NewBadQueryParamsError(causes interface{}) RestErr {
	result := RestError{
		ErrStatus: http.StatusUnprocessableEntity,
		ErrError:  BadQueryParams.Error(),
		ErrCauses: causes,
	}
	return result
}

// Parser of error string messages returns RestError
func ParseErrors(err error) RestErr {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return NewRestError(http.StatusNotFound, NotFound.Error(), err)
	case errors.Is(err, Forbidden):
		return NewRestError(http.StatusForbidden, Forbidden.Error(), err)
	case errors.Is(err, BadQueryParams):
		return NewRestError(http.StatusUnprocessableEntity, BadQueryParams.Error(), err)
	case errors.Is(err, context.DeadlineExceeded):
		return NewRestError(http.StatusRequestTimeout, RequestTimeoutError.Error(), err)
	case strings.Contains(err.Error(), "SQLSTATE"):
		return parseSqlErrors(err)
	case strings.Contains(err.Error(), "Field validation"):
		return parseValidatorError(err)
	case errors.Is(err, BadRequest) || strings.Contains(err.Error(), "Unmarshal"):
		return NewRestError(http.StatusBadRequest, BadRequest.Error(), err)
	case strings.Contains(err.Error(), "UUID"):
		return NewRestError(http.StatusBadRequest, err.Error(), err)
	case strings.Contains(strings.ToLower(err.Error()), "password"):
		return NewRestError(http.StatusUnauthorized, Unauthorized.Error(), err)
	case strings.Contains(strings.ToLower(err.Error()), "token"):
		return NewRestError(http.StatusUnauthorized, Unauthorized.Error(), err)
	case strings.Contains(strings.ToLower(err.Error()), "bcrypt"):
		return NewRestError(http.StatusBadRequest, BadRequest.Error(), err)
	default:
		if restErr, ok := err.(RestErr); ok {
			return restErr
		}
		return NewInternalServerError(err)
	}
}

func parseSqlErrors(err error) RestErr {
	if strings.Contains(err.Error(), "23505") {
		return NewRestError(http.StatusBadRequest, ExistsEmailError.Error(), err)
	}

	return NewRestError(http.StatusBadRequest, BadRequest.Error(), err)
}

func parseValidatorError(err error) RestErr {
	if strings.Contains(err.Error(), "Password") {
		return NewRestError(http.StatusBadRequest, "Invalid password, min length 6", err)
	}

	if strings.Contains(err.Error(), "Email") {
		return NewRestError(http.StatusBadRequest, "Invalid email", err)
	}

	return NewRestError(http.StatusBadRequest, BadRequest.Error(), err)
}

// Error response
func ErrorResponse(err error) (int, interface{}) {
	return ParseErrors(err).Status(), ParseErrors(err)
}
