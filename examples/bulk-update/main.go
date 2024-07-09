package main

import (
	"bytes"
	"cmp"
	"embed"
	_ "embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strconv"
	"sync"
	"text/template"
)

var (
	//go:embed *.gohtml
	templatesFS embed.FS
)

func main() {
	srv := newServer([]Record{
		{Name: "blue"},
		{Name: "green"},
		{Name: "yellow"},
		{Name: "red"},
		{Name: "purple"},
		{Name: "pink"},
		{Name: "teal"},
		{Name: "brown"},
		{Name: "olive"},
	})
	slog.Error("the server closed", "err", http.ListenAndServe(":"+cmp.Or(os.Getenv("PORT"), "8080"), srv.handler()))
}

type Record struct {
	ID     int
	Name   string
	Active bool
}

type server struct {
	templates *template.Template

	storage storage
}

func newServer(records []Record) *server {
	srv := server{
		templates: template.Must(template.New("").ParseFS(templatesFS, "*")),
	}
	srv.storage.insert(records)
	return &srv
}

func (srv *server) handler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", srv.index)
	mux.HandleFunc("PUT /activate", srv.activate)
	mux.HandleFunc("PUT /deactivate", srv.deactivate)
	return mux
}

func (srv *server) index(res http.ResponseWriter, req *http.Request) {
	list := srv.storage.list()
	srv.writePage(res, req, "templates.gohtml", http.StatusOK, list)
}

func (srv *server) activate(res http.ResponseWriter, req *http.Request) {
	srv.setActiveStatus(res, req, true)
}

func (srv *server) deactivate(res http.ResponseWriter, req *http.Request) {
	srv.setActiveStatus(res, req, false)
}

func (srv *server) setActiveStatus(res http.ResponseWriter, req *http.Request, active bool) {
	_ = req.ParseForm()
	ids, err := parseIDs(req.Form["ids"], srv.storage.count())
	if err != nil {
		http.Error(res, "failed to parse ids", http.StatusBadRequest)
		return
	}
	list := srv.storage.update(active, ids...)
	srv.writePage(res, req, "rows", http.StatusOK, list)
}

func parseIDs(values []string, max int) ([]int, error) {
	if len(values) >= max {
		return nil, fmt.Errorf("too many ids, max: %d", max)
	}
	ids := make([]int, 0, len(values))
	for _, idStr := range values {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (srv *server) writePage(res http.ResponseWriter, _ *http.Request, name string, status int, data any) {
	var buf bytes.Buffer
	err := srv.templates.ExecuteTemplate(&buf, name, data)
	if err != nil {
		log.Println(err)
		http.Error(res, "failed to write page", http.StatusInternalServerError)
		return
	}
	res.Header().Set("content-type", "text/html")
	res.WriteHeader(status)
	_, _ = res.Write(buf.Bytes())
}

type storage struct {
	mut  sync.RWMutex
	data []Record
}

func (s *storage) insert(records []Record) {
	s.mut.Lock()
	defer s.mut.Unlock()
	start := len(s.data) + 1
	for i := range records {
		records[i].ID = i + start
	}
	s.data = append(s.data, records...)
}

func (s *storage) list() []Record {
	s.mut.RLock()
	defer s.mut.RUnlock()
	return slices.Clone(s.data)
}

func (s *storage) count() int {
	s.mut.RLock()
	defer s.mut.RUnlock()
	return len(s.data)
}

func (s *storage) update(active bool, ids ...int) []Record {
	s.mut.Lock()
	defer s.mut.Unlock()
	for _, id := range ids {
		i := slices.IndexFunc(s.data, func(record Record) bool {
			return record.ID == id
		})
		if i >= 0 {
			s.data[i].Active = active
		}
	}
	return slices.Clone(s.data)
}
