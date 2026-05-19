package domain

import (
	"time"
)

type Task struct {
	Id          uint64
	UserId      uint64
	Title       string
	Description string
	Status      TaskStatus
	Deadline    *time.Time
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

type TaskStatus string

const (
	NewTaskStatus        TaskStatus = "NEW"
	InProgressTaskStatus TaskStatus = "IN_PROGRESS"
	DoneTaskStatus       TaskStatus = "DONE"
)

// TaskFilter містить параметри для пошуку та сортування
type TaskFilter struct {
	Status   TaskStatus // Фільтр за статусом (наприклад: "NEW")
	Deadline *time.Time // Фільтр за дедлайном
	SortBy   string     // Поле для сортування (наприклад: "deadline", "-created_date")
}
