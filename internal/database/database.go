package database

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/arthben/http_jwt_crud/internal/config"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func BoostrapDatabase(cfg *config.EnvParams) (*sqlx.DB, error) {
	dsn := mysql.Config{
		User:                 cfg.DB.Username,
		Passwd:               cfg.DB.Password,
		Net:                  "tcp",
		Addr:                 cfg.DB.Host,
		DBName:               cfg.DB.Name,
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
	}
	dbase, err := sqlx.Connect("mysql", dsn.FormatDSN())
	if err != nil || dbase == nil {
		return nil, err
	}

	maxIdleConn, _ := strconv.Atoi(cfg.DB.MinPool)
	maxOpenConn, _ := strconv.Atoi(cfg.DB.MaxPool)
	dbase.SetMaxIdleConns(maxIdleConn)
	dbase.SetMaxOpenConns(maxOpenConn)

	return dbase, nil
}

func CloseDatabase(db *sqlx.DB) error {
	return db.Close()
}

func ListTodo(db *sqlx.DB, ctx context.Context, id string) ([]*TableTodos, error) {
	var args []string

	sql := `SELECT id, title, detail, 
			CONCAT('', created_date) AS created_date, 
			CONCAT('', updated_date) AS updated_date, st_completed
			FROM todos
			WHERE 1=1`

	if id != "" {
		args = append(args, id)
		sql = strings.Join([]string{sql, "AND id=?"}, " ")
	}

	sql = strings.Join([]string{sql, "ORDER BY updated_date DESC"}, " ")

	resp := []*TableTodos{}
	err := db.SelectContext(ctx, &resp, sql, args)
	return resp, err
}

func AddTodo(db *sqlx.DB, ctx context.Context, tb *TableTodos) error {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	// any error will be rollback
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	sql := `INSERT INTO todos(id, title, detail, created_date, updated_date, st_completed)
			VALUES(:id, :title, :detail, :created_date, :updated_date, :st_completed)`
	prep, err := tx.PrepareNamedContext(ctx, sql)
	if err != nil {
		return err
	}

	_, err = prep.ExecContext(ctx, tb)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
