package goensemblrest

import (
	"context"
)

// GetMapCdnaToRegion converts cDNA coordinates to genomic coordinates.
func (c *Client) GetMapCdnaToRegion(ctx context.Context, id, region string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"id":     id,
		"region": region,
	}
	return c.Call(ctx, "getMapCdnaToRegion", params, nil, target, opts...)
}

// GetMapCdsToRegion converts CDS coordinates to genomic coordinates.
func (c *Client) GetMapCdsToRegion(ctx context.Context, id, region string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"id":     id,
		"region": region,
	}
	return c.Call(ctx, "getMapCdsToRegion", params, nil, target, opts...)
}

// GetMapAssemblyOneToTwo converts coordinates of one assembly to another.
func (c *Client) GetMapAssemblyOneToTwo(ctx context.Context, species, asmOne, region, asmTwo string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"asm_one": asmOne,
		"region":  region,
		"asm_two": asmTwo,
	}
	return c.Call(ctx, "getMapAssemblyOneToTwo", params, nil, target, opts...)
}

// GetMapTranslationToRegion converts translation (protein) coordinates to genomic coordinates.
func (c *Client) GetMapTranslationToRegion(ctx context.Context, id, region string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"id":     id,
		"region": region,
	}
	return c.Call(ctx, "getMapTranslationToRegion", params, nil, target, opts...)
}
