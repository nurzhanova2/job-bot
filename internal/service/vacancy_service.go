package service

import (
	"strings"

	"github.com/nurzhanova2/job-bot.git/internal/model"
	"github.com/nurzhanova2/job-bot.git/internal/repository"
)

type VacancyService struct {
	repo     repository.VacancyRepository
	keywords []string
}

func NewVacancyService(repo repository.VacancyRepository, keywords []string) *VacancyService {
	return &VacancyService{
		repo:     repo,
		keywords: keywords,
	}
}

func (s *VacancyService) SaveIfRelevant(v model.Vacancy) error {
	if s.isRelevant(v) {
		return s.repo.Save(v)
	}
	return nil
}

func (s *VacancyService) isRelevant(v model.Vacancy) bool {
	text := strings.ToLower(v.Title + " " + v.Description)
	for _, keyword := range s.keywords {
		if strings.Contains(text, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func (s *VacancyService) GetAll() []model.Vacancy {
	return s.repo.GetAll()
}

func (s *VacancyService) Clear() {
	s.repo.Clear()
}
