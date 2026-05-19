package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
)

type TaskController struct {
	taskService app.TaskService
}

func NewTaskController(ts app.TaskService) TaskController {
	return TaskController{
		taskService: ts,
	}
}
func (c TaskController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController.Save(request.Bind): %s", err)
			BadRequest(w, errors.New("invalid request body"))
			return
		}

		task.UserId = user.Id
		task.Status = domain.NewTaskStatus

		task, err = c.taskService.Save(task)
		if err != nil {
			log.Printf("TaskController.Save(c.taskService.Save): %s", err)
			InternalServerError(w, err)
			return

		}

		taskDto := resources.TaskDto{}
		taskDto = taskDto.DomainToDto(task)

		Success(w, taskDto)
	}
}
func (c TaskController) FindList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		// Зчитуємо параметри фільтрації та сортування з URL (наприклад: ?status=NEW&sortBy=-created_date)
		queryStatus := r.URL.Query().Get("status")
		querySort := r.URL.Query().Get("sortBy")

		// Наповнюємо фільтр живими даними
		filter := domain.TaskFilter{
			Status: domain.TaskStatus(queryStatus),
			SortBy: querySort,
		}

		// Передаємо заповнений фільтр у сервіс
		tasks, err := c.taskService.FindList(user.Id, filter)
		if err != nil {
			log.Printf("TaskController.FindList(c.taskService.FindList): %s", err)
			InternalServerError(w, err)
			return
		}

		taskDto := resources.TaskDto{}
		tasksDto := taskDto.DomainToDtoCollection(tasks)

		Success(w, tasksDto)
	}
}
func (c TaskController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		task := r.Context().Value(TaskKey).(domain.Task)

		if task.UserId != user.Id {
			Forbidden(w, errors.New("access denied"))
			return
		}

		taskDto := resources.TaskDto{}
		taskDto = taskDto.DomainToDto(task)

		Success(w, taskDto)
	}
}

//todo: add method for chang (update) Task status

func (c TaskController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		task := r.Context().Value(TaskKey).(domain.Task)

		if task.UserId != user.Id {
			Forbidden(w, errors.New("access denied"))
			return
		}

		// Розпаковуємо JSON із Postman (завдяки попередньому кроку тут уже є статус)
		updTask, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController.Update(requests.Bind): %s", err)
			BadRequest(w, errors.New("invalid request body"))
			return
		}

		// Копіюємо оновлені поля
		task.Title = updTask.Title
		task.Description = updTask.Description
		task.Deadline = updTask.Deadline
		task.Status = updTask.Status // <-- ОСЬ ЦЕЙ РЯДОК ВСЕ ОЖИВЛЯЄ!

		// Зберігаємо в базу даних
		task, err = c.taskService.Update(task)
		if err != nil {
			log.Printf("TaskController.Update(c.taskService.Update): %s", err)
			InternalServerError(w, err)
			return
		}

		taskDto := resources.TaskDto{}
		taskDto = taskDto.DomainToDto(task)

		Success(w, taskDto)
	}

}

func (c TaskController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		task := r.Context().Value(TaskKey).(domain.Task)

		if task.UserId != user.Id {
			Forbidden(w, errors.New("access denied"))

			return
		}
		err := c.taskService.Delete(task.Id)
		if err != nil {
			log.Printf("TaskController.Delete(c.taskService.Delete): %s", err)
			InternalServerError(w, err)
			return
		}

		noContent(w)
	}
}
