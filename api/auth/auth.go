package auth

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/arthben/http_jwt_crud/internal/config"
	res "github.com/arthben/http_jwt_crud/internal/response"
	"github.com/golang-jwt/jwt"
)

type AuthService struct {
	cfg *config.EnvParams
}

func NewAuthService(cfg *config.EnvParams) *AuthService {
	return &AuthService{cfg: cfg}
}

func (a *AuthService) ValidateAccessToken(accessToken string, serviceCode string) *res.Message {
	pubKey := []byte(a.cfg.Server.PublicKey)
	if errCode := validateToken(a.cfg.Token.Issuer, accessToken, serviceCode, pubKey); errCode != nil {
		return errCode
	}

	return nil
}

func (a *AuthService) GetAccessToken(w http.ResponseWriter, r *http.Request) (*TokenResponse, *res.Message) {
	header, errHeader := validateHeader(r)
	if errHeader != nil {
		return nil, errHeader
	}

	_, errBody := validateBody(w, r)
	if errBody != nil {
		return nil, errBody
	}

	// validate client key
	if a.cfg.Client.Key != header.ClientKey {
		return nil, res.BadResponse(http.StatusUnauthorized, ServiceCode, "00", "Unathorized. Unknown Client")
	}

	pubKey := []byte(a.cfg.Client.PublicKey)
	strToSign := strings.Join([]string{a.cfg.Client.Key, header.Timestamp}, "|")
	if errSignature := validateSignature(header, strToSign, pubKey); errSignature != nil {
		return nil, errSignature
	}

	// generate access token
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(a.cfg.Server.PrivateKey))
	if err != nil {
		return nil, res.BadResponse(http.StatusUnauthorized, ServiceCode, "00", "Unauthorized. Signature")
	}

	claims := make(jwt.MapClaims)
	expTime, _ := strconv.Atoi(a.cfg.Token.Expire)
	now := time.Now()
	exp := time.Duration(expTime) * time.Second
	claims["iat"] = now.Unix()
	claims["iss"] = a.cfg.Token.Issuer
	claims["exp"] = now.Add(exp).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	strToken, err := token.SignedString(privateKey)
	if err != nil {
		return nil, res.BadResponse(http.StatusUnauthorized, ServiceCode, "00", "Unauthorized. Signature")
	}

	resp := &TokenResponse{
		ResponseCode:    "2007300",
		ResponseMessage: "Success",
		AccessToken:     strToken,
		TokenType:       "Bearer",
		ExpiresIn:       a.cfg.Token.Expire,
	}
	return resp, nil
}
