package errors

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var (
	Conflict = &Error{
		Code:    "CONFLICT",
		Message: "Provided slug conflicts with another link.",
	}
	Expired = &Error{
		Code:    "EXPIRED",
		Message: "This link expired.",
	}
	InternalServerError = &Error{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "An unknown internal server error occurred while processing your request. Please try again later.",
	}
	InvalidJSONBody = &Error{
		Code:    "INVALID_JSON_BODY",
		Message: "Your request body contains invalid JSON.",
	}
	NotFound = &Error{
		Code:    "NOT_FOUND",
		Message: "Link with provided slug not found.",
	}
)
