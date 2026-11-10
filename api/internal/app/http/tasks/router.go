package tasks

import (
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	tasksHandler "github.com/ArthurMilligan/charlie/api/internal/app/handlers/tasks"
	"github.com/go-chi/chi/v5"
)

func NewTasksRouter(r chi.Router, t interfaces.TasksService) {
	r.Post("/add", tasksHandler.AddTask(t))
	r.Post("/get", tasksHandler.GetTasks(t))
}
