package app

import (
	"context"
	"log"
	"net/http"
	"pr-service/internal/config"
)

type App struct {
	serviceProvider *serviceProvider
	server          *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	return a.runHTTPServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initHTTPServer(_ context.Context) error {
	router := a.setupRouter()
	a.server = &http.Server{
		Addr:    a.serviceProvider.HTTPConfig().Address(),
		Handler: router,
	}
	return nil
}

func (a *App) setupRouter() http.Handler {
    // Генератор создаёт функцию NewRouter(...)
    return openapi.NewRouter(
        openapi.NewStrictHandler(
            a.serviceProvider.Controllers(),
            nil, // сюда можно передавать middleware from generator
        ),
    )
}

func (a *App) runHTTPServer() error {
	log.Printf("HTTP server is running on %s", a.serviceProvider.HTTPConfig().Address())

	err := a.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	log.Println("Shutting down server")
	return a.server.Shutdown(ctx)
}
