package cmd

import (
	"context"
	"reflect"

	"github.com/openclaw/gogcli/internal/outfmt"
)

type pageFetchFunc[T any] func(pageToken string) ([]T, string, error)

func loadPagedItems[T any](page string, all bool, fetch pageFetchFunc[T]) ([]T, string, error) {
	if all {
		items, err := collectAllPages(page, fetch)
		if err != nil {
			var zero []T
			return zero, "", err
		}
		return items, "", nil
	}
	return fetch(page)
}

// addPagedAggregates injects count/has_more into every paged JSON envelope.
// count is the number of slice-valued payload entries in this page; has_more
// reflects the nextPageToken. Existing keys (per-command overrides) win.
func addPagedAggregates(payload map[string]any) {
	count := 0
	for _, v := range payload {
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Slice && rv.Type().Elem().Kind() != reflect.Uint8 {
			count += rv.Len()
		}
	}
	if _, exists := payload["count"]; !exists {
		payload["count"] = count
	}
	if _, exists := payload["has_more"]; !exists {
		tok, _ := payload["nextPageToken"].(string)
		payload["has_more"] = tok != ""
	}
}

func writePagedJSONResult(ctx context.Context, payload map[string]any, emptyCount int, failEmpty bool) error {
	addPagedAggregates(payload)
	if err := outfmt.WriteJSON(ctx, stdoutWriter(ctx), payload); err != nil {
		return err
	}
	if emptyCount == 0 {
		return failEmptyExit(failEmpty)
	}
	return nil
}
