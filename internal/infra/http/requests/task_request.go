package requests

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type TaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Deadline    *int64 `json:"deadline"`
	Status      string `json:"status"` // <-- ДОДАНО: тепер статус приймається з Postman
}

func (r TaskRequest) ToDomainModel() (interface{}, error) {
	var timeUnix int64
	if r.Deadline != nil {
		timeUnix = *r.Deadline
	}
	var deadline *time.Time
	if timeUnix != 0 {
		d1 := time.Unix(timeUnix, 0)
		deadline = &d1
	}

	return domain.Task{
		Title:       r.Title,
		Description: r.Description,
		Deadline:    deadline,
		Status:      domain.TaskStatus(r.Status), // <-- ДОДАНО: передаємо статус у доменну модель
	}, nil
}
