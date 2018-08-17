package gmws

import "io"

// Request to Amazon MWS
type MWSRequest interface {
	GenerateURL(host string) (string, error)
	Method() string
	Body() io.Reader
}
