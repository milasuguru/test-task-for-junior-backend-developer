package task

import (
	"context"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type ScheduleRepository interface {
	Create(ctx context.Context, s *taskdomain.Schedule) (*taskdomain.Schedule, error)
	GetByID(ctx context.Context, id int64) (*taskdomain.Schedule, error)
	List(ctx context.Context) ([]taskdomain.Schedule, error)
	Delete(ctx context.Context, id int64) error
}

type ScheduleUsecase interface {
	Create(ctx context.Context, input CreateScheduleInput) (*taskdomain.Schedule, error)
	GetByID(ctx context.Context, id int64) (*taskdomain.Schedule, error)
	List(ctx context.Context) ([]taskdomain.Schedule, error)
	Delete(ctx context.Context, id int64) error
	GenerateTasks(ctx context.Context) error
}

type CreateScheduleInput struct {
	Title             string
	Description       string
	Type              taskdomain.ScheduleType
	EveryNDays        *int
	MonthDay          *int
	SpecificDates     []string
	Parity            *taskdomain.Parity
	GenerateDaysAhead *int
}
