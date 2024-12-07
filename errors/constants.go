package errors

const (
	ErrNotFound                    = "resource_not_found"
	ErrUnauthorized                = "unauthorized_access"
	ErrInvalidInput                = "invalid_input"
	ErrTimeout                     = "operation_timed_out"
	ErrInternalServer              = "internal_server_error"
	ErrConflict                    = "conflict_occurred"
	ErrAlreadyExists               = "resource_already_exists"
	ErrNotImplemented              = "feature_not_implemented"
	ErrTooManyRequests             = "too_many_requests"
	ErrServiceUnavailable          = "service_unavailable"
	ErrSignatureKeyOrFileIsMissing = "signature_key_or_file_is_missing"
)

const (
	Signature = "signiture file invalid or not found"
)
