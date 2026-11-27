package apperror

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

const (
	CodeBadRequest         = 40000
	CodeUnauthorized       = 40100
	CodeForbidden          = 40300
	CodeNotFound           = 40400
	CodeConflict           = 40900
	CodeInternalServer     = 50000
	CodeBuiltinProtect     = 50010
	CodeDeleteRestricted   = 50011
	CodeDuplicateResource  = 40901
	CodeInvalidMenuRule    = 40010
	CodeInvalidMenuTree    = 40011
	CodeInvalidCredential  = 40101
	CodeInvalidCurrentUser = 40102
)
