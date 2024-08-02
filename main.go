package main

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
)

func main() {
	repo, err := NewRepository()
	if err != nil {
		panic(err)
	}

	// Add some default demo todos
	if _, err = repo.Add("Buy milk", false, time.Now()); err != nil {
		panic(err)
	}
	if _, err = repo.Add("Wash the car", false, time.Now()); err != nil {
		panic(err)
	}
	if _, err = repo.Add("Feed the cat", true, time.Now()); err != nil {
		panic(err)
	}

	s := NewServer(repo)
	hostAddr := ":8080"
	slog.Info("listening", slog.String("host", hostAddr))
	if err := http.ListenAndServe(hostAddr, s); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server error", slog.Any("err", err))
		}
	}
}

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
	repo *Repository
}

var _ http.Handler = new(Server)

func NewServer(repo *Repository) *Server {
	s := &Server{repo: repo}
	m := http.NewServeMux()

	// The following endpoints render navigable pages.
	m.HandleFunc("GET /{$}", s.page_index)

	// The following endpoints render HTMX components for partial reloads of frames.
	// Non-HTMX requests are rejected with 400 Bad Request.
	// Every endpoint represents an action with parameters on a particular frame:
	//
	//   METHOD /hx/<frame_name>[;var=<variant_name>][;act=<action_name>]/...
	m.HandleFunc("GET /hx/list;act=search/{$}",
		s.get_hx_list__search)

	// The matrix parameter "var=all" specifies that this endpoint is supposed to be
	// used only from within the variant "all" of the frame "list" when action "toggle"
	// is triggered.
	m.HandleFunc("POST /hx/list;var=all;act=toggle/{id}/",
		s.post_hx_list_all_toggle)
	m.HandleFunc("POST /hx/list;var=all;act=add/{$}",
		s.post_hx_list_all_add)
	m.HandleFunc("DELETE /hx/list;var=all/{id}/{$}",
		s.delete_hx_list_all)

	// The matrix parameter "var=search-result" specifies that this endpoint is supposed
	// to be used only from within the variant "search-result" of the frame "list".
	m.HandleFunc("POST /hx/list;var=search-result;act=toggle/{id}/{$}",
		s.post_hx_list_searchResult_toggle)
	m.HandleFunc("DELETE /hx/list;var=search-result/{id}/{$}",
		s.delete_hx_list_searchResult)
	s.mux = m

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("access",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path))
	s.mux.ServeHTTP(w, r)
}

func (s *Server) page_index(w http.ResponseWriter, r *http.Request) {
	list, err := s.repo.All()
	if err != nil {
		internalErr(w, err, "getting all todos", slog.Default())
		return
	}
	render(w, r, page_index(list), "pageIndex")
}

func (s *Server) get_hx_list__search(w http.ResponseWriter, r *http.Request) {
	if !requireHTMXRequest(w, r) {
		return
	}

	term := r.FormValue("term")
	if term == "" {
		render_list_all(w, r, s.repo)
		return
	}

	render_list_searchResult(w, r, s.repo)
}

func (s *Server) post_hx_list_all_add(w http.ResponseWriter, r *http.Request) {
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

	render_list_all(w, r, s.repo)
}

func (s *Server) delete_hx_list_all(w http.ResponseWriter, r *http.Request) {
	if !requireHTMXRequest(w, r) {
		return
	}

	id := r.PathValue("id")
	if err := s.repo.Remove(id); err != nil {
		internalErr(w, err,
			"getting all todos",
			slog.With(slog.String("id", id)))
		return
	}

	render_list_all(w, r, s.repo)
}

func (s *Server) post_hx_list_all_toggle(w http.ResponseWriter, r *http.Request) {
	if !requireHTMXRequest(w, r) {
		return
	}

	id := r.PathValue("id")
	_, err := s.repo.Toggle(id)
	if err != nil {
		internalErr(w, err, "toggling", slog.With(slog.String("id", id)))
	}
	slog.Info("toggled", slog.String("id", id))

	render_list_all(w, r, s.repo)
}

func (s *Server) post_hx_list_searchResult_toggle(w http.ResponseWriter, r *http.Request) {
	if !requireHTMXRequest(w, r) {
		return
	}

	id := r.PathValue("id")
	_, err := s.repo.Toggle(id)
	if err != nil {
		internalErr(w, err, "toggling", slog.With(slog.String("id", id)))
	}
	slog.Info("toggled", slog.String("id", id))

	render_list_searchResult(w, r, s.repo)
}

func (s *Server) delete_hx_list_searchResult(w http.ResponseWriter, r *http.Request) {
	if !requireHTMXRequest(w, r) {
		return
	}

	id := r.PathValue("id")
	if err := s.repo.Remove(id); err != nil {
		internalErr(w, err,
			"getting all todos",
			slog.With(slog.String("id", id)))
		return
	}

	render_list_searchResult(w, r, s.repo)
}

func internalErr(w http.ResponseWriter, err error, msg string, log *slog.Logger) {
	log.Error(msg, slog.Any("err", err))
	const code = http.StatusInternalServerError
	http.Error(w, http.StatusText(code), code)
}

func requireHTMXRequest(w http.ResponseWriter, r *http.Request) (ok bool) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "not an HTMX request", http.StatusBadRequest)
		return false
	}
	return true
}

func render_list_all(w http.ResponseWriter, r *http.Request, repo *Repository) {
	todos, err := repo.All()
	if err != nil {
		internalErr(w, err, "getting all todos", slog.Default())
	}
	render(w, r, frame_list_all(todos), "frame_list_all")
}

func render_list_searchResult(w http.ResponseWriter, r *http.Request, repo *Repository) {
	todos, err := repo.Find(r.FormValue("term"))
	if err != nil {
		internalErr(w, err, "getting all todos", slog.Default())
	}
	render(w, r, frame_list_searchResult(todos), "frame_list_searchResult")
}
