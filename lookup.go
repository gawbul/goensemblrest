package goensemblrest

import (
	"context"
)

// GetLookupByID finds the species and database for a single identifier (e.g. gene, transcript, protein).
func (c *Client) GetLookupByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getLookupById", params, nil, target, opts...)
}

// GetLookupByMultipleIDs finds the species and database for several identifiers.
func (c *Client) GetLookupByMultipleIDs(ctx context.Context, ids []string, target any, opts ...RequestOption) error {
	body := map[string]any{"ids": ids}
	return c.Call(ctx, "getLookupByMultipleIds", nil, body, target, opts...)
}

// GetLookupBySymbol finds the species and database for a symbol in a linked external database.
func (c *Client) GetLookupBySymbol(ctx context.Context, species, symbol string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"symbol":  symbol,
	}
	return c.Call(ctx, "getLookupBySymbol", params, nil, target, opts...)
}

// GetLookupByMultipleSymbols finds the species and database for a set of symbols in a linked external database.
func (c *Client) GetLookupByMultipleSymbols(ctx context.Context, species string, symbols []string, target any, opts ...RequestOption) error {
	params := map[string]string{"species": species}
	body := map[string]any{"symbols": symbols}
	return c.Call(ctx, "getLookupByMultipleSymbols", params, body, target, opts...)
}
