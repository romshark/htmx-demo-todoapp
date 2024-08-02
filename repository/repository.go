package repository

import (
	"fmt"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/blevesearch/bleve"
)

type Todo struct {
	ID      string
	Title   string
	Done    bool
	Created time.Time
}

type Repository struct {
	lock      sync.Mutex
	idCounter uint64
	index     bleve.Index
	todos     []Todo
}

// NewRepository creates a new repository instance.
// Use path="" for in-memory search index.
func NewRepository() (*Repository, error) {
	// Create a new in-memory bleve search index.
	indexMapping := bleve.NewIndexMapping()

	index, err := bleve.NewUsing("",
		indexMapping,
		bleve.Config.DefaultIndexType,
		bleve.Config.DefaultMemKVStore,
		nil)
	if err != nil {
		return nil, fmt.Errorf("creating new bleve index: %w", err)
	}
	return &Repository{index: index}, nil
}

func (s *Repository) findByID(id string) (index int) {
	for i := range s.todos {
		if s.todos[i].ID == id {
			return i
		}
	}
	return -1
}

func (s *Repository) Close() error { return s.index.Close() }

// Len returns the number of todo items stored.
func (s *Repository) Len() int { return len(s.todos) }

// Add adds a new todo item.
func (s *Repository) Add(title string, done bool, now time.Time) (id string, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.idCounter++
	id = strconv.FormatInt(int64(s.idCounter), 16)

	t := Todo{ID: id, Title: title, Done: done, Created: now}
	if err := s.index.Index(id, t); err != nil {
		return "", err
	}
	s.todos = append(s.todos, t)
	return id, nil
}

var ErrNotFound = fmt.Errorf("not found")

// Toggle toggles the "done" field of the given todo.
// Returns ErrNotFound if id isn't found.
func (s *Repository) Toggle(id string) (newState Todo, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	i := s.findByID(id)
	if i < 0 {
		return Todo{}, ErrNotFound
	}
	s.todos[i].Done = !s.todos[i].Done
	return s.todos[i], nil
}

// Remove removes a todo item. No-op if id doesn't exist.
func (s *Repository) Remove(id string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	i := s.findByID(id)
	if i < 0 {
		return nil
	}

	if err := s.index.Delete(id); err != nil {
		return err
	}
	s.todos = append(s.todos[:i], s.todos[i+1:]...)
	return nil
}

// All calls retuens all stored todo sorted by index DESC.
func (s *Repository) All() ([]Todo, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	cp := make([]Todo, len(s.todos))
	copy(cp, s.todos)
	slices.Reverse(cp) // Newest first
	return cp, nil
}

// Find calls fn for every item that matches term,
// unless fn returns an error, in which case Find returns immediately.
func (s *Repository) Find(term string) ([]Todo, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	disj := bleve.NewDisjunctionQuery(
		bleve.NewPrefixQuery(term),
		bleve.NewTermQuery(term),
	)
	req := bleve.NewSearchRequest(disj)
	res, err := s.index.Search(req)
	if err != nil {
		return nil, err
	}
	r := make([]Todo, len(res.Hits))

	for i := range res.Hits {
		r[i] = s.todos[s.findByID(res.Hits[i].ID)]
	}
	return r, nil
}
