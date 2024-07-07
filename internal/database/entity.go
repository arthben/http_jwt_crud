package database

type TableTodos struct {
	ID              string `db:"id"`
	Title           string `db:"title"`
	Detail          string `db:"detail"`
	CreatedDate     string `db:"created_date"`
	UpdatedDate     string `db:"updated_date"`
	StatusCompleted string `db:"st_completed"`
	CompletedDate   string `db:"completed_date"`
}
