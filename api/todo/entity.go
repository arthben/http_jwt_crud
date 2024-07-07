package todo

type RequestHeader struct {
	ContentType   string `header:"Content-Type" binding:"required"`
	Authorization string `header:"Authorization" binding:"required"`
	ClientKey     string `header:"X-Client-Key" binding:"required"`
	Timestamp     string `header:"X-Timestamp" binding:"required"`
	Signature     string `header:"X-Signature" binding:"required"`
}

type AddTodosRequest struct {
	Title      string `json:"title" binding:"required,min=1,max=256"`
	DetailTodo string `json:"detail_todo" binding:"omitempty,max=1024"`
}
