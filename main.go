package main

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/romshark/htmx-demo-todoapp/config"
	"github.com/romshark/htmx-demo-todoapp/repository"
	"github.com/romshark/htmx-demo-todoapp/server"
)

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	conf := config.MustLoad("config.yaml")

	repo, err := repository.NewRepository()
	panicOnErr(err)

	// Add some default demo todos
	_, err = repo.Add("Buy milk", false, time.Now())
	panicOnErr(err)
	_, err = repo.Add("Wash the car", false, time.Now())
	panicOnErr(err)
	_, err = repo.Add("Feed the cat", true, time.Now())
	panicOnErr(err)
	_, err = repo.Add("Buy more cat food", false, time.Now())
	panicOnErr(err)

	s := server.New(repo, conf)
	slog.Info("listening", slog.String("host", conf.Host))
	if err := http.ListenAndServe(conf.Host, s); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server error", slog.Any("err", err))
		}
	}
}
