package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// var (
// 	// ErrUserNotFound is user not found.
// 	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
// )

// Todo is a Todo model.
type Todo struct {
	ID     int64
	Title  string
	Status bool
}

// TodoRepo is a Todo repo.
// biz 层对数据操作层提出了以下接口要求，不管实际存储的是Mysql还是Redis
type TodoRepo interface {
	Save(context.Context, *Todo) (*Todo, error)
	Update(context.Context, *Todo) error
	Delete(context.Context, int64) error
	FindByID(context.Context, int64) (*Todo, error)
	ListAll(context.Context) ([]*Todo, error)
}

// TodoUsecase is a Todo usecase.
type TodoUsecase struct {
	repo TodoRepo
	log  *log.Helper
}

// NewTodoUsecase new a Todo usecase.
func NewTodoUsecase(repo TodoRepo, logger log.Logger) *TodoUsecase {
	return &TodoUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateTodo creates a Todo, and returns the new Todo.
// 对外提供的业务函数，实现复杂的业务逻辑
func (uc *TodoUsecase) CreateTodo(ctx context.Context, t *Todo) (*Todo, error) {
	uc.log.WithContext(ctx).Infof("CreateTodo: %#v", t)
	return uc.repo.Save(ctx, t) // 调用下一层的Save方法
}

func (uc *TodoUsecase) Get(ctx context.Context, id int64) (*Todo, error) {
	uc.log.WithContext(ctx).Infof("Get: %#v", id)
	return uc.repo.FindByID(ctx, id) // 调用下一层的 FindByID 方法
}
