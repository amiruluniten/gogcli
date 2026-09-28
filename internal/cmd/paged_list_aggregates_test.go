package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/openclaw/gogcli/internal/outfmt"
)

func TestWritePagedJSONResultAggregates(t *testing.T) {
	for _, tc := range []struct {
		name        string
		payload     map[string]any
		wantCount   int
		wantHasMore bool
	}{
		{"more", map[string]any{"messages": []map[string]any{{"id": "1"}, {"id": "2"}}, "nextPageToken": "tok"}, 2, true},
		{"last", map[string]any{"messages": []map[string]any{{"id": "1"}}, "nextPageToken": ""}, 1, false},
		{"empty", map[string]any{"messages": []map[string]any{}, "nextPageToken": ""}, 0, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			ctx := outfmt.WithMode(context.Background(), outfmt.Mode{JSON: true})
			if err := outfmt.WriteJSON(ctx, &buf, func() map[string]any {
				addPagedAggregates(tc.payload)
				return tc.payload
			}()); err != nil {
				t.Fatal(err)
			}
			var got map[string]any
			if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if got["count"].(float64) != float64(tc.wantCount) {
				t.Fatalf("count = %v, want %d", got["count"], tc.wantCount)
			}
			if got["has_more"].(bool) != tc.wantHasMore {
				t.Fatalf("has_more = %v, want %v", got["has_more"], tc.wantHasMore)
			}
			if _, ok := got["nextPageToken"]; !ok {
				t.Fatal("nextPageToken missing")
			}
		})
	}
}
