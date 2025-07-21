package handler

import (
	"fmt"
	"strings"

	 tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nurzhanova2/job-bot.git/internal/service"
)

type TelegramHandler struct {
	Bot     *tgbotapi.BotAPI
	Service *service.VacancyService
	ChatID  int64
}

func NewTelegramHandler(bot *tgbotapi.BotAPI, svc *service.VacancyService, chatID int64) *TelegramHandler {
	return &TelegramHandler{
		Bot:     bot,
		Service: svc,
		ChatID:  chatID,
	}
}

func (h *TelegramHandler) SendAllVacancies() error {
	vacancies := h.Service.GetAll()

	if len(vacancies) == 0 {
		msg := tgbotapi.NewMessage(h.ChatID, "Нет подходящих вакансий.")
		_, err := h.Bot.Send(msg)
		return err
	}

	for _, v := range vacancies {
		text := fmt.Sprintf("*%s*\n%s\n%s\n%s", escapeMarkdown(v.Title), v.Company, v.Salary, v.URL)
		msg := tgbotapi.NewMessage(h.ChatID, text)
		msg.ParseMode = "Markdown"
		_, err := h.Bot.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}

func escapeMarkdown(text string) string {
	replacer := strings.NewReplacer("_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]", "`", "\\`")
	return replacer.Replace(text)
}
