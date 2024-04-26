package sse

import (
	"context"
	"net/http"
)

type RequestCreateFactory func(ctx context.Context, stream string, cli *Client) (*http.Request, error)

func DefaultRequestFactory(ctx context.Context, stream string, cli *Client) (*http.Request, error) {
	req, err := http.NewRequest("GET", cli.URL, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	// Setup request, specify stream to connect to
	if stream != "" {
		query := req.URL.Query()
		query.Add("stream", stream)
		req.URL.RawQuery = query.Encode()
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	lastID, exists := cli.LastEventID.Load().([]byte)
	if exists && lastID != nil {
		req.Header.Set("Last-Event-ID", string(lastID))
	}

	// Add user specified headers
	for k, v := range cli.Headers {
		req.Header.Set(k, v)
	}
	return req, nil
}
