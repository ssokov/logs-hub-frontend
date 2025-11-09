package frontend

import (
	"embed"
	"fmt"
	"html/template"

	"logs-hub-frontend/pkg/client"

	"github.com/labstack/echo/v4"
	"github.com/vmkteam/embedlog"
)

type WidgetManager struct {
	embedlog.Logger

	template *template.Template

	client *client.Client
}

func NewWidgetManager(logger embedlog.Logger, client *client.Client) *WidgetManager {
	return &WidgetManager{Logger: logger, client: client}
}

//go:embed main.html
var f embed.FS

var funcMap = template.FuncMap{}

func (wm *WidgetManager) Init() error {
	// parse template
	tmp, err := f.ReadFile("main.html")
	if err != nil {
		return fmt.Errorf("parse docs err=%w", err)
	}
	kpTemplate, err := template.New("main").Funcs(funcMap).Parse(string(tmp))
	if err != nil {
		return fmt.Errorf("parse main template err=%w", err)
	}

	wm.template = kpTemplate

	return nil
}

func (wm *WidgetManager) MainHandler(c echo.Context) error {

	services, err := wm.client.Services.GetAllServices(c.Request().Context())

	if err != nil {
		wm.Error(c.Request().Context(), "services get failed", "err", err)
		return err
	}

	// execute template with parsed data
	err = wm.template.Execute(c.Response().Writer, services)
	if err != nil {
		wm.Error(c.Request().Context(), "render widget failed", "err", err)
		return err
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return nil
}
