package domain

import "fmt"

// ErrorCode представляет код ошибки
type ErrorCode string

const (
	ErrorCodeTeamExists   ErrorCode = "TEAM_EXISTS"
	ErrorCodePRExists     ErrorCode = "PR_EXISTS"
	ErrorCodePRMerged     ErrorCode = "PR_MERGED"
	ErrorCodeNotAssigned  ErrorCode = "NOT_ASSIGNED"
	ErrorCodeNoCandidate  ErrorCode = "NO_CANDIDATE"
	ErrorCodeNotFound     ErrorCode = "NOT_FOUND"
)

// DomainError представляет доменную ошибку
type DomainError struct {
	Code    ErrorCode
	Message string
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewDomainError создает новую доменную ошибку
func NewDomainError(code ErrorCode, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

