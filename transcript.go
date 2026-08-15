package goensemblrest

import (
	"context"
)

// GetTranscriptHaplotypes computes observed transcript haplotype sequences based on phased genotype data.
func (c *Client) GetTranscriptHaplotypes(ctx context.Context, species, id string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"id":      id,
	}
	return c.Call(ctx, "getTranscriptHaplotypes", params, nil, target, opts...)
}
