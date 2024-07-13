package tests

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/arthben/http_jwt_crud/api/auth"
	"github.com/arthben/http_jwt_crud/api/handlers"
	"github.com/arthben/http_jwt_crud/api/todo"
	"github.com/arthben/http_jwt_crud/internal/config"
	"github.com/arthben/http_jwt_crud/internal/database"
	"github.com/arthben/http_jwt_crud/internal/response"
	"github.com/jmoiron/sqlx"
)

const (
	chars    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	TSLayout = "2006-01-02T15:04:05+07:00"
)

var (
	db  *sqlx.DB
	cfg *config.EnvParams
)

func TestGetToken(t *testing.T) {
	now := time.Now()
	endpoint := "/v1.0/access-token"
	payload := &auth.TokenRequest{
		GrantType: "client_credentials",
	}

	scenario := []struct {
		name           string
		payload        *auth.TokenRequest
		header         *auth.TokenHeader
		expectRespCode string
		statusCode     int
	}{
		{
			name:    "Success get token Access",
			payload: payload,
			header: &auth.TokenHeader{
				ContentType: "application/json",
				ClientKey:   cfg.Client.Key,
				Timestamp:   now.Format(TSLayout),
				Signature:   "", // will be generate
			},
			expectRespCode: "2007300",
			statusCode:     http.StatusOK,
		},
		{
			name:    "Failed. Empty Body",
			payload: &auth.TokenRequest{},
			header: &auth.TokenHeader{
				ContentType: "application/json",
				ClientKey:   cfg.Client.Key,
				Timestamp:   now.Format(TSLayout),
				Signature:   "", // will be generate
			},
			expectRespCode: "4007300",
			statusCode:     400,
		},
		{
			name:    "Failed. Empty Client Key",
			payload: payload,
			header: &auth.TokenHeader{
				ContentType: "application/json",
				ClientKey:   "",
				Timestamp:   now.Format(TSLayout),
				Signature:   "",
			},
			expectRespCode: "4007302",
			statusCode:     400,
		},
		{
			name:    "Failed. Wrong Client Key",
			payload: payload,
			header: &auth.TokenHeader{
				ContentType: "application/json",
				ClientKey:   "ahjksdkjashdkjahskjdhaskjdha",
				Timestamp:   now.Format(TSLayout),
				Signature:   "",
			},
			expectRespCode: "4017300",
			statusCode:     401,
		},
		{
			name:    "Failed. Invalid Timestamp Format",
			payload: payload,
			header: &auth.TokenHeader{
				ContentType: "application/json",
				ClientKey:   cfg.Client.Key,
				Timestamp:   now.Format("2006-01-02T15:04:05"),
				Signature:   "",
			},
			expectRespCode: "4007301",
			statusCode:     400,
		},
		{
			name:    "Failed. Invalid Signature",
			payload: payload,
			header: &auth.TokenHeader{
				ContentType: "application/json",
				ClientKey:   cfg.Client.Key,
				Timestamp:   now.Format(TSLayout),
				Signature:   "signature123123123",
			},
			expectRespCode: "4017300",
			statusCode:     401,
		},
	}

	for _, ts := range scenario {
		t.Run(ts.name, func(t *testing.T) {

			signature := strings.TrimSpace(ts.header.Signature)
			if signature == "" {
				strToSign := strings.Join([]string{ts.header.ClientKey, ts.header.Timestamp}, "|")
				genSig, err := generateSignature(strToSign)
				if err != nil {
					t.Errorf("Error Generate Signature - %v", err)
					return
				}
				signature = genSig
			}

			var buff bytes.Buffer
			_ = json.NewEncoder(&buff).Encode(ts.payload)

			request := httptest.NewRequest(http.MethodPost, endpoint, &buff)
			request.Header.Add("Content-Type", ts.header.ContentType)
			request.Header.Add("X-CLIENT-KEY", ts.header.ClientKey)
			request.Header.Add("X-TIMESTAMP", ts.header.Timestamp)
			request.Header.Add("X-SIGNATURE", signature)
			responseRecorder := httptest.NewRecorder()

			authSerice := handlers.NewHandlers(db, cfg)
			authSerice.GetAccessToken(responseRecorder, request)

			if responseRecorder.Code != ts.statusCode {
				t.Errorf("Expected status code %d, not HTTP %d", ts.statusCode, responseRecorder.Code)
				return
			}

			if responseRecorder.Code != http.StatusOK {
				var errCode response.Message
				json.Unmarshal(responseRecorder.Body.Bytes(), &errCode)
				if errCode.ResponseCode != ts.expectRespCode {
					t.Errorf("Expected response code '%s', got '%s - %s'", ts.expectRespCode, errCode.ResponseCode, errCode.ResponseMessage)
					return
				}
			}
		})
	}
}

func TestGetTodo(t *testing.T) {
	now := time.Now()
	srv := handlers.NewHandlers(db, cfg)

	accessToken, err := getToken(srv, now.Format(TSLayout))
	if err != nil {
		t.Error("Failed Get Access Token")
		return
	}

	scenario := []struct {
		name           string
		header         *todo.RequestHeader
		todoID         string
		expectRespCode string
		statusCode     int
	}{
		{
			name: "Success. Get All Todos ID",
			header: &todo.RequestHeader{
				ContentType:   "application/json",
				Authorization: "Bearer " + accessToken,
				ClientKey:     cfg.Client.Key,
				Timestamp:     now.Format(TSLayout),
				Signature:     "signature",
			},
			todoID:         "",
			expectRespCode: "2002400",
			statusCode:     http.StatusOK,
		},
		{
			name: "Success. Get Single Todos ID",
			header: &todo.RequestHeader{
				ContentType:   "application/json",
				Authorization: "Bearer " + accessToken,
				ClientKey:     cfg.Client.Key,
				Timestamp:     now.Format(TSLayout),
				Signature:     "signature",
			},
			todoID:         "uuid",
			expectRespCode: "2002400",
			statusCode:     http.StatusOK,
		},
		{
			name: "Failed. Invalid Access Token",
			header: &todo.RequestHeader{
				ContentType:   "application/json",
				Authorization: "Bearer ABC" + accessToken,
				ClientKey:     cfg.Client.Key,
				Timestamp:     now.Format(TSLayout),
				Signature:     "signature",
			},
			todoID:         "",
			expectRespCode: "4012401",
			statusCode:     http.StatusUnauthorized,
		},
	}

	for _, ts := range scenario {
		var tempID string

		t.Run(ts.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/v1.0/todo", nil)
			request.Header.Add("Content-Type", ts.header.ContentType)
			request.Header.Add("Authorization", ts.header.Authorization)
			request.Header.Add("X-CLIENT-KEY", ts.header.ClientKey)
			request.Header.Add("X-TIMESTAMP", ts.header.Timestamp)
			request.Header.Add("X-SIGNATURE", "signature")

			if ts.todoID != "" {
				request.SetPathValue("ID", tempID)
			}

			responseRecorder := httptest.NewRecorder()

			srv.GetTodoList(responseRecorder, request)
			if responseRecorder.Code != ts.statusCode {
				t.Errorf("Expected status code %d, not HTTP %d", ts.statusCode, responseRecorder.Code)
				return
			}

			if responseRecorder.Code != http.StatusOK {
				var errCode response.Message
				json.Unmarshal(responseRecorder.Body.Bytes(), &errCode)
				if errCode.ResponseCode != ts.expectRespCode {
					t.Errorf("Expected response code '%s', got '%s - %s'", ts.expectRespCode, errCode.ResponseCode, errCode.ResponseMessage)
					return
				}
			}

			// take single todo ID for later use on scenario get single data todo
			var data []*database.TableTodos
			json.Unmarshal(responseRecorder.Body.Bytes(), &data)
			if len(data) > 0 {
				tempID = data[0].ID
			}
		})
	}
}

func TestAddTodo(t *testing.T) {
	now := time.Now()
	endpoint := "/v1.0/todo"
	srv := handlers.NewHandlers(db, cfg)

	accessToken, err := getToken(srv, now.Format(TSLayout))
	if err != nil {
		t.Error("Failed Get Access Token")
		return
	}

	scenario := []struct {
		name           string
		header         *todo.RequestHeader
		body           *todo.AddTodosRequest
		expectRespCode string
		statusCode     int
	}{
		// header signature if empty, will be generated
		{
			name: "Success",
			header: &todo.RequestHeader{
				ContentType:   "application/json",
				Authorization: "Bearer " + accessToken,
				ClientKey:     cfg.Client.Key,
				Timestamp:     now.Format(TSLayout),
				Signature:     "",
			},
			body: &todo.AddTodosRequest{
				Title:      "Make unit testing",
				DetailTodo: "Scenario get access token & add todos",
			},
			expectRespCode: "2002400",
			statusCode:     http.StatusOK,
		},
		{
			name: "Success. Without Detail",
			header: &todo.RequestHeader{
				ContentType:   "application/json",
				Authorization: "Bearer " + accessToken,
				ClientKey:     cfg.Client.Key,
				Timestamp:     now.Format(TSLayout),
				Signature:     "",
			},
			body: &todo.AddTodosRequest{
				Title: "Todos Without Detail",
			},
			expectRespCode: "2002400",
			statusCode:     http.StatusOK,
		},
		{
			name: "Failed. Invalid Access Token",
			header: &todo.RequestHeader{
				ContentType:   "application/json",
				Authorization: "Bearer xyz" + accessToken,
				ClientKey:     cfg.Client.Key,
				Timestamp:     now.Format(TSLayout),
				Signature:     "",
			},
			body: &todo.AddTodosRequest{
				Title:      "Just trying",
				DetailTodo: "Nothing happen",
			},
			expectRespCode: "4012401",
			statusCode:     http.StatusUnauthorized,
		},
		{
			name: "Failed. Invalid Signature",
			header: &todo.RequestHeader{
				ContentType:   "application/json",
				Authorization: "Bearer " + accessToken,
				ClientKey:     cfg.Client.Key,
				Timestamp:     now.Format(TSLayout),
				Signature:     "signature12315234",
			},
			body: &todo.AddTodosRequest{
				Title:      "Just trying",
				DetailTodo: "Nothing happen",
			},
			expectRespCode: "4012400",
			statusCode:     http.StatusUnauthorized,
		},
		{
			name: "Failed. Invalid Client Key",
			header: &todo.RequestHeader{
				ContentType:   "application/json",
				Authorization: "Bearer " + accessToken,
				ClientKey:     cfg.Client.Key + "ABC",
				Timestamp:     now.Format(TSLayout),
				Signature:     "",
			},
			body: &todo.AddTodosRequest{
				Title:      "Just trying",
				DetailTodo: "Nothing happen",
			},
			expectRespCode: "4012400",
			statusCode:     http.StatusUnauthorized,
		},
	}

	for _, ts := range scenario {
		t.Run(ts.name, func(t *testing.T) {
			bytePayload, _ := json.Marshal(ts.body)

			var buff bytes.Buffer
			_ = json.Compact(&buff, bytePayload)

			signature := strings.TrimSpace(ts.header.Signature)
			// generate signature
			if signature == "" {
				h := sha256.New()
				h.Write(buff.Bytes())

				strToSign := strings.Join([]string{
					http.MethodPost,
					endpoint,
					accessToken,
					strings.ToLower(hex.EncodeToString(h.Sum(nil))),
					now.Format(TSLayout),
				}, ":")

				mac := hmac.New(sha512.New, []byte(cfg.Client.Secret))
				mac.Write([]byte(strToSign))

				signature = base64.StdEncoding.EncodeToString(mac.Sum(nil))
			}

			request := httptest.NewRequest(http.MethodPost, "/v1.0/todo", &buff)
			request.Header.Add("Content-Type", ts.header.ContentType)
			request.Header.Add("Authorization", ts.header.Authorization)
			request.Header.Add("X-CLIENT-KEY", ts.header.ClientKey)
			request.Header.Add("X-TIMESTAMP", ts.header.Timestamp)
			request.Header.Add("X-SIGNATURE", signature)
			responseRecorder := httptest.NewRecorder()

			srv.NewTodo(responseRecorder, request)
			if responseRecorder.Code != ts.statusCode {
				t.Errorf("Expected status code %d, not HTTP %d", ts.statusCode, responseRecorder.Code)
				return
			}

			if responseRecorder.Code != http.StatusOK {
				var errCode response.Message
				json.Unmarshal(responseRecorder.Body.Bytes(), &errCode)
				if errCode.ResponseCode != ts.expectRespCode {
					t.Errorf("Expected response code '%s', got '%s - %s'", ts.expectRespCode, errCode.ResponseCode, errCode.ResponseMessage)
					return
				}
			}
		})
	}
}

func getToken(srv *handlers.Handlers, ts string) (string, error) {
	payload := &auth.TokenRequest{
		GrantType: "client_credentials",
	}
	var buff bytes.Buffer
	_ = json.NewEncoder(&buff).Encode(payload)

	signature, err := generateSignature(strings.Join([]string{
		cfg.Client.Key, ts,
	}, "|"))
	if err != nil {
		return "", err
	}

	request := httptest.NewRequest(http.MethodPost, "/v1.0/access-token", &buff)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-CLIENT-KEY", cfg.Client.Key)
	request.Header.Add("X-TIMESTAMP", ts)
	request.Header.Add("X-SIGNATURE", signature)
	responseRecorder := httptest.NewRecorder()

	srv.GetAccessToken(responseRecorder, request)

	var tokenResp auth.TokenResponse
	json.Unmarshal(responseRecorder.Body.Bytes(), &tokenResp)
	return tokenResp.AccessToken, nil
}

func generateSignature(strToSign string) (string, error) {
	privKey, err := os.ReadFile("tests/partner_private_key.pem")
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(privKey)
	if block == nil {
		fmt.Println("Failed Decode Private Key")
		return "", errors.New("failed decode private key")
	}

	parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	privateKey := parseResult.(*rsa.PrivateKey)

	msghHash := sha256.New()
	msghHash.Write([]byte(strToSign))

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, msghHash.Sum(nil))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func init() {
	// up to parent for reading configs folder
	os.Chdir("..")

	c, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(0)
	}
	cfg = &c

	repo, err := database.BoostrapDatabase(&c)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(0)
	}
	db = repo
}
