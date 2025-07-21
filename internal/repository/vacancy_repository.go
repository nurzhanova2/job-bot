package repository

import (
	"sync"

	"github.com/nurzhanova2/job-bot.git/internal/model"
)

type VacancyRepository interface {
	Save(v model.Vacancy) error
	GetAll() []model.Vacancy
	Clear()
}

type InMemoryVacancyRepository struct {
	mu       sync.RWMutex
	vacancies []model.Vacancy
}

func NewInMemoryVacancyRepository() *InMemoryVacancyRepository {
	return &InMemoryVacancyRepository{
		vacancies: make([]model.Vacancy, 0),
	}
}

func (r *InMemoryVacancyRepository) Save(v model.Vacancy) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.vacancies = append(r.vacancies, v)
	return nil
}

func (r *InMemoryVacancyRepository) GetAll() []model.Vacancy {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.vacancies
}

func (r *InMemoryVacancyRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.vacancies = []model.Vacancy{}
}

