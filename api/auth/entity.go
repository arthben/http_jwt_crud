package auth

// time_format:"2006-01-02T15:04:05+07:00"
type TokenHeader struct {
	ContentType string `header:"Content-Type" validate:"required"`
	ClientKey   string `header:"X-Client-Key" validate:"required,max=64"`
	Timestamp   string `header:"X-Timestamp" validate:"required" `
	Signature   string `header:"X-Signature" validate:"required"`
}

type TokenRequest struct {
	GrantType string `json:"grantType" validate:"required"`
}

type TokenResponse struct {
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	AccessToken     string `json:"accessToken"`
	TokenType       string `json:"tokenType"`
	ExpiresIn       string `json:"expiresIn"`
}
