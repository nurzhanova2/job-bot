package parser

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/nurzhanova2/job-bot.git/internal/dto"
)

func ParseHabrVacancies() ([]dto.RawVacancyDTO, error) {
	url := "https://career.habr.com/vacancies?q=golang&type=all"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching Habr: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading HTML: %w", err)
	}

	var vacancies []dto.RawVacancyDTO

	doc.Find(".vacancy-card").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".vacancy-card__title").Text()
		company := s.Find(".vacancy-card__company-title").Text()
		location := s.Find(".vacancy-card__meta").Text()
		url, _ := s.Find("a.vacancy-card__title-link").Attr("href")

		vacancies = append(vacancies, dto.RawVacancyDTO{
			Title:       title,
			Company:     company,
			Location:    location,
			URL:         "https://career.habr.com" + url,
			Source:      "habr",
			Description: "", 
			PublishedAt: "", 
		})
	})

	return vacancies, nil
}
