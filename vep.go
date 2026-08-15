package goensemblrest

import (
	"context"
)

// GetVariantConsequencesByHGVSNotation fetches variant consequences based on an HGVS notation.
func (c *Client) GetVariantConsequencesByHGVSNotation(ctx context.Context, species, hgvsNotation string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species":       species,
		"hgvs_notation": hgvsNotation,
	}
	return c.Call(ctx, "getVariantConsequencesByHGVSNotation", params, nil, target, opts...)
}

// GetVariantConsequencesByMultipleHGVSNotations fetches variant consequences for multiple HGVS notations.
func (c *Client) GetVariantConsequencesByMultipleHGVSNotations(ctx context.Context, species string, hgvsNotations []string, target any, opts ...RequestOption) error {
	params := map[string]string{"species": species}
	body := map[string]any{"hgvs_notations": hgvsNotations}
	return c.Call(ctx, "getVariantConsequencesByMultipleHGVSNotations", params, body, target, opts...)
}

// GetVariantConsequencesByID fetches variant consequences based on a variant identifier.
func (c *Client) GetVariantConsequencesByID(ctx context.Context, species, id string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"id":      id,
	}
	return c.Call(ctx, "getVariantConsequencesById", params, nil, target, opts...)
}

// GetVariantConsequencesByMultipleIDs fetches variant consequences for multiple IDs.
func (c *Client) GetVariantConsequencesByMultipleIDs(ctx context.Context, species string, ids []string, target any, opts ...RequestOption) error {
	params := map[string]string{"species": species}
	body := map[string]any{"ids": ids}
	return c.Call(ctx, "getVariantConsequencesByMultipleIds", params, body, target, opts...)
}

// GetVariantConsequencesByRegion fetches variant consequences for a given region and allele.
func (c *Client) GetVariantConsequencesByRegion(ctx context.Context, species, region, allele string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"region":  region,
		"allele":  allele,
	}
	return c.Call(ctx, "getVariantConsequencesByRegion", params, nil, target, opts...)
}

// GetVariantConsequencesByMultipleRegions fetches variant consequences for multiple regions.
func (c *Client) GetVariantConsequencesByMultipleRegions(ctx context.Context, species string, variants []string, target any, opts ...RequestOption) error {
	params := map[string]string{"species": species}
	body := map[string]any{"variants": variants}
	return c.Call(ctx, "getVariantConsequencesByMultipleRegions", params, body, target, opts...)
}
