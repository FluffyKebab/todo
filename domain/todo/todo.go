package todo

type User struct {
	ID   string
	Name string
}

type Todo struct {
	ID     string
	UserID string
	Body   string
	Done   bool
}

type UserService interface {
	CreateUser(u User) (string, error)
}

type TodoService interface {
	CreateTodo(t Todo) (string, error)
	UpdateTodo(t Todo) error
	DeleteTodo(id string) error
	GetUserTodos(userId string) error
}
