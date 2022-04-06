package zendesk

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLink(t *testing.T) {
	mockAPI := newMockAPI(http.MethodGet, "link.json")
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	link, err := client.GetLink(ctx, 1)
	if err != nil {
		t.Fatalf("Failed to get ticket link: %s", err)
	}

	expectedID := int64(1)
	if link.ID != expectedID {
		t.Fatalf("Returned link does not have the expected id %d. Link id is %d", expectedID, link.ID)
	}
}

func TestCreateLink(t *testing.T) {
	mockAPI := newMockAPIWithStatus(http.MethodPost, "link.json", http.StatusCreated)
	client := newTestClient(mockAPI)
	defer mockAPI.Close()

	issueId := int64(123456)
	issueKey := "PC-123456"
	ticketId := int64(1)
	link, err := client.CreateLink(ctx, Link{
		IssueID:  &issueId,
		IssueKey: &issueKey,
		TicketID: &ticketId,
	})
	if err != nil {
		t.Fatalf("Failed to link ticket to jira issue %s", err)
	}

	expectedId := int64(5)
	if link.ID != expectedId {
		t.Fatalf("Returned link does not have expected ID %d. Link id is %d", expectedId, link.ID)
	}
}

func TestDeleteLink(t *testing.T) {
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		w.Write(nil)
	}))

	c := newTestClient(mockAPI)
	err := c.DeleteLink(ctx, 420)
	if err != nil {
		t.Fatalf("Failed to delete link field: %s", err)
	}
}
