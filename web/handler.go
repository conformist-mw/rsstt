package web

import (
	"embed"
	"html/template"
	"net/http"
	"strconv"

	"rsstt/logger"
	"rsstt/repository"
	"rsstt/service"
)

//go:embed templates
var templateFS embed.FS

var tmpl = template.Must(template.ParseFS(templateFS, "templates/index.html"))

type feedRow struct {
	ID          uint
	Title       string
	URL         string
	IsActive    bool
	LastFetched string
}

func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /", handleIndex)
	mux.HandleFunc("POST /feeds/add", handleAddFeed)
	mux.HandleFunc("POST /feeds/{id}/toggle", handleToggle)
	mux.HandleFunc("POST /feeds/{id}/delete", handleDelete)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	feeds := repository.FetchAllFeeds()
	rows := make([]feedRow, len(feeds))
	for i, f := range feeds {
		title := f.Url
		if f.Title != nil && *f.Title != "" {
			title = *f.Title
		}
		lastFetched := "never"
		if f.LastFetchedAt.Valid {
			lastFetched = f.LastFetchedAt.Time.Format("2006-01-02 15:04")
		}
		rows[i] = feedRow{
			ID:          f.ID,
			Title:       title,
			URL:         f.Url,
			IsActive:    f.IsActive,
			LastFetched: lastFetched,
		}
	}
	if err := tmpl.Execute(w, rows); err != nil {
		logger.Log.Errorf("template error: %v", err)
	}
}

func handleAddFeed(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "URL required", http.StatusBadRequest)
		return
	}
	if _, err := service.CreateFeed(url); err != nil {
		http.Error(w, "Failed to fetch feed: "+err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleToggle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	repository.ToggleFeedActive(uint(id))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	repository.DeleteFeed(uint(id))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
