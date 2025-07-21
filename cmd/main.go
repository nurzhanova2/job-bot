package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/nurzhanova2/job-bot.git/internal/handler"
	"github.com/nurzhanova2/job-bot.git/internal/repository"
	"github.com/nurzhanova2/job-bot.git/internal/scheduler"
	"github.com/nurzhanova2/job-bot.git/internal/service"
)

func main() {
	log.Println("Загрузка переменных окружения...")
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла:", err)
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	keywords := os.Getenv("KEYWORDS")

	if botToken == "" || chatIDStr == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN или TELEGRAM_CHAT_ID не заданы")
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Некорректный TELEGRAM_CHAT_ID: %v", err)
	}

	// Telegram Bot
	log.Println("Инициализация Telegram-бота...")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Ошибка инициализации Telegram-бота: %v", err)
	}
	bot.Debug = false
	log.Printf("Бот авторизован: %s", bot.Self.UserName)

	// Получить chatID от пользователя
	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates := bot.GetUpdatesChan(u)

		log.Println("Напиши что-нибудь боту в Telegram (чтобы получить chatID)...")

		for update := range updates {
			if update.Message != nil {
				log.Printf("chatID: %d, from: %s", update.Message.Chat.ID, update.Message.From.UserName)
			}
		}
	}()

	log.Println("Инициализация компонентов...")
	repo := repository.NewInMemoryVacancyRepository()
	svc := service.NewVacancyService(repo, parseKeywords(keywords))
	notifier := handler.NewTelegramHandler(bot, svc, chatID)
	sched := scheduler.NewScheduler(svc, notifier)

	log.Println("Запуск планировщика задач...")
	sched.Start()

	log.Println("Сервис запущен. Ожидание задач...")

	// Живой лог каждые 30 секунд
	go func() {
		for {
			log.Println("Живой статус: бот активен, планировщик работает...")
			time.Sleep(30 * time.Second)
		}
	}()

	select {}
}

func parseKeywords(input string) []string {
	if input == "" {
		return nil
	}
	parts := strings.Split(input, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	log.Printf("Ключевые слова для фильтрации: %v", parts)
	return parts
}
