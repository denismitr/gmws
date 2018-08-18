package gmws

// MWSRequest - request to Amazon MWS API
type MWSRequest interface {
	GetParams() *RequestParams
}

// RequestParams - all parameters of particular request
type RequestParams struct {
	QueryParams map[string]string
	Action      string
	ActionPath  string
	Method      string
}
