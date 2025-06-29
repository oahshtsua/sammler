package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/oahshtsua/sammler/internal/data"
	"github.com/oahshtsua/sammler/internal/syndication"
)

type Result struct {
	feedID    int64
	feedTitle string
	entries   []syndication.FeedEntry
	err       error
}

func (app *application) refreshFeeds() error {
	feeds, err := app.queries.GetFeeds(context.Background())
	if err != nil {
		app.logger.Error("Failed to get feeds from database.")
		return err
	}

	app.logger.Info("Starting feed refresh", "feed_count", len(feeds))

	var wg sync.WaitGroup
	resultChan := make(chan Result)

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for _, feed := range feeds {
		wg.Add(1)
		go worker(feed, &wg, resultChan)
	}

	successCount := 0
	errorCount := 0
	for result := range resultChan {
		if result.err != nil {
			app.logger.Error("Fetching feed failed",
				"feed_title", result.feedTitle,
				"error", result.err,
			)
			errorCount++
			continue
		}

		now := time.Now().UTC().Format(time.RFC3339)
		if len(result.entries) > 0 {
			entries := buildCreateEntryParams(result.feedID, now, result.entries)
			err = app.queries.CreateMultipleEntry(context.Background(), entries)
			if err != nil {
				app.logger.Error("Fetching feed failed",
					"feed_title", result.feedTitle,
					"entry_count", len(result.entries),
					"error", result.err,
				)
				errorCount++
				continue
			}
		}
		err = app.queries.UpdateFeedCheckedAt(context.Background(), data.UpdateFeedCheckedAtParams{
			ID:        result.feedID,
			CheckedAt: now,
		})
		if err != nil {
			log.Println("Error updating feed checked timestamp", result.feedID, err)
			app.logger.Error("Error updating feed checked timestamp",
				"feed_title", result.feedTitle,
				"error", result.err,
			)
			continue
		}
		app.logger.Info("Successfully updated feed",
			"feed_title", result.feedTitle,
			"new_entries", len(result.entries))
		successCount++
	}
	return nil
}

func worker(f data.Feed, wg *sync.WaitGroup, rc chan Result) {
	defer wg.Done()
	entries, err := syndication.GetNewEntries(f.FeedUrl, f.Type, f.CheckedAt)
	rc <- Result{
		feedID:    f.ID,
		feedTitle: f.Title,
		entries:   entries,
		err:       err,
	}
}
