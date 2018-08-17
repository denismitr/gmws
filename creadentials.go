package gmws

// Credentials - Amazon MWS credentials
type Credentials struct {
	AccessKey string
	SecretKey string
	AuthToken string
	SellerID  string
}

// NewCredentials creates a new set of Amazon MWS credentials
func NewCredentials(accessKey, secretKey, authToken, sellerID string) Credentials {
	return Credentials{
		AccessKey: accessKey,
		SecretKey: secretKey,
		AuthToken: authToken,
		SellerID:  sellerID,
	}
}
