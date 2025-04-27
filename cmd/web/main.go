package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/oahshtsua/sammler/internal/data"
	"github.com/oahshtsua/sammler/internal/syndication"
)

type application struct {
	logger  *slog.Logger
	queries *data.Queries
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	port := flag.Int("port", 3456, "Network port")
	dsn := flag.String("dsn", "sammler.db", "Sqlite database file")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := application{
		logger:  logger,
		queries: data.New(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      app.router(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Info("Refreshing feeds...")

	feeds, err := app.queries.GetFeeds(context.Background())
	if err != nil {
		log.Fatal("Error refreshing feed:", err)
	}

	for _, feed := range feeds {
		newEntries, err := syndication.GetNewEntries(feed)
		if err != nil {
			log.Printf("Error refreshing feed: %s", feed.Title)
			log.Fatal(err)
		}
		now := time.Now().UTC().Format(time.RFC3339)
		for _, entry := range newEntries {
			_, err := app.queries.CreateEntry(context.Background(), data.CreateEntryParams{
				FeedID:      feed.ID,
				Title:       entry.Title,
				Subtitle:    sql.NullString{String: entry.Subtitle, Valid: entry.Subtitle != ""},
				Author:      sql.NullString{String: entry.Author.Name, Valid: entry.Author.Name != ""},
				Content:     entry.Content,
				ExternalUrl: entry.Link.Href,
				PublishedAt: entry.Published,
				CreatedAt:   now,
			})

			if err != nil {
				log.Fatal(err)
			}

		}
		err = app.queries.UpdateFeedCheckedAt(context.Background(), data.UpdateFeedCheckedAtParams{
			ID:        feed.ID,
			CheckedAt: now,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	logger.Info("Starting server", "port", *port)
	err = srv.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)

	}
}
