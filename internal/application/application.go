package application

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/vladgrskkh/pr_service/config"
	"github.com/vladgrskkh/pr_service/internal/handlers/healthcheck"
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
}

func NewAppllication(cfgFile string) *Application {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, loggerOpts))

	cfg, err := config.NewConfig(cfgFile)
	if err != nil {
		logger.Log(context.Background(), LevelFatal, err.Error())
		os.Exit(1)
	}

	pullReqService := service.NewPullReqService()

	return &Application{
		Cfg:            cfg,
		Logger:         logger,
		PullReqService: pullReqService,
	}
}

func (app *Application) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/healthcheck", healthcheck.New(app.Logger, app.Cfg.Env, app.Cfg.Version))

	return r
}
