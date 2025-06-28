package data

import (
	"context"
	"fmt"
	"strings"
)

func (q *Queries) CreateMultipleEntry(ctx context.Context, args []CreateEntryParams) error {
	baseQuery := `INSERT INTO entries (
	feed_id,
	title,
	author,
	content,
	external_url,
	published_at,
	created_at
	) VALUES`

	placeholders := []string{}
	arguments := []any{}

	for _, arg := range args {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?)")
		arguments = append(arguments, arg.FeedID)
		arguments = append(arguments, arg.Title)
		arguments = append(arguments, arg.Author)
		arguments = append(arguments, arg.Content)
		arguments = append(arguments, arg.ExternalUrl)
		arguments = append(arguments, arg.PublishedAt)
		arguments = append(arguments, arg.CreatedAt)
	}
	finalQuery := fmt.Sprintf("%s %s;", baseQuery, strings.Join(placeholders, ","))
	_, err := q.db.ExecContext(ctx, finalQuery, arguments...)
	return err
}
