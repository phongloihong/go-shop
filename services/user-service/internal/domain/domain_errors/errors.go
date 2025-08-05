package domain_error

import "connectrpc.com/connect"

type DomainError interface {
	error
	Code() connect.Code
}

type domainError struct {
	message string
	code    connect.Code
}

func (e *domainError) Error() string {
	return e.message
}

func (e *domainError) Code() connect.Code {
	return e.code
}

func MapError(err error) *connect.Error {
	if domainErr, ok := err.(DomainError); ok {
		return connect.NewError(domainErr.Code(), domainErr)
	}

	return connect.NewError(connect.CodeInternal, err)
}

// Constructors
func NewNotFoundError(msg string) DomainError {
	return &domainError{
		message: msg,
		code:    connect.CodeNotFound,
	}
}

func NewAlreadyExistsError(msg string) DomainError {
	return &domainError{
		message: msg,
		code:    connect.CodeAlreadyExists,
	}
}

func NewUnauthorizedError(msg string) DomainError {
	return &domainError{
		message: msg,
		code:    connect.CodeUnauthenticated,
	}
}

func NewInvalidData(msg string) DomainError {
	return &domainError{
		message: msg,
		code:    connect.CodeInvalidArgument,
	}
}

func NewInternalError(msg string) DomainError {
	return &domainError{
		message: msg,
		code:    connect.CodeInternal,
	}
}
