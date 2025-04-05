package models

// Task модель задачи пользователя
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// TasksList модель списка задач пользователя
type TasksList struct {
	Tasks []*Task `json:"tasks" default:"[]"`
}
