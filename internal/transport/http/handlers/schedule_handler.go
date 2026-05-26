package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	taskdomain "example.com/taskservice/internal/domain/task"
	taskusecase "example.com/taskservice/internal/usecase/task"
)

type ScheduleHandler struct {
	usecase taskusecase.ScheduleUsecase
}

func NewScheduleHandler(usecase taskusecase.ScheduleUsecase) *ScheduleHandler {
	return &ScheduleHandler{usecase: usecase}
}

type createScheduleDTO struct {
	Title             string                `json:"title"`
	Description       string                `json:"description"`
	Type              taskdomain.ScheduleType `json:"type"`
	EveryNDays        *int                  `json:"every_n_days,omitempty"`
	MonthDay          *int                  `json:"month_day,omitempty"`
	SpecificDates     []string              `json:"specific_dates,omitempty"`
	Parity            *taskdomain.Parity    `json:"parity,omitempty"`
	GenerateDaysAhead *int                  `json:"generate_days_ahead,omitempty"`
}

type scheduleResponseDTO struct {
	ID                int64                   `json:"id"`
	Title             string                  `json:"title"`
	Description       string                  `json:"description"`
	Type              taskdomain.ScheduleType `json:"type"`
	EveryNDays        *int                    `json:"every_n_days,omitempty"`
	MonthDay          *int                    `json:"month_day,omitempty"`
	SpecificDates     []string                `json:"specific_dates,omitempty"`
	Parity            *taskdomain.Parity      `json:"parity,omitempty"`
	GenerateDaysAhead int                     `json:"generate_days_ahead"`
	CreatedAt         time.Time               `json:"created_at"`
	UpdatedAt         time.Time               `json:"updated_at"`
}

func newScheduleDTO(s *taskdomain.Schedule) scheduleResponseDTO {
	dto := scheduleResponseDTO{
		ID:                s.ID,
		Title:             s.Title,
		Description:       s.Description,
		Type:              s.Type,
		EveryNDays:        s.EveryNDays,
		MonthDay:          s.MonthDay,
		Parity:            s.Parity,
		GenerateDaysAhead: s.GenerateDaysAhead,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}

	for _, d := range s.SpecificDates {
		dto.SpecificDates = append(dto.SpecificDates, d.Format("2006-01-02"))
	}

	return dto
}

func (h *ScheduleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createScheduleDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	created, err := h.usecase.Create(r.Context(), taskusecase.CreateScheduleInput{
		Title:             req.Title,
		Description:       req.Description,
		Type:              req.Type,
		EveryNDays:        req.EveryNDays,
		MonthDay:          req.MonthDay,
		SpecificDates:     req.SpecificDates,
		Parity:            req.Parity,
		GenerateDaysAhead: req.GenerateDaysAhead,
	})
	if err != nil {
		writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, newScheduleDTO(created))
}

func (h *ScheduleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getScheduleID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	s, err := h.usecase.GetByID(r.Context(), id)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newScheduleDTO(s))
}

func (h *ScheduleHandler) List(w http.ResponseWriter, r *http.Request) {
	schedules, err := h.usecase.List(r.Context())
	if err != nil {
		writeUsecaseError(w, err)
		return
	}

	response := make([]scheduleResponseDTO, 0, len(schedules))
	for i := range schedules {
		response = append(response, newScheduleDTO(&schedules[i]))
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *ScheduleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getScheduleID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.usecase.Delete(r.Context(), id); err != nil {
		writeUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getScheduleID(r *http.Request) (int64, error) {
	rawID := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil || id <= 0 {
		return 0, strconv.ErrSyntax
	}
	return id, nil
}
