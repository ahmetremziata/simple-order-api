package response

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type ErrorBuilder interface {
	SetError(statusCode int, message string) ErrorBuilder
	Build() ErrorResponse
}

type ErrorBuilderImp struct {
	errorResponse ErrorResponse
}

func NewErrorBuilder() ErrorBuilder {
	response := ErrorResponse{}
	return &ErrorBuilderImp{
		errorResponse: response,
	}
}

func (builder *ErrorBuilderImp) SetError(statusCode int, message string) ErrorBuilder {
	builder.errorResponse.Message = message
	builder.errorResponse.StatusCode = statusCode
	return builder
}

func (builder *ErrorBuilderImp) Build() ErrorResponse {
	return builder.errorResponse
}
