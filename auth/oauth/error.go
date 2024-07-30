package oauth

// Error is used for generic OAuth errors.
type Error string

func (e Error) Error() string {
	return "oauth: " + string(e)
}

type ConvertCodeError string

func (e ConvertCodeError) Error() string {
	return "oauth converting access code to access token: " + string(e)
}

// Some specific errors

var (
	EmptyClientIDError    = Error("ClientID is empty")
	EmptyCodeError        = Error("Access Code/Token is empty")
	NilTokenError         = Error("Token is nil")
	TooManyCallbacksError = Error("Too many callback requests")
)
