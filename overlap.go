package goensemblrest

import (
	"context"
)

// GetOverlapByID retrieves features that overlap a region defined by the given identifier.
func (c *Client) GetOverlapByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getOverlapById", params, nil, target, opts...)
}

// GetOverlapByRegion retrieves features that overlap a given region.
func (c *Client) GetOverlapByRegion(ctx context.Context, species, region string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"region":  region,
	}
	return c.Call(ctx, "getOverlapByRegion", params, nil, target, opts...)
}

// GetOverlapByTranslation retrieves features related to a specific Translation (e.g. domains, variants).
func (c *Client) GetOverlapByTranslation(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getOverlapByTranslation", params, nil, target, opts...)
}
