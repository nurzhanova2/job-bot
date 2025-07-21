package scheduler

import (
	"log"
	"time"

	"github.com/nurzhanova2/job-bot.git/internal/handler"
	"github.com/nurzhanova2/job-bot.git/internal/model"
	"github.com/nurzhanova2/job-bot.git/internal/parser"
	"github.com/nurzhanova2/job-bot.git/internal/service"
	"github.com/robfig/cron/v3"
)

type JobScheduler struct {
	Cron     *cron.Cron
	Service  *service.VacancyService
	Notifier *handler.TelegramHandler
}

func NewScheduler(svc *service.VacancyService, notifier *handler.TelegramHandler) *JobScheduler {
	return &JobScheduler{
		Cron:     cron.New(cron.WithSeconds()), 
		Service:  svc,
		Notifier: notifier,
	}
}

func (j *JobScheduler) Start() {
    j.Cron.AddFunc("0 0 * * * *", j.FetchAndNotify)
    j.Cron.Start()
    log.Println("Планировщик запущен: FetchAndNotify будет вызываться каждый час")
}


func (j *JobScheduler) FetchAndNotify() {
	log.Println("FetchAndNotify: старт выполнения", time.Now())

	rawVacancies, err := parser.ParseHabrVacancies()
	if err != nil {
		log.Println("Ошибка парсинга:", err)
		return
	}

	log.Printf("Найдено сырых вакансий: %d", len(rawVacancies))

	for _, rv := range rawVacancies {
		v := model.Vacancy{
			ID:          rv.URL,
			Title:       rv.Title,
			Company:     rv.Company,
			Location:    rv.Location,
			Salary:      rv.Salary,
			Description: rv.Description,
			URL:         rv.URL,
			Source:      rv.Source,
			PublishedAt: rv.PublishedAt,
		}
		err := j.Service.SaveIfRelevant(v)
		if err != nil {
			log.Println("⚠️ Ошибка при сохранении вакансии:", err)
		}
	}

	filtered := j.Service.GetAll()
	log.Printf("Отправляется %d релевантных вакансий", len(filtered))

	if err := j.Notifier.SendAllVacancies(); err != nil {
		log.Println("Ошибка при отправке вакансий в Telegram:", err)
	}

	j.Service.Clear() // очищаем после отправки
	log.Println("Завершено: список вакансий очищен до следующего запуска")
}
