package application

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vladgrskkh/pr_service/config"
	"github.com/vladgrskkh/pr_service/internal/handlers/healthcheck"
	"github.com/vladgrskkh/pr_service/internal/handlers/pr"
	"github.com/vladgrskkh/pr_service/internal/handlers/team"
	"github.com/vladgrskkh/pr_service/internal/handlers/users"
	"github.com/vladgrskkh/pr_service/internal/repository"
	"github.com/vladgrskkh/pr_service/internal/service"
)

const (
	LevelTrace = slog.Level(-8)
	LevelFatal = slog.Level(12)
)

var LevelNames = map[slog.Leveler]string{
	LevelTrace: "TRACE",
	LevelFatal: "FATAL",
}

var (
	loggerOpts = &slog.HandlerOptions{
		Level: LevelTrace,

		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := LevelNames[level]
				if !exists {
					levelLabel = level.String()
				}

				a.Value = slog.StringValue(levelLabel)
			}
			return a
		},
	}
)

type Application struct {
	Cfg            *config.Config
	Logger         *slog.Logger
	PullReqService *service.PullReqService
	DB             *pgxpool.Pool
}

func NewAppllication(cfgFile string) *Application {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, loggerOpts))

	cfg, err := config.NewConfig(cfgFile)
	if err != nil {
		logger.Log(context.Background(), LevelFatal, err.Error())
		os.Exit(1)
	}

	dbpool, err := openDB(cfg)
	if err != nil {
		logger.Log(context.Background(), LevelFatal, err.Error())
		os.Exit(1)
	}

	defer dbpool.Close()

	logger.Info("db connection pool established")

	pullReqsRepo := repository.NewPullRequestRepo(dbpool)
	teamsRepo := repository.NewTeamRepository(dbpool)
	usersRepo := repository.NewUsersRepo(dbpool)

	pullReqService := service.NewPullReqService(logger, pullReqsRepo, teamsRepo, usersRepo)

	return &Application{
		Cfg:            cfg,
		Logger:         logger,
		PullReqService: pullReqService,
		DB:             dbpool,
	}
}

func openDB(cfg *config.Config) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration(cfg.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	config.MaxConns = cfg.DB.MaxOpenConns
	config.MaxConnIdleTime = duration

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return pool, nil
}

func (app *Application) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/healthcheck", healthcheck.New(app.Logger, app.Cfg.Env, app.Cfg.Version))

	r.Get("/users/{userID}/getReview", users.NewGetReviewsHandler(app.Logger, app.PullReqService))
	r.Post("/users/setIsActive", users.NewPostSetIsActiveHandler(app.Logger, app.PullReqService))

	r.Get("/team/{teamName}", team.NewGetTeamHandler(app.Logger, app.PullReqService))
	r.Post("/team/add", team.NewPostTeamHandler(app.Logger, app.PullReqService))

	r.Post("/pullRequest/merge", pr.NewPostMergeHandler(app.Logger, app.PullReqService))
	r.Post("/pullRequest/create", pr.NewPostPullReqHandler(app.Logger, app.PullReqService))
	r.Post("/pullRequest/reassign", pr.NewPostReassignHandler(app.Logger, app.PullReqService))

	return r
}
