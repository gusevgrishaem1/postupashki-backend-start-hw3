package httprepository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"postupashki-backend-start-hw3/internal/contracts"
)

type Committer struct {
	url    string
	client *http.Client
}

func NewCommitter(taskServiceURL string) *Committer {
	return &Committer{url: strings.TrimRight(taskServiceURL, "/") + "/commit", client: &http.Client{Timeout: 5 * time.Second}}
}

func (c *Committer) Commit(taskID string, result contracts.Result) error {
	body, err := json.Marshal(struct {
		TaskID string `json:"task_id"`
		contracts.Result
	}{TaskID: taskID, Result: result})
	if err != nil {
		return err
	}
	response, err := c.client.Post(c.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("commit returned %s", response.Status)
	}
	return nil
}
