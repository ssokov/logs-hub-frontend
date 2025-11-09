package app

import (
	"context"

	"logs-hub-frontend/pkg/client"
	"logs-hub-frontend/pkg/frontend"

	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
	"github.com/vmkteam/embedlog"
	"github.com/vmkteam/vfs"
)

type Config struct {
	Database *pg.Options
	Server   struct {
		Host      string
		Port      int
		IsDevel   bool
		EnableVFS bool
	}
	Sentry struct {
		Environment string
		DSN         string
	}
	Client struct {
		Endpoint string
	}
	VFS vfs.Config
}

type App struct {
	embedlog.Logger
	appName string
	cfg     Config
	echo    *echo.Echo

	wm     *frontend.WidgetManager
	client *client.Client
}

func New(appName string, sl embedlog.Logger, cfg Config) *App {
	a := &App{
		appName: appName,
		cfg:     cfg,
		echo:    echo.New(),
		Logger:  sl,
	}

	a.client = client.NewDefaultClient(a.cfg.Client.Endpoint)
	a.wm = frontend.NewWidgetManager(a.Logger, a.client)

	return a
}

// Run is a function that runs application.
func (a *App) Run(ctx context.Context) error {
	a.registerHandlers()

	err := a.wm.Init()
	if err != nil {
		a.Error(ctx, "init widgetManager failed", "err", err)
		return err
	}

	return a.runHTTPServer(ctx, a.cfg.Server.Host, a.cfg.Server.Port)
}
