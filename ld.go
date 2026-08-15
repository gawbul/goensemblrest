package goensemblrest

import (
	"context"
)

// GetLdID computes and returns LD values between the given variant and all other variants in a window around it.
func (c *Client) GetLdID(ctx context.Context, species, id, populationName string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species":         species,
		"id":              id,
		"population_name": populationName,
	}
	return c.Call(ctx, "getLdId", params, nil, target, opts...)
}

// GetLdPairwise computes and returns LD values between two given variants.
func (c *Client) GetLdPairwise(ctx context.Context, species, id1, id2 string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"id1":     id1,
		"id2":     id2,
	}
	return c.Call(ctx, "getLdPairwise", params, nil, target, opts...)
}

// GetLdRegion computes and returns LD values between all pairs of variants in the defined region.
func (c *Client) GetLdRegion(ctx context.Context, species, region, populationName string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species":         species,
		"region":          region,
		"population_name": populationName,
	}
	return c.Call(ctx, "getLdRegion", params, nil, target, opts...)
}
