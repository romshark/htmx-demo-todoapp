package server

import (
	_ "embed"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/a-h/templ"

	"github.com/romshark/htmx-demo-todoapp/repository"
)

//go:embed icon.ico
var fileFaviconIco []byte

//go:embed assets/htmx.js
var fileHTMXJS []byte

//go:embed assets/surreal.js
var fileSurrealJS []byte

func render(w http.ResponseWriter, r *http.Request, c templ.Component, name string) {
	err := c.Render(r.Context(), w)
	if err != nil {
		slog.Error("rendering template",
			slog.Any("err", err),
			slog.String("name", name))
		const code = http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
	}
}

type Server struct {
	mux  *http.ServeMux
	repo *repository.Repository
}

var _ http.Handler = new(Server)

func New(repo *repository.Repository) *Server {
	s := &Server{repo: repo}
	m := http.NewServeMux()

	m.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(fileFaviconIco)
	})
	m.HandleFunc("GET /assets/htmx.js", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(fileHTMXJS)
	})
	m.HandleFunc("GET /assets/surreal.js", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(fileSurrealJS)
	})

	// The following endpoints render navigable pages.
	m.HandleFunc("GET /{$}", s.handleIndex)

	// The following endpoints render HTMX components for partial reloads of frames.
	// Non-HTMX requests are rejected with 400 Bad Request.

	// The matrix parameter "var=all" specifies that this endpoint is supposed to be
	// used only from within the variant "all" of the frame "list" when action "toggle"
	// is triggered.
	m.HandleFunc("POST /{id}/toggle/{$}",
		s.handlePostToggleTodo)
	m.HandleFunc("POST /{$}",
		s.handlePostIndex)

	// This endpoint doesn't use the DELETE method because
	// HTML form can't issue a DELETE request.
	m.HandleFunc("POST /{id}/delete/{$}",
		s.handlePostTodoDelete)

	s.mux = m

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("access",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("query", r.URL.Query().Encode()))
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.FormValue("term")

	list, err := fetchTodos(s.repo, searchTerm)
	if err != nil {
		internalErr(w, err, "getting all todos", slog.Default())
		return
	}

	if searchTerm != "" {
		headersHXReplaceURL(w, fmt.Sprintf("/?term=%s", url.QueryEscape(searchTerm)))
	} else {
		headersHXReplaceURL(w, "/")
	}
	if isHXRequest(r) {
		render(w, r, comList(list, searchTerm), "comList")
		return
	}

	headersNoCache(w)
	render(w, r, pageIndex(list, searchTerm), "pageIndex")
}

func (s *Server) handlePostIndex(w http.ResponseWriter, r *http.Request) {
	if !requireHTMXRequest(w, r) {
		return
	}

	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	if _, err := s.repo.Add(title, false, time.Now()); err != nil {
		internalErr(w, err, "addind new todo", slog.Default())
		return
	}

	redirectIndex(w, r)
}

func (s *Server) handlePostTodoDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.repo.Remove(id); err != nil {
		internalErr(w, err, "removing todo", slog.With(slog.String("id", id)))
		return
	}

	if isHXRequest(r) {
		renderList(w, r, s.repo, r.FormValue("term"))
		return
	}

	redirectIndex(w, r)
}

func (s *Server) handlePostToggleTodo(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.FormValue("term")
	id := r.PathValue("id")
	_, err := s.repo.Toggle(id)
	if err != nil {
		internalErr(w, err, "toggling todo", slog.With(slog.String("id", id)))
	}
	slog.Info("toggled", slog.String("id", id))

	if isHXRequest(r) {
		renderList(w, r, s.repo, searchTerm)
		return
	}

	redirectIndex(w, r)
}

func internalErr(w http.ResponseWriter, err error, msg string, log *slog.Logger) {
	log.Error(msg, slog.Any("err", err))
	const code = http.StatusInternalServerError
	http.Error(w, http.StatusText(code), code)
}

func getPercentDone(todos []repository.Todo) string {
	countDone := 0
	for i := range todos {
		if todos[i].Done {
			countDone++
		}
	}
	f := float64(countDone) / float64(len(todos))
	return fmt.Sprintf("%d", int(f*100))
}

func requireHTMXRequest(_ http.ResponseWriter, _ *http.Request) (ok bool) {
	return true
	// http.Error(w, "not an HTMX request", http.StatusBadRequest)
	// return false
}

func isHXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func headersHXReplaceURL(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Replace-Url", url)
}

func headersNoCache(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

func fetchTodos(
	repo *repository.Repository, searchTerm string,
) ([]repository.Todo, error) {
	if searchTerm == "" {
		todos, err := repo.All()
		if err != nil {
			return nil, fmt.Errorf("getting all todos: %w", err)
		}
		return todos, nil
	}
	todos, err := repo.Find(searchTerm)
	if err != nil {
		return nil, fmt.Errorf("searching todos: %w", err)
	}
	return todos, nil
}

func renderList(
	w http.ResponseWriter, r *http.Request,
	repo *repository.Repository, searchTerm string,
) {
	todos, err := fetchTodos(repo, searchTerm)
	if err != nil {
		internalErr(w, err, "fetching todos", slog.Default())
	}
	render(w, r, comList(todos, searchTerm), "comList")
}

func redirectIndex(w http.ResponseWriter, r *http.Request) {
	u := "/"
	searchTerm := r.FormValue("term")
	if searchTerm != "" {
		u = fmt.Sprintf("/?term=%s", url.QueryEscape(searchTerm))
	}
	http.Redirect(w, r, u, http.StatusSeeOther)
}
