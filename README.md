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
№ Быстрый запуск

### 1. Клонируй репозиторий
```bash
git clone https://github.com/yourname/alatau-city-bank.git
cd alatau-city-bank
2. Настрой .env
3. docker-compose up --build

API: http://localhost:8080
Swagger: http://localhost:8080/swagger/index.html
