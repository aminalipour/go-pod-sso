package errors

const (
	ErrNotFound                    = "resourceNotFound"
	ErrUnauthorized                = "unauthorizedAccess"
	ErrInvalidInput                = "invalidInput"
	ErrTimeout                     = "operationTimedOut"
	ErrInternalServer              = "internalServerError"
	ErrConflict                    = "conflictOccurred"
	ErrAlreadyExists               = "resourceAlreadyExists"
	ErrNotImplemented              = "featureNotImplemented"
	ErrTooManyRequests             = "tooManyRequests"
	ErrServiceUnavailable          = "serviceUnavailable"
	ErrSignatureKeyOrFileIsMissing = "signatureKeyOrFileIsMissing"
)

const (
	Signature = "signiture file invalid or not found"
)
