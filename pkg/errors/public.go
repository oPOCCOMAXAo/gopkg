package errors

import "github.com/gin-gonic/gin"

type PublicTextError struct {
	Message string
}

func NewPublicError(message string) error {
	return &PublicTextError{Message: message}
}

func (e *PublicTextError) Error() string {
	return e.Message
}

func NewGinPublicError(message string) *gin.Error {
	return &gin.Error{
		Err:  NewPublicError(message),
		Type: gin.ErrorTypePublic,
	}
}
