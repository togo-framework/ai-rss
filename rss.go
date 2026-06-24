// Package rss is a togo AI data-source plugin: fetch + parse RSS/Atom feeds so
// ai-rag ingest and agents can pull article content. Registers an "ai-rss"
// service on the kernel and a REST endpoint: POST /api/ai/rss {"url":"…"}.
package rss

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mmcdole/gofeed"
	"github.com/togo-framework/togo"
)

// Item is a normalized feed entry.
type Item struct {
	Title     string `json:"title"`
	Link      string `json:"link"`
	Content   string `json:"content"`
	Published string `json:"published"`
}

// Source fetches and parses RSS/Atom/JSON feeds.
type Source struct{ fp *gofeed.Parser }

// New returns a feed Source.
func New() *Source { return &Source{fp: gofeed.NewParser()} }

// Fetch parses the feed at url and returns its items.
func (s *Source) Fetch(ctx context.Context, url string) ([]Item, error) {
	feed, err := s.fp.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, err
	}
	items := make([]Item, 0, len(feed.Items))
	for _, it := range feed.Items {
		content := it.Content
		if content == "" {
			content = it.Description
		}
		pub := it.Published
		if it.PublishedParsed != nil {
			pub = it.PublishedParsed.Format(time.RFC3339)
		}
		items = append(items, Item{Title: it.Title, Link: it.Link, Content: content, Published: pub})
	}
	return items, nil
}

// FromKernel returns the registered Source, or nil.
func FromKernel(k *togo.Kernel) *Source {
	if v, ok := k.Get("ai-rss"); ok {
		if s, ok := v.(*Source); ok {
			return s
		}
	}
	return nil
}

func init() {
	togo.RegisterProviderFunc("ai-rss", togo.PriorityService, func(k *togo.Kernel) error {
		s := New()
		k.Set("ai-rss", s)
		mount(k.Router, s)
		return nil
	})
}

func mount(r chi.Router, s *Source) {
	r.Post("/api/ai/rss", func(w http.ResponseWriter, req *http.Request) {
		var body struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil || body.URL == "" {
			http.Error(w, `{"error":"url required"}`, http.StatusBadRequest)
			return
		}
		items, err := s.Fetch(req.Context(), body.URL)
		if err != nil {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"items": items})
	})
}
