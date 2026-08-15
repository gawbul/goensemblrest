package goensemblrest

import (
	"context"
)

// GetVariationRecoderByID translates a variant identifier, HGVS notation, or genomic SPDI notation to all possible variant IDs.
func (c *Client) GetVariationRecoderByID(ctx context.Context, species, id string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"id":      id,
	}
	return c.Call(ctx, "getVariationRecoderById", params, nil, target, opts...)
}

// GetVariationRecoderByMultipleIDs translates a list of variant identifiers, HGVS notations, or SPDI notations.
func (c *Client) GetVariationRecoderByMultipleIDs(ctx context.Context, species string, ids []string, target any, opts ...RequestOption) error {
	params := map[string]string{"species": species}
	body := map[string]any{"ids": ids}
	return c.Call(ctx, "getVariationRecoderByMultipleIds", params, body, target, opts...)
}

// GetVariationByID uses a variant identifier (e.g. rsID) to return variation features.
func (c *Client) GetVariationByID(ctx context.Context, species, id string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"id":      id,
	}
	return c.Call(ctx, "getVariationById", params, nil, target, opts...)
}

// GetVariationByPMCID returns variation features associated with a PMCID.
func (c *Client) GetVariationByPMCID(ctx context.Context, species, pmcid string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"pmcid":   pmcid,
	}
	return c.Call(ctx, "getVariationByPMCID", params, nil, target, opts...)
}

// GetVariationByPMID returns variation features associated with a PubMed ID.
func (c *Client) GetVariationByPMID(ctx context.Context, species, pmid string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"pmid":    pmid,
	}
	return c.Call(ctx, "getVariationByPMID", params, nil, target, opts...)
}

// GetVariationByMultipleIDs uses a list of variant identifiers to return variation features.
func (c *Client) GetVariationByMultipleIDs(ctx context.Context, species string, ids []string, target any, opts ...RequestOption) error {
	params := map[string]string{"species": species}
	body := map[string]any{"ids": ids}
	return c.Call(ctx, "getVariationByMultipleIds", params, body, target, opts...)
}
