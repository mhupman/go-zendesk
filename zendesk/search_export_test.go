package zendesk

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestSearchExportUser(t *testing.T) {
	mockAPI := newMockAPI(http.MethodGet, "search_export_user.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	results, _, err := client.SearchExport(ctx, &SearchExportOptions{
	    FilterType: "user",
    })
	if err != nil {
		t.Fatalf("Failed to get search results: %s", err)
	}

	list := results.List()
	if len(list) != 1 {
		t.Fatalf("expected length of sla policies is , but got %d", len(list))
	}

	result, ok := list[0].(User)
	if !ok {
		t.Fatalf("Cannot assert %v as a user", list[0])
	}

	if result.ID != 1234 {
		t.Fatalf("User did not have the expected id %v", result)
	}
}

func TestSearchExportQueryParam(t *testing.T) {
	expected := "query string"
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryString := r.URL.Query().Get("query")
		if queryString != expected {
			t.Fatalf(`Did not get the expect query string: "%s". Was: "%s"`, expected, queryString)
		}
		w.Write(readFixture(filepath.Join(http.MethodGet, "search_export_user.json")))
	}))

	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	opts := SearchExportOptions{
		Query: expected,
	}

	_, _, err := client.SearchExport(ctx, &opts)
	if err != nil {
		t.Fatalf("Received error from search export api")
	}
}
