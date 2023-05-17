package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rifqoi/aws-project/todo-backend/config"
)

func main() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8000"
	}

	dyDB := config.NewDynamoDB()
	dyRepo := NewDynamoRepo(dyDB)
	todoService := NewTodoService(dyRepo)
	todoHandlers := NewTodoHandlers(todoService)
	todoRoute := NewTodoRoutes(todoHandlers)

	mux := chi.NewRouter()
	SetupRoutes(mux).AddRoute(todoRoute)

	s := http.Server{
		Addr:        ":" + PORT,
		Handler:     mux,
		ReadTimeout: 2 * time.Second,
	}

	go func() {
		log.Println("Starting server at http://localhost:" + PORT)
		err := s.ListenAndServe()
		if err != nil {
			log.Panicf("error starting server")
		}
	}()

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, os.Kill)

	sig := <-ch
	log.Println("Got signal ", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.Shutdown(ctx)
}
