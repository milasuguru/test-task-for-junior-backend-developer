package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type ScheduleRepository struct {
	pool *pgxpool.Pool
}

func NewScheduleRepository(pool *pgxpool.Pool) *ScheduleRepository {
	return &ScheduleRepository{pool: pool}
}

func (r *ScheduleRepository) Create(ctx context.Context, s *taskdomain.Schedule) (*taskdomain.Schedule, error) {
	const query = `
		INSERT INTO schedules (title, description, status, type, every_n_days, month_day, specific_dates, parity, generate_days_ahead, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, title, description, status, type, every_n_days, month_day, specific_dates, parity, generate_days_ahead, created_at, updated_at
	`

	dates := timesToStrings(s.SpecificDates)

	row := r.pool.QueryRow(ctx, query,
		s.Title, s.Description, s.Status, s.Type,
		s.EveryNDays, s.MonthDay, dates, s.Parity,
		s.GenerateDaysAhead, s.CreatedAt, s.UpdatedAt,
	)

	return scanSchedule(row)
}

func (r *ScheduleRepository) GetByID(ctx context.Context, id int64) (*taskdomain.Schedule, error) {
	const query = `
		SELECT id, title, description, status, type, every_n_days, month_day, specific_dates, parity, generate_days_ahead, created_at, updated_at
		FROM schedules WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	s, err := scanSchedule(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}
		return nil, err
	}

	return s, nil
}

func (r *ScheduleRepository) List(ctx context.Context) ([]taskdomain.Schedule, error) {
	const query = `
		SELECT id, title, description, status, type, every_n_days, month_day, specific_dates, parity, generate_days_ahead, created_at, updated_at
		FROM schedules ORDER BY id DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []taskdomain.Schedule
	for rows.Next() {
		s, err := scanSchedule(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, *s)
	}

	return result, rows.Err()
}

func (r *ScheduleRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM schedules WHERE id = $1`

	res, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return taskdomain.ErrNotFound
	}

	return nil
}

type scheduleScanner interface {
	Scan(dest ...any) error
}

func scanSchedule(scanner scheduleScanner) (*taskdomain.Schedule, error) {
	var (
		s             taskdomain.Schedule
		status        string
		scheduleType  string
		parity        *string
		specificDates []time.Time
	)

	if err := scanner.Scan(
		&s.ID,
		&s.Title,
		&s.Description,
		&status,
		&scheduleType,
		&s.EveryNDays,
		&s.MonthDay,
		&specificDates,
		&parity,
		&s.GenerateDaysAhead,
		&s.CreatedAt,
		&s.UpdatedAt,
	); err != nil {
		return nil, err
	}

	s.Status = taskdomain.Status(status)
	s.Type = taskdomain.ScheduleType(scheduleType)
	s.SpecificDates = specificDates

	if parity != nil {
		p := taskdomain.Parity(*parity)
		s.Parity = &p
	}

	return &s, nil
}

func timesToStrings(dates []time.Time) []string {
	result := make([]string, 0, len(dates))
	for _, d := range dates {
		result = append(result, d.Format("2006-01-02"))
	}
	return result
}
