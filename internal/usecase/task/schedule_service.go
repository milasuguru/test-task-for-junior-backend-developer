package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type ScheduleService struct {
	scheduleRepo ScheduleRepository
	taskRepo     Repository
	now          func() time.Time
}

func NewScheduleService(scheduleRepo ScheduleRepository, taskRepo Repository) *ScheduleService {
	return &ScheduleService{
		scheduleRepo: scheduleRepo,
		taskRepo:     taskRepo,
		now:          func() time.Time { return time.Now().UTC() },
	}
}

func (s *ScheduleService) Create(ctx context.Context, input CreateScheduleInput) (*taskdomain.Schedule, error) {
	if err := validateScheduleInput(input); err != nil {
		return nil, err
	}

	now := s.now()

	schedule := &taskdomain.Schedule{
		Title:             strings.TrimSpace(input.Title),
		Description:       strings.TrimSpace(input.Description),
		Status:            taskdomain.StatusNew,
		Type:              input.Type,
		EveryNDays:        input.EveryNDays,
		MonthDay:          input.MonthDay,
		Parity:            input.Parity,
		GenerateDaysAhead: 30,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if input.GenerateDaysAhead != nil && *input.GenerateDaysAhead > 0 {
		schedule.GenerateDaysAhead = *input.GenerateDaysAhead
	}

	for _, ds := range input.SpecificDates {
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid date format %q, use YYYY-MM-DD", ErrInvalidInput, ds)
		}
		schedule.SpecificDates = append(schedule.SpecificDates, t)
	}

	created, err := s.scheduleRepo.Create(ctx, schedule)
	if err != nil {
		return nil, err
	}

	if err := s.generateTasksForSchedule(ctx, created); err != nil {
		return nil, err
	}

	return created, nil
}

func (s *ScheduleService) GetByID(ctx context.Context, id int64) (*taskdomain.Schedule, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}
	return s.scheduleRepo.GetByID(ctx, id)
}

func (s *ScheduleService) List(ctx context.Context) ([]taskdomain.Schedule, error) {
	return s.scheduleRepo.List(ctx)
}

func (s *ScheduleService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}
	return s.scheduleRepo.Delete(ctx, id)
}

func (s *ScheduleService) GenerateTasks(ctx context.Context) error {
	schedules, err := s.scheduleRepo.List(ctx)
	if err != nil {
		return err
	}

	for i := range schedules {
		if err := s.generateTasksForSchedule(ctx, &schedules[i]); err != nil {
			return err
		}
	}

	return nil
}

func (s *ScheduleService) generateTasksForSchedule(ctx context.Context, schedule *taskdomain.Schedule) error {
	now := s.now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	for i := 0; i < schedule.GenerateDaysAhead; i++ {
		date := today.AddDate(0, 0, i)

		if !schedule.ShouldCreateOn(date) {
			continue
		}

		task := &taskdomain.Task{
			Title:       schedule.Title,
			Description: schedule.Description,
			Status:      taskdomain.StatusNew,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		if _, err := s.taskRepo.Create(ctx, task); err != nil {
			return err
		}
	}

	return nil
}

func validateScheduleInput(input CreateScheduleInput) error {
	if strings.TrimSpace(input.Title) == "" {
		return fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Type.Valid() {
		return fmt.Errorf("%w: invalid schedule type", ErrInvalidInput)
	}

	switch input.Type {
	case taskdomain.ScheduleTypeDaily:
		if input.EveryNDays == nil || *input.EveryNDays <= 0 {
			return fmt.Errorf("%w: every_n_days must be positive for daily schedule", ErrInvalidInput)
		}
	case taskdomain.ScheduleTypeMonthly:
		if input.MonthDay == nil || *input.MonthDay < 1 || *input.MonthDay > 30 {
			return fmt.Errorf("%w: month_day must be between 1 and 30", ErrInvalidInput)
		}
	case taskdomain.ScheduleTypeSpecific:
		if len(input.SpecificDates) == 0 {
			return fmt.Errorf("%w: specific_dates must not be empty", ErrInvalidInput)
		}
	case taskdomain.ScheduleTypeParity:
		if input.Parity == nil || !input.Parity.Valid() {
			return fmt.Errorf("%w: parity must be 'even' or 'odd'", ErrInvalidInput)
		}
	}

	return nil
}
