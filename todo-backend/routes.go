package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rifqoi/aws-project/todo-backend/common/logs"
	"go.uber.org/zap"
)

type Routes struct {
	mux *chi.Mux
}

func SetupRoutes(mux *chi.Mux) Routes {
	setupMiddleware(mux)
	addCorsMiddleware(mux)

	return Routes{mux}
}

type routeFunc interface {
	Set(mux *chi.Mux)
}

func (r Routes) AddRoute(routeFunc ...routeFunc) Routes {
	for _, route := range routeFunc {
		route.Set(r.mux)
	}

	return r
}

func setupMiddleware(mux *chi.Mux) {
	log, _ := zap.NewProduction()
	sugar := log.Sugar()

	mux.Use(middleware.RealIP)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.StripSlashes)
	mux.Use(logs.NewMiddleware(sugar))
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.AllowContentType("application/json"))
}

func addCorsMiddleware(mux *chi.Mux) {
	allowedOrigins := "*"

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{allowedOrigins},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	mux.Use(corsMiddleware.Handler)
}

type TodoRoutes struct {
	todoHandlers *TodoHandlers
}

func NewTodoRoutes(todoHandlers *TodoHandlers) *TodoRoutes {
	return &TodoRoutes{
		todoHandlers: todoHandlers,
	}
}

func (t *TodoRoutes) Set(mux *chi.Mux) {
	mux.Post("/task", t.todoHandlers.AddTask)
}
