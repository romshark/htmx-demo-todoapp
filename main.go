package main

import (
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"time"

	"github.com/romshark/htmx-demo-todoapp/repository"
	"github.com/romshark/htmx-demo-todoapp/server"
)

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fHost := flag.String("host", ":8080", "server host address")
	flag.Parse()

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

	s := server.New(repo)
	slog.Info("listening", slog.String("host", *fHost))
	if err := http.ListenAndServe(*fHost, s); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server error", slog.Any("err", err))
		}
	}
}
