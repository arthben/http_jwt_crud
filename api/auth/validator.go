package auth

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/http"
	"time"

	res "github.com/arthben/http_jwt_crud/internal/response"
	bindhttp "github.com/arthben/http_jwt_crud/pkg/bind_http"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
)

const (
	ServiceCode = "73"
	GrantType   = "client_credentials"
	TSLayout    = "2006-01-02T15:04:05+07:00"
)

func validateToken(issuer string, accessToken string, serviceCode string, pubKey []byte) *res.Message {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKey)
	if err != nil {
		return res.BadResponse(http.StatusUnauthorized, serviceCode, "01", "Unauthorized. Invalid Token")
	}

	token, err := jwt.Parse(accessToken, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method %s", jwtToken.Header["alg"])
		}
		return publicKey, nil
	})

	// if error occured, token has been expired
	if err != nil {
		return res.BadResponse(http.StatusUnauthorized, serviceCode, "01", "Unauthorized. Invalid Token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return res.BadResponse(http.StatusUnauthorized, serviceCode, "01", "Unauthorized. Invalid Token")
	}

	if claims["iss"] != issuer {
		return res.BadResponse(http.StatusUnauthorized, serviceCode, "01", "Unauthorized. Invalid Token")
	}

	return nil
}

func validateSignature(header *TokenHeader, strToSign string, pubKey []byte) *res.Message {
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return res.BadResponse(http.StatusUnauthorized, ServiceCode, "00", "Unauthorized. Signature")
	}

	parseResult, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return res.BadResponse(http.StatusUnauthorized, ServiceCode, "00", "Unauthorized. Signature")
	}
	publicKey := parseResult.(*rsa.PublicKey)

	msgHash := sha256.New()
	msgHash.Write([]byte(strToSign))
	msgHashSum := msgHash.Sum(nil)

	signature, err := base64.StdEncoding.DecodeString(header.Signature)
	if err != nil {
		return res.BadResponse(http.StatusUnauthorized, ServiceCode, "00", "Unauthorized. Signature")
	}

	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, msgHashSum, signature); err != nil {
		return res.BadResponse(http.StatusUnauthorized, ServiceCode, "00", "Unauthorized. Signature")
	}
	return nil
}

func validateBody(w http.ResponseWriter, r *http.Request) (*TokenRequest, *res.Message) {
	var payload TokenRequest

	if err := bindhttp.BindBody(w, r, &payload); err != nil {
		return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "00", "Bad Request")
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(&payload); err != nil {
		// check if the struct is nill
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "02", "Invalid Mandatory Field grantType")
		}
		return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "00", "Bad Request")
	}

	if payload.GrantType != GrantType {
		return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "00", "Unsupported grantType")
	}

	return &payload, nil
}

func validateHeader(r *http.Request) (*TokenHeader, *res.Message) {
	var header TokenHeader

	if err := bindhttp.BindHeader(r.Header, &header); err != nil {
		return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "00", "Unauthorized. Bad Request")
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(&header); err != nil {
		// check if the struct is nill
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "00", "Unauthorized. Bad Request")
		}

		for _, ve := range err.(validator.ValidationErrors) {
			switch ve.Field() {
			case "ContentType":
				if ve.Tag() == "required" {
					return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "02", "Missing Mandatory Field Content-Type")
				}

			case "ClientKey":
				if ve.Tag() == "required" {
					return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "02", "Missing Mandatory Field X-CLIENT-KEY")
				}
				if ve.Tag() == "max" {
					return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "01", "Invalid Field Format X-CLIENT-KEY")
				}

			case "Timestamp":
				if ve.Tag() == "required" {
					return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "02", "Missing Mandatory Field X-TIMESTAMP")
				}

			case "Signature":
				if ve.Tag() == "required" {
					return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "02", "Missing Mandatory Field X-SIGNATURE")
				}

			default:
				return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "00", "Unauthorized. Bad Request")
			}
		}
	}

	// validate content-type
	if header.ContentType != "application/json" {
		return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "01", "Invalid Field Format Content-Type")
	}

	// validate timestamp format
	if _, err := time.Parse(TSLayout, header.Timestamp); err != nil {
		return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "01", "Invalid Field Format X-TIMESTAMP")
	}

	return &header, nil
}
