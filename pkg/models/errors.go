package models

const (
	ErrCodeFileNotFound      = "FILE_NOT_FOUND"
	ErrCodeHashMismatch      = "HASH_MISMATCH"
	ErrCodeInvalidPath       = "INVALID_PATH"
	ErrCodeReadFailed        = "READ_FAILED"
	ErrCodeWriteFailed       = "WRITE_FAILED"
	ErrCodeEditFailed        = "EDIT_FAILED"
	ErrCodeExecFailed        = "EXEC_FAILED"
	ErrCodeProcessNotFound   = "PROCESS_NOT_FOUND"
	ErrCodeInvalidCategory   = "INVALID_CATEGORY"
	ErrCodeQueryFailed       = "QUERY_FAILED"
	ErrCodeSearchFailed      = "SEARCH_FAILED"
	ErrCodeMutationNotAllowed = "MUTATION_NOT_ALLOWED"
	ErrCodeInvalidRequest    = "INVALID_REQUEST"
	ErrCodeInternalError     = "INTERNAL_ERROR"
)

func NewErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
	}
}
