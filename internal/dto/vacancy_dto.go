package dto

// Для передачи внутрь системы после парсинга/API
type RawVacancyDTO struct {
	Title       string
	Company     string
	Location    string
	Salary      string
	Description string
	URL         string
	Source      string
	PublishedAt string
}

// Для Telegram-бота
type VacancyMessageDTO struct {
	Title    string
	Company  string
	Location string
	Salary   string
	URL      string
}
