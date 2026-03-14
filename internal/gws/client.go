package gws

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// Client wraps the gws CLI for making Google Workspace API calls.
type Client struct {
	account string
}

func NewClient(account string) *Client {
	return &Client{account: account}
}

// GetDocument fetches a Google Doc with all tab content.
func (c *Client) GetDocument(docID string) (*Document, error) {
	params := map[string]any{
		"documentId":        docID,
		"includeTabsContent": true,
	}
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshaling params: %w", err)
	}

	args := []string{"docs", "documents", "get", "--params", string(paramsJSON)}
	out, err := c.run(args...)
	if err != nil {
		return nil, err
	}

	var doc Document
	if err := json.Unmarshal(out, &doc); err != nil {
		return nil, fmt.Errorf("parsing document: %w", err)
	}
	return &doc, nil
}

func (c *Client) run(args ...string) ([]byte, error) {
	cmd := exec.Command("gws", args...)
	if c.account != "" {
		cmd.Env = append(cmd.Environ(), "GOOGLE_WORKSPACE_CLI_ACCOUNT="+c.account)
	}

	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("gws command failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("running gws: %w", err)
	}
	return out, nil
}
