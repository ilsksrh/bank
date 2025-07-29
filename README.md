# Alatau City Bank — Backend API

##  Описание

Проект представляет собой backend-систему для автоматизации банковских операций: управление клиентами, аккаунтами, картами, транзакциями, кредитами, а также аудит изменений. 
Разработан в рамках производственной практики Digital Engineering (Нархоз) с использованием Go и PostgreSQL. 
В системе реализована микроархитектура с разделением на сервисы и хранением бизнес-логики в БД через процедуры.

---

##  Технологии

- **Язык:** Golang 1.22+
- **База данных:** PostgreSQL 15+
- **Фреймворк:** Echo
- **ORM:** GORM
- **Документация API:** Swagger/OpenAPI
- **JWT аутентификация:** Access/Refresh токены
- **Docker:** Compose для поднятия сервиса и БД
- **Миграции:** SQL-скрипты

---

## Структура проекта
.
├── cmd/ # main.go — точка входа
├── pkg/
│ ├── app/ # Инициализация Echo, middleware
│ ├── db/ # Подключение к БД
│ ├── handlers/ # Обработчики HTTP-запросов
│ ├── services/ # Бизнес-логика (вызов SQL-функций)
│ └── utils/ # JWT, логгеры, хелперы
├── migrations/ # SQL-файлы создания таблиц и функций
├── .env # Конфигурация окружения
├── docker-compose.yml # Поднятие API и PostgreSQL
├── go.mod / go.sum # Зависимости
└── README.md # Этот файл



№ Быстрый запуск

### 1. Клонируй репозиторий
```bash
git clone https://github.com/yourname/alatau-city-bank.git
cd alatau-city-bank
2. Настрой .env
3. docker-compose up --build

API: http://localhost:8080
Swagger: http://localhost:8080/swagger/index.html
