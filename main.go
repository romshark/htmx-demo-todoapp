package main

import (
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"time"

	"github.com/romshark/httpsim"
	httpsimconf "github.com/romshark/httpsim/config"

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
	conf := config.MustLoad("config.yml")
	fHTTPSimConfig := flag.String(
		"httpsim-conf",
		"httpsim.yml",
		"httpsim config file",
	)
	flag.Parse()

	repo, err := repository.NewRepository()
	panicOnErr(err)

	// Add some default demo todos.
	_, err = repo.Add("Buy milk", false, time.Now())
	panicOnErr(err)
	_, err = repo.Add("Wash the car", false, time.Now())
	panicOnErr(err)
	_, err = repo.Add("Feed the cat", true, time.Now())
	panicOnErr(err)
	_, err = repo.Add("Buy more cat food", false, time.Now())
	panicOnErr(err)
	_, err = repo.Add("Make search faster", false, time.Now())
	panicOnErr(err)

	s := server.New(repo)

	// Use httpsim middleware for simulating error responses and delays.
	httpsimConf, err := httpsimconf.LoadFile(*fHTTPSimConfig)
	if err != nil {
		panicOnErr(err)
	}
	withHTTPSim := httpsim.NewMiddleware(
		server.WithLog(s), *httpsimConf, httpsim.DefaultSleep, httpsim.DefaultRand,
	)

	slog.Info("listening", slog.String("host", conf.Host))
	if err := http.ListenAndServe(conf.Host, withHTTPSim); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server error", slog.Any("err", err))
		}
	}
}
