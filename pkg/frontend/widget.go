package frontend

import (
	"embed"
	"fmt"
	"html/template"

	logshub "logs-hub-frontend/pkg/client"

	"github.com/labstack/echo/v4"
	"github.com/vmkteam/embedlog"
)

type WidgetManager struct {
	embedlog.Logger

	servicesTmpl *template.Template
	logsTmpl     *template.Template

	client *logshub.Client
}

func NewWidgetManager(logger embedlog.Logger, client *logshub.Client) *WidgetManager {
	return &WidgetManager{Logger: logger, client: client}
}

//go:embed services.html
var f embed.FS

//go:embed logs_service.html
var l embed.FS

var funcMap = template.FuncMap{}

func (wm *WidgetManager) Init() error {
	// services.html
	srvBytes, err := f.ReadFile("services.html")
	if err != nil {
		return fmt.Errorf("read services.html err=%w", err)
	}

	servicesTmpl, err := template.New("services").Funcs(funcMap).Parse(string(srvBytes))
	if err != nil {
		return fmt.Errorf("parse services.html err=%w", err)
	}

	wm.servicesTmpl = servicesTmpl

	// logs_service.html
	logsBytes, err := l.ReadFile("logs_service.html")
	if err != nil {
		return fmt.Errorf("read logs_service.html err=%w", err)
	}

	logsTmpl, err := template.New("logs").Funcs(funcMap).Parse(string(logsBytes))
	if err != nil {
		return fmt.Errorf("parse logs_service.html err=%w", err)
	}

	wm.logsTmpl = logsTmpl

	return nil
}

func (wm *WidgetManager) ServicesHandler(c echo.Context) error {
	services, err := wm.client.Logs.Get(c.Request().Context())

	if err != nil {
		wm.Error(c.Request().Context(), "services get failed", "err", err)
		return err
	}

	// execute template with parsed data
	err = wm.servicesTmpl.Execute(c.Response().Writer, services)
	if err != nil {
		wm.Error(c.Request().Context(), "render widget failed", "err", err)
		return err
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return nil
}

func (wm *WidgetManager) LogsByServiceIDHandler(c echo.Context) error {
	serviceIDParam := c.Param("service_id")
	var serviceID int
	_, err := fmt.Sscanf(serviceIDParam, "%d", &serviceID)
	if err != nil {
		wm.Error(c.Request().Context(), "parse service_id failed", "err", err, "service_id", serviceIDParam)
		return err
	}

	logsService, err := wm.client.Logs.GetLogsByServiceID(c.Request().Context(), serviceID)
	if err != nil {
		wm.Error(c.Request().Context(), "logs get failed", "err", err, "service_id", serviceID)
		return err
	}

	// execute template with parsed data
	err = wm.logsTmpl.Execute(c.Response().Writer, logsService)
	if err != nil {
		wm.Error(c.Request().Context(), "render widget failed", "err", err)
		return err
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return nil
}
