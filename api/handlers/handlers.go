package handlers

import (
	"net/http"

	"github.com/arthben/http_jwt_crud/api/auth"
	"github.com/arthben/http_jwt_crud/api/todo"

	"github.com/arthben/http_jwt_crud/docs"
	"github.com/arthben/http_jwt_crud/internal/config"
	"github.com/arthben/http_jwt_crud/internal/response"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/jmoiron/sqlx"
)

type Handlers struct {
	mux  *http.ServeMux
	auth *auth.AuthService
	todo *todo.TodoService
}

func NewHandlers(
	db *sqlx.DB,
	cfg *config.EnvParams,
) *Handlers {
	return &Handlers{
		mux:  http.NewServeMux(),
		auth: auth.NewAuthService(cfg),
		todo: todo.NewTodoService(db, cfg),
	}
}

// @title HTTP JWT CRUD
// @version 1.0
// @description Demonstrate HTTP with Middleware, JWT, SQLX and slog package
// @contact.name Yohanes Catur
// @contact.url www.linkedin.com/in/yohanescatur
// @contact.email yohanescatur@gmail.com
func (h *Handlers) BuildRouter() (http.Handler, error) {
	// Load swagger
	h.InitSwagger()

	h.mux.HandleFunc("GET /halo", halo)
	h.mux.HandleFunc("POST /v1.0/access-token", h.GetAccessToken)
	h.mux.HandleFunc("POST /v1.0/todo", h.NewTodo)

	return http.Handler(h.mux), nil
}

func (h *Handlers) InitSwagger() {
	docs.SwaggerInfo.Title = "HTTP JWT CRUD Documentation"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}
	docs.SwaggerInfo.Description = `API For Todo.<br />
	Response Code format : 2007300<br />
	Respon Code pattern : HTTP Status - Service Code - Response Code`

	h.mux.HandleFunc("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

}

// GetAccessToken godoc
// @Summary Get Access Token
// @Tags Token
// @Accept json
// @Produce json
// @Param Content-Type header string true "application/json"
// @Param X-Client-Key header string true "Client key provided by server"
// @Param X-Timestamp  header string true "format: 2006-01-02T15:04:05+07:00"
// @Param X-Signature  header string true "Generated Signature"
// @Success 200 {object} auth.TokenResponse
// @Failure 400 {object} response.Message
// @Failure 401 {object} response.Message
// @Router /v1.0/access-token [POST]
func (h *Handlers) GetAccessToken(w http.ResponseWriter, r *http.Request) {
	resp, errCode := h.auth.GetAccessToken(w, r)
	if errCode != nil {
		response.Write(w).AbortWithJSON(errCode)
		return
	}

	response.Write(w).JSON(resp)
}

// NewTodo godoc
// @Summary Add Todo item
// @Description.markdown todo_add
// @Tags Todo
// @Accept json
// @Produce json
// @Param Content-Type  header string true "application/json"
// @Param Authorization header string true "Bearer token"
// @Param X-Client-Key  header string true "Client Key provided by server"
// @Param X-Timestamp   header string true "format: 2006-01-02T15:04:05+07:00"
// @Param X-Signature   header string true "123425234"
// @Param request       body   todo.AddTodosRequest true "request"
// @Success 200 {object} database.TableTodos
// @Failure 400 {object} response.Message
// @Router /v1.0/todo [POST]
func (h *Handlers) NewTodo(w http.ResponseWriter, r *http.Request) {
	tb, errCode := h.todo.Add(w, r)
	if errCode != nil {
		response.Write(w).AbortWithJSON(errCode)
		return
	}

	response.Write(w).JSON(tb)
}

func halo(w http.ResponseWriter, r *http.Request) {
	response.Write(w).JSON(&response.Message{
		HttpStatus:      http.StatusOK,
		ResponseCode:    "00",
		ResponseMessage: "Halo",
	})
}
