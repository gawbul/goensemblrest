package goensemblrest

import (
	"context"
)

// GetSequenceByID requests sequence data by a stable identifier.
func (c *Client) GetSequenceByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getSequenceById", params, nil, target, opts...)
}

// GetSequenceByMultipleIDs requests multiple sequences by a list of stable identifiers.
func (c *Client) GetSequenceByMultipleIDs(ctx context.Context, ids []string, target any, opts ...RequestOption) error {
	body := map[string]any{"ids": ids}
	return c.Call(ctx, "getSequenceByMultipleIds", nil, body, target, opts...)
}

// GetSequenceByRegion returns the genomic sequence of the specified region of the given species.
func (c *Client) GetSequenceByRegion(ctx context.Context, species, region string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"region":  region,
	}
	return c.Call(ctx, "getSequenceByRegion", params, nil, target, opts...)
}

// GetSequenceByMultipleRegions requests multiple sequences for a list of regions in a species.
func (c *Client) GetSequenceByMultipleRegions(ctx context.Context, species string, regions []string, target any, opts ...RequestOption) error {
	params := map[string]string{"species": species}
	body := map[string]any{"regions": regions}
	return c.Call(ctx, "getSequenceByMultipleRegions", params, body, target, opts...)
}
