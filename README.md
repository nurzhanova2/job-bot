# Telegram Job Bot

`Telegram-бот` для автоматического парсинга и отправки свежих вакансий с career.habr.com по заданным ключевым словам.

*Проект создан как pet-проект для практики: архитектуры бота, парсинга, планировщика (cron), работы с API*

### Возможности
- Парсинг вакансий с Habr Career
- Фильтрация по ключевым словам
- Отправка подходящих вакансий в Telegram-чат
- Планировщик на CRON (по умолчанию — каждый час)
- Простая настройка через .env
- Поддержка Docker

### Установка
```python
git clone https://github.com/nurzhanova2/http-scanner.git
cd http-scanner
go mod tidy

// запуск
go run cmd/main.go --input=./data/urls.txt --report=json
```