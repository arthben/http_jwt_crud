package todo

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/arthben/http_jwt_crud/internal/config"
	dbs "github.com/arthben/http_jwt_crud/internal/database"
	"github.com/arthben/http_jwt_crud/internal/logging"
	res "github.com/arthben/http_jwt_crud/internal/response"
	"github.com/jmoiron/sqlx"
)

const (
	DateLayout  = "2006-01-02 15:04:05"
	UnCompleted = "0"
)

type TodoService struct {
	db  *sqlx.DB
	cfg *config.EnvParams
}

func NewTodoService(db *sqlx.DB, cfg *config.EnvParams) *TodoService {
	return &TodoService{db: db, cfg: cfg}
}

func (t *TodoService) Get(w http.ResponseWriter, r *http.Request) ([]*dbs.TableTodos, *res.Message) {
	_, errCode := validateHeader(r, t.cfg)
	if errCode != nil {
		return nil, errCode
	}

	// get the ID
	id := r.PathValue("ID")
	resp, err := dbs.ListTodo(t.db, context.TODO(), id)
	if err != nil && id != "" {
		return nil, res.BadResponse(http.StatusNotFound, ServiceCode, "00", "No Data Found")
	}

	return resp, nil
}

func (t *TodoService) Add(w http.ResponseWriter, r *http.Request) (*dbs.TableTodos, *res.Message) {
	logger, _ := logging.FromContext(r.Context())
	header, errCode := validateHeader(r, t.cfg)
	if errCode != nil {
		return nil, errCode
	}

	payload, errCode := validateAddBody(w, r)
	if errCode != nil {
		return nil, errCode
	}

	// skip err because on validateBody already handle it
	body, _ := json.Marshal(payload)

	// validate signature
	if errCode := validateSignature(header, body, r.URL.Path, r.Method, t.cfg.Client.Secret); errCode != nil {
		return nil, errCode
	}

	// skip err because already check in validateHeader
	createdDate, _ := time.Parse(TSLayout, header.Timestamp)

	now := time.Now()
	tbl := &dbs.TableTodos{
		ID:              generateIDTodos(),
		Title:           payload.Title,
		Detail:          payload.DetailTodo,
		CreatedDate:     createdDate.Format(DateLayout),
		UpdatedDate:     now.Format(DateLayout),
		StatusCompleted: UnCompleted,
	}
	if err := dbs.AddTodo(t.db, context.TODO(), tbl); err != nil {
		if logger != nil {
			logger.Error("AddTodo", slog.String("error", err.Error()))
		} else {
			log.Fatalf("AddTodo %v", err)
		}
		return nil, res.BadResponse(http.StatusInternalServerError, ServiceCode, "00", "Error Internal Server")
	}

	return tbl, nil
}
