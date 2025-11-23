# PR Service — Avito TA For Backend Internship (autumn 2025)

PR Service - это сервис создания и управления Pull Request'ами, при создании Pull Request ему автоматически присваивается до двух reviewers из команды автора.

Cтек:
- Go 1.25.3
- PostgreSQL 18 
- Docker / docker compose для контейнирезации и развертывания
- Prometheus для сбора метрик
- Vegeta для нагрузочного тестирования
- GitHub Actions для CI

## Поднятие сервиса

Склонируйте репозиторий в удобную вам дирикторию.
```bash
git clone https://github.com/vladgrskkh/pr_service .
```

Для запуска сервиса необходим docker. (https://docs.docker.com/engine/install/)

Перед запуском необходимо настроить переменные окружения .env и .env_db.

.env

```env
DB_DSN='postgres://prs_user:pa55word@db/prs?sslmode=disable'
DB_DSN_LOCAL='postgres://prs_user:pa55word@localhost/prs?sslmode=disable'
```

.env_db

```env
POSTGRES_USER='prs_user'
POSTGRES_PASSWORD='pa55word'
POSTGRES_DB='prs'
```

Запуск контейнеров

```bash
make run/docker-compose/up
```

## Makefile

Чтобы посмотреть все Makefile rules
```bash
make help
```

Примерный аутпут:

-Usage:
-  help:                        print this help message
-  run/api:                     run the API application
-  db/psql:                     connect to the database using psql
-  db/migrations/new name=$1:   create a new database migration
-  db/migrations/up:            apply all up database migrations
-  run/docker/api:              run the docker image
-  run/docker-compose/up:       run the docker-compose(docker-compose.yml) stack in detached mode
-  stop/docker-compose/down:   stop the docker-compose(docker-compose.yml) stack
-  audit:                       tidy and vendor dependencies and format, vet and test all code
-  vendor:                     tidy and vendor dependencies
-  build/api:                   build the cmd/api application
-  build/docker:                build the docker image
-  build-and-push/docker:       build the docker image and push it to docker hub
-  e2e/up:                      run the docker-compose(e2e/docker-compose.yml) stack
-  e2e/test:                    run the e2e tests
-  e2e/down:                    stop the docker-compose(e2e/docker-compose.yml) stack
-  e2e:                         run the env and e2e test



---

## Структура проекта

```text
├── cmd
│   └── api
│       └── main.go
├── config
│   ├── config_docker-compose.toml
│   ├── config.go
│   ├── config.toml
│   └── prometheus.yml
├── docker-compose.yml
├── Dockerfile
├── docs
│   └── openapi.yml
├── env
├── go.mod
├── go.sum
├── internal
│   ├── apierrors
│   │   └── apierrors.go
│   ├── application
│   │   └── application.go
│   ├── domain
│   │   ├── pr.go
│   │   ├── team.go
│   │   └── user.go
│   ├── e2e
│   │   └── e2e_test.go
│   ├── handlers
│   │   ├── healthcheck
│   │   │   └── healthcheck.go
│   │   ├── pr
│   │   │   ├── mocks
│   │   │   │   ├── PRMerger.go
│   │   │   │   ├── PullReqCreater.go
│   │   │   │   └── PullReqReassigner.go
│   │   │   └── pr.go
│   │   ├── team
│   │   │   ├── mocks
│   │   │   │   ├── TeamCreater.go
│   │   │   │   └── TeamGetter.go
│   │   │   ├── team_test.go
│   │   │   └── team.go
│   │   └── users
│   │       ├── mocks
│   │       │   ├── IsActiveSetter.go
│   │       │   ├── MassDeactivater.go
│   │       │   └── ReviewsGetter.go
│   │       └── users.go
│   ├── middleware
│   │   ├── metrics.go
│   │   └── recover.go
│   ├── repository
│   │   ├── pullrequest.go
│   │   ├── team.go
│   │   └── users.go
│   ├── server
│   │   └── server.go
│   └── service
│       └── pr.go
├── load - dir с инструментами и результатами load testring
│   ├── attack-30.bin
│   ├── attack-30.json
│   ├── docker-compose.test.yml
│   ├── Makefile
│   ├── synt.sql
│   ├── target.list
│   ├── userSetIsActive.json
│   └── userSetIsActive1.json
├── Makefile
├── migrations
│   ├── 000001_create_teams_table.down.sql
│   ├── 000001_create_teams_table.up.sql
│   ├── 000002_create_users_table.down.sql
│   ├── 000002_create_users_table.up.sql
│   ├── 000003_create_pull_requests_table.down.sql
│   ├── 000003_create_pull_requests_table.up.sql
│   ├── 000004_create_pr_merged_trigger.down.sql
│   ├── 000004_create_pr_merged_trigger.up.sql
│   ├── 000005_create_indexes.down.sql
│   └── 000005_create_indexes.up.sql
├── pkg
│   └── helpers
│       └── json
│           └── json.go
├── README.md
└── vendor

```

---

## API

Endpoints:

Users:
- `GET  /users/getReview` — получить список PR, где пользователь назначен ревьюером.
- `POST /users/setIsActive` — активировать/деактивировать конкретного пользователя.
- `POST /users/massDeactivate` - деактивирует заданных пользователей одной команды и переназначить Pull Requests.

Team:
- `GET  /team/get` — получить команду со списком участников.
- `POST /team/add` — создать команду со списком участников.

PullRequest:
- `POST /pullRequest/create` — создать Pull Request и автоматически назначить до 2 ревьюеров из команды автора.
- `POST /pullRequest/merge` — сделать Pull Request MERGED.
- `POST /pullRequest/reassign` — переназначить ревьюера(меняет указанного ревьювера на активного из команды автора).

System:
- `GET  /health` — healthcheck(status, env, version).
- `GET  /metrics` — метрики Prometheus.

---

## CI

Подключил к проекту Github Actions. При пуше и пул реквесте на ветки main или feat/** запускается пайплайн, который прогоняет все тесты(в том числе e2e) и запускает линтер(golangci-lint).

---

## Реализованные требования

Основная часть:

- Реализован API на Golang в соответствии с документацией openapi.
- БД выбрана PostgreSQL, все миграции содержатся в директории ./migrations, все миграции автоматически применяются при поднятии docker-compose.
- Структура проекта по DDD (domain, service, repository, handlers).
- Метрики собираются при помощи Prometheus (endpoint /metrics).
- Сервис выполняет все требования по SLI.
- Используется docker-compose для поднятия контейнеров(pr_service доступен на порту 8080).

Дополнительная часть:

- Проведено нагрузочное тестирование (vegeta, пробивал два endopoint team/get и users/setIsActive). Отчёт приложен в папке load(attack-30.json).
- Описана конфигурация линтера для golangci-lint
- Реализован endpoint users/massDeactivate, он деактивирует заданных пользователей одной команды и переназначает открытые pr на активных участников команды.
- Реализован сценарий E2E-тестирования(internal/e2e). В этом сценарии создается команда, после создаются PR'ы, один пользователь деактивируется, после merge всех PR'ов.

---

## Вопросы/проблемы и принятые решения
- Изначально возник вопрос такого типа: а может ли один user состоять в нескольких командах. В тз ничего про это сказано не было, поэтому принял решение, что не может, user 
привязан к одной команде.
- Для упрощение задачи самому себе сделал так, что в table pull_requests хранится []text ревьюверов, те без отдельной таблицы. Так как ревьеверов максимум 2, думаю, что сильно на
производительность это не повлияет.
- По endpoint массовой деактивации было непонятно условие, всех ли пользователей мы деактивируем в команде и если всех, тогда смысла переназначать pr нет. Было принято решение
деактивировать часть пользователей и после перераспределять pr. 
- В предложенной документации иногда непонятно почему возвращаются те или иные коды ошибок, 500 вообще не предусмотрена. Принял решение добавить 500, а остальное не трогать.
- Не успел покрыть код unit тестами(и не написал нормальных e2e и intergration), документация тоже оставляет желать лучшего.
