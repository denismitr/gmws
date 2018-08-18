package gmws

type fakeGetRequest struct {
	Action      string
	Method      string
	ActionPath  string
	QueryParams map[string]string
}

func (fgr fakeGetRequest) GetParams() *RequestParams {
	return &RequestParams{
		Method:      fgr.Method,
		Action:      fgr.Action,
		ActionPath:  fgr.ActionPath,
		QueryParams: fgr.QueryParams,
	}
}
