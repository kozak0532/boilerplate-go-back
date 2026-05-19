package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const TasksTableName = "tasks"

type task struct {
	Id          uint64            `db:"id,omitempty"`
	UserId      uint64            `db:"user_id"`
	Title       string            `db:"title"`
	Description string            `db:"description"`
	Status      domain.TaskStatus `db:"status"`
	Deadline    *time.Time        `db:"deadline"`
	CreatedDate time.Time         `db:"created_date"`
	UpdatedDate time.Time         `db:"updated_date"`
	DeletedDate *time.Time        `db:"deleted_date"`
}

type TaskRepository interface {
	Save(t domain.Task) (domain.Task, error)
	FindList(uId uint64, filter domain.TaskFilter) ([]domain.Task, error) // Змінено тут
	Find(id uint64) (domain.Task, error)
	Update(t domain.Task) (domain.Task, error)
	Delete(id uint64) error
}

type taskRepository struct {
	coll db.Collection
	sess db.Session
}

func NewTaskRepository(session db.Session) TaskRepository {
	return taskRepository{
		coll: session.Collection(TasksTableName),
		sess: session,
	}
}
func (r taskRepository) Save(t domain.Task) (domain.Task, error) {
	tsk := r.mapDomainToModel(t)
	tsk.CreatedDate = time.Now()
	tsk.UpdatedDate = time.Now()

	err := r.coll.InsertReturning(&tsk)
	if err != nil {
		return domain.Task{}, err
	}

	t = r.mapModelToDomain(tsk)
	return t, nil
}
func (r taskRepository) FindList(uId uint64, filter domain.TaskFilter) ([]domain.Task, error) {
	var tasks []task

	// 1. Базові умови (шукаємо завдання конкретного юзера, які не видалені)
	cond := db.Cond{
		"user_id":      uId,
		"deleted_date": nil,
	}

	// 2. Додаємо фільтри (якщо вони передані)
	if filter.Status != "" {
		cond["status"] = filter.Status
	}
	if filter.Deadline != nil {
		// Наприклад, шукаємо завдання, дедлайн яких менше або дорівнює переданому
		cond["deadline <="] = filter.Deadline
	}

	// 3. Формуємо запит із нашими умовами
	query := r.coll.Find(cond)

	// 4. Додаємо сортування
	if filter.SortBy != "" {
		// upper/db розуміє мінус "-" як сортування за спаданням (DESC)
		// Наприклад: "deadline" (за зростанням), "-created_date" (за спаданням)
		query = query.OrderBy(filter.SortBy)
	} else {
		// Сортування за замовчуванням (нові зверху)
		query = query.OrderBy("-created_date")
	}

	// 5. Виконуємо фінальний запит
	err := query.All(&tasks)
	if err != nil {
		return nil, err
	}

	return r.mapModelToDomainCollection(tasks), nil
}

func (r taskRepository) Find(id uint64) (domain.Task, error) {
	var t task

	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&t)
	if err != nil {
		return domain.Task{}, err
	}
	return r.mapModelToDomain(t), nil

}
func (r taskRepository) Update(t domain.Task) (domain.Task, error) {
	tsk := r.mapDomainToModel(t)

	tsk.UpdatedDate = time.Now()

	err := r.coll.Find(db.Cond{"id": t.Id, "deleted_date": nil}).Update(&tsk)

	if err != nil {
		return domain.Task{}, err
	}

	t = r.mapModelToDomain(tsk)
	return t, nil
}

func (r taskRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r taskRepository) mapDomainToModel(t domain.Task) task {
	return task{
		Id:          t.Id,
		UserId:      t.UserId,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Deadline:    t.Deadline,
		CreatedDate: t.CreatedDate,
		UpdatedDate: t.UpdatedDate,
		DeletedDate: t.DeletedDate,
	}
}

func (r taskRepository) mapModelToDomain(t task) domain.Task {
	return domain.Task{
		Id:          t.Id,
		UserId:      t.UserId,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Deadline:    t.Deadline,
		CreatedDate: t.CreatedDate,
		UpdatedDate: t.UpdatedDate,
		DeletedDate: t.DeletedDate,
	}
}

func (r taskRepository) mapModelToDomainCollection(ts []task) []domain.Task {
	tasks := make([]domain.Task, len(ts))
	for i := range ts {
		tasks[i] = r.mapModelToDomain(ts[i])
	}
	return tasks
}
