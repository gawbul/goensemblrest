package goensemblrest

import (
	"context"
)

// GetXrefsBySymbol looks up an external symbol and returns all Ensembl objects linked to it.
func (c *Client) GetXrefsBySymbol(ctx context.Context, species, symbol string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"symbol":  symbol,
	}
	return c.Call(ctx, "getXrefsBySymbol", params, nil, target, opts...)
}

// GetXrefsByID performs lookups of Ensembl Identifiers and retrieves their external references in other databases.
func (c *Client) GetXrefsByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getXrefsById", params, nil, target, opts...)
}

// GetXrefsByName performs a lookup based upon the primary accession or display label of an external reference.
func (c *Client) GetXrefsByName(ctx context.Context, species, name string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"name":    name,
	}
	return c.Call(ctx, "getXrefsByName", params, nil, target, opts...)
}
