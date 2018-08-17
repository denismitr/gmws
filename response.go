package gmws

// Response - Amazon MWS response body
type Response interface {
	Err() error
	Body() string
}

type ErrorResponse struct {
	err error
}

func (er ErrorResponse) Err() error {
	return er.err
}

func (ErrorResponse) Body() string {
	return ""
}

// NewErrorResponse - create new error response
func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{err: err}
}

// SuccessResponse contains a body of successfull response
// from Amazon API
type SuccessResponse struct {
	body string
}

func (SuccessResponse) Err() error {
	return nil
}

func (sr SuccessResponse) Body() string {
	return sr.body
}

// NewSuccessResponse - create new success response
func NewSuccessResponse(body string) *SuccessResponse {
	return &SuccessResponse{body: body}
}
