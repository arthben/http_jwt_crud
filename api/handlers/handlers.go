package handlers

import (
	"log/slog"
	"net/http"

	"github.com/arthben/http_jwt_crud/api/auth"
	"github.com/arthben/http_jwt_crud/api/todo"
	"github.com/arthben/http_jwt_crud/internal/config"
	"github.com/arthben/http_jwt_crud/internal/response"
	"github.com/jmoiron/sqlx"
)

type Handlers struct {
	mux    *http.ServeMux
	logger *slog.Logger
	auth   *auth.AuthService
	todo   *todo.TodoService
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

func (h *Handlers) BuildRouter() (http.Handler, error) {
	h.mux.HandleFunc("GET /halo", halo)
	h.mux.HandleFunc("POST /v1.0/access-token", h.GetAccessToken)
	h.mux.HandleFunc("POST /v1.0/todo", h.NewTodo)

	return http.Handler(h.mux), nil
}

func (h *Handlers) GetAccessToken(w http.ResponseWriter, r *http.Request) {
	resp, errCode := h.auth.GetAccessToken(w, r)
	if errCode != nil {
		response.Write(w).AbortWithJSON(errCode)
		return
	}

	response.Write(w).JSON(resp)
}

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
