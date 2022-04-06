package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Link struct {
	ID       int64   `json:"id,omitempty"`
	IssueID  *int64  `json:"issue_id,omitempty"`
	IssueKey *string `json:"issue_key,omitempty"`
	TicketID *int64  `json:"ticket_id,omitempty"`
	URL      *string `json:"url,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type LinkAPI interface {
	GetLink(ctx context.Context, linkID int64) (Link, error)
	CreateLink(ctx context.Context, link Link) (Link, error)
	DeleteLink(ctx context.Context, link Link)
}

func (z *Client) GetLink(ctx context.Context, linkID int64) (Link, error) {
	var result struct {
		Link Link `json:"link"`
	}

	body, err := z.get(ctx, fmt.Sprintf("/jira/links/%d", linkID))
	if err != nil {
		return Link{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Link{}, err
	}

	return result.Link, err
}

func (z *Client) CreateLink(ctx context.Context, link Link) (Link, error) {
	var data, result struct {
		Link Link `json:"link"`
	}
	data.Link = link

	body, err := z.post(ctx, "/jira/links", data)
	if err != nil {
		return Link{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Link{}, err
	}

	return result.Link, nil
}

func (z *Client) DeleteLink(ctx context.Context, linkID int64) error {
	err := z.delete(ctx, fmt.Sprintf("/jira/links/%d", linkID))
	if err != nil {
		return err
	}
	return nil
}
