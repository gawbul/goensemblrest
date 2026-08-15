package goensemblrest

import (
	"context"
)

// GetArchiveByID uses the given identifier to return its latest version.
func (c *Client) GetArchiveByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getArchiveById", params, nil, target, opts...)
}

// GetArchiveByMultipleIDs retrieves the latest version for a set of identifiers.
func (c *Client) GetArchiveByMultipleIDs(ctx context.Context, ids []string, target any, opts ...RequestOption) error {
	body := map[string]any{"id": ids}
	return c.Call(ctx, "getArchiveByMultipleIds", nil, body, target, opts...)
}
