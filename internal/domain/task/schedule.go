package task

import "time"

type ScheduleType string

const (
	ScheduleTypeDaily    ScheduleType = "daily"
	ScheduleTypeMonthly  ScheduleType = "monthly"
	ScheduleTypeSpecific ScheduleType = "specific"
	ScheduleTypeParity   ScheduleType = "parity"
)

type Parity string

const (
	ParityEven Parity = "even"
	ParityOdd  Parity = "odd"
)

type Schedule struct {
	ID                int64
	Title             string
	Description       string
	Status            Status
	Type              ScheduleType
	EveryNDays        *int
	MonthDay          *int
	SpecificDates     []time.Time
	Parity            *Parity
	GenerateDaysAhead int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (s ScheduleType) Valid() bool {
	switch s {
	case ScheduleTypeDaily, ScheduleTypeMonthly, ScheduleTypeSpecific, ScheduleTypeParity:
		return true
	default:
		return false
	}
}

func (p Parity) Valid() bool {
	switch p {
	case ParityEven, ParityOdd:
		return true
	default:
		return false
	}
}

// ShouldCreateOn возвращает true если задача должна быть создана на указанную дату
func (s *Schedule) ShouldCreateOn(date time.Time) bool {
	day := date.Day()

	switch s.Type {
	case ScheduleTypeDaily:
		if s.EveryNDays == nil {
			return false
		}
		// Считаем дни от CreatedAt
		diff := int(date.Sub(s.CreatedAt.Truncate(24*time.Hour)).Hours() / 24)
		return diff >= 0 && diff%*s.EveryNDays == 0

	case ScheduleTypeMonthly:
		if s.MonthDay == nil {
			return false
		}
		return day == *s.MonthDay

	case ScheduleTypeSpecific:
		for _, d := range s.SpecificDates {
			if d.Year() == date.Year() && d.Month() == date.Month() && d.Day() == date.Day() {
				return true
			}
		}
		return false

	case ScheduleTypeParity:
		if s.Parity == nil {
			return false
		}
		if *s.Parity == ParityEven {
			return day%2 == 0
		}
		return day%2 != 0
	}

	return false
}
