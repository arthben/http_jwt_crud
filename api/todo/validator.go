package todo

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/arthben/http_jwt_crud/api/auth"
	"github.com/arthben/http_jwt_crud/internal/config"
	res "github.com/arthben/http_jwt_crud/internal/response"
	bindhttp "github.com/arthben/http_jwt_crud/pkg/bind_http"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const (
	ServiceCode = "24"
	TSLayout    = "2006-01-02T15:04:05+07:00"
)

func generateIDTodos() string {
	return uuid.New().String()
}

func validateClientKey(cfg *config.EnvParams, clientKey string) *res.Message {
	if cfg.Client.Key != clientKey {
		return res.BadResponse(http.StatusUnauthorized, ServiceCode, "00", "Unauthorized. Unknown Client Key")
	}

	return nil
}

func validateAccessToken(cfg *config.EnvParams, accessToken string) *res.Message {
	authService := auth.NewAuthService(cfg)
	bearerToken := strings.Split(accessToken, "Bearer ")
	return authService.ValidateAccessToken(strings.TrimSpace(bearerToken[1]), ServiceCode)
}

func validateSignature(header *RequestHeader, body []byte, reqPath string, method string, clientSecret string) *res.Message {
	// signature :
	// HTTP Method + ":" + request Path + ":" + access token + ":" + lowercase(hexencode(SHA256(minify(body)))) + ":" + timestamp

	dst := &bytes.Buffer{}
	if err := json.Compact(dst, body); err != nil {
		return res.BadResponse(http.StatusBadRequest, ServiceCode, "00", "Bad Request")
	}

	h := sha256.New()
	h.Write(dst.Bytes())

	strSign := strings.Join([]string{
		method,
		reqPath,
		header.Authorization[7:],
		strings.ToLower(hex.EncodeToString(h.Sum(nil))),
		header.Timestamp,
	}, ":")

	mac := hmac.New(sha512.New, []byte(clientSecret))
	mac.Write([]byte(strSign))

	internalSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	if internalSignature != header.Signature {
		return res.BadResponse(http.StatusUnauthorized, ServiceCode, "00", "Unauthorized Signature")
	}

	return nil
}

func validateBody(w http.ResponseWriter, r *http.Request) (*AddTodosRequest, *res.Message) {
	var payload AddTodosRequest

	if err := bindhttp.BindBody(w, r, &payload); err != nil {
		return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "00", "Bad Request")
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(&payload); err != nil {
		// check if the struct is nill
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "02", "Unauthorized. Bad Request")
		}

		for _, ve := range err.(validator.ValidationErrors) {
			switch ve.Field() {
			case "Title":
				if ve.Tag() == "required" {
					return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "02", "Missing Mandatory Field title")
				}
				if ve.Tag() == "min" || ve.Tag() == "max" {
					return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "01", "Invalid Field Format. Length of field value")
				}

			case "DetailTodo":
				if ve.Tag() == "max" {
					return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "01", "Invalid Field Format. Length of field value")
				}

			default:
				return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "00", "Unauthorized. Bad Request")

			}
		}
	}

	return &payload, nil
}

func validateHeader(r *http.Request) (*RequestHeader, *res.Message) {
	var header RequestHeader

	if err := bindhttp.BindHeader(r.Header, &header); err != nil {
		return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "00", "Bad Request")
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(&header); err != nil {
		// check if struct is nill
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "00", "Unauthorized. Bad Request")
		}

		for _, ve := range err.(validator.ValidationErrors) {
			switch ve.Field() {
			case "ContentType":
				if ve.Tag() == "required" {
					return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "02", "Missing Mandatory Field Content-Type")
				}

			case "Authorization":
				if ve.Tag() == "required" {
					return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "02", "Missing Mandatory Field Authorization")
				}

			case "ClientKey":
				if ve.Tag() == "required" {
					return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "02", "Missing Mandatory Field ClientKey")
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

	// validate Bearer
	authBearer := strings.Split(header.Authorization, "Bearer")
	if len(authBearer) == 0 {
		return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "01", "Invalid Field Format Authorization")
	}

	if len(strings.TrimSpace(authBearer[1])) == 0 {
		return nil, res.BadResponse(http.StatusBadRequest, ServiceCode, "01", "Invalid Field Format Authorization")
	}

	return &header, nil
}
