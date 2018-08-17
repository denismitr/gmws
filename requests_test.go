package gmws

import "io"

type fakeGetRequest struct {
	url        string
	httpMethod string
}

func (fgr fakeGetRequest) GenerateURL(host string) (string, error) {
	return "https://" + host + fgr.url, nil
}
func (fgr fakeGetRequest) Method() string {
	return fgr.httpMethod
}
func (fgr fakeGetRequest) Body() io.Reader {
	return nil
}
