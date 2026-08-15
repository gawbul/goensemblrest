package goensemblrest

import (
	"context"
)

// GetCafeGeneTreeByID retrieves a cafe tree of the gene tree using the gene tree stable identifier.
func (c *Client) GetCafeGeneTreeByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getCafeGeneTreeById", params, nil, target, opts...)
}

// GetCafeGeneTreeMemberBySymbol retrieves the cafe tree of the gene tree that contains the gene identified by a symbol.
func (c *Client) GetCafeGeneTreeMemberBySymbol(ctx context.Context, species, symbol string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"symbol":  symbol,
	}
	return c.Call(ctx, "getCafeGeneTreeMemberBySymbol", params, nil, target, opts...)
}

// GetCafeGeneTreeMemberByID retrieves the cafe tree of the gene tree that contains the gene / transcript / translation stable identifier in the given species.
func (c *Client) GetCafeGeneTreeMemberByID(ctx context.Context, species, id string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"id":      id,
	}
	return c.Call(ctx, "getCafeGeneTreeMemberById", params, nil, target, opts...)
}

// GetGeneTreeByID retrieves a gene tree for a gene tree stable identifier.
func (c *Client) GetGeneTreeByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getGeneTreeById", params, nil, target, opts...)
}

// GetGeneTreeMemberBySymbol retrieves the gene tree that contains the gene identified by a symbol.
func (c *Client) GetGeneTreeMemberBySymbol(ctx context.Context, species, symbol string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"symbol":  symbol,
	}
	return c.Call(ctx, "getGeneTreeMemberBySymbol", params, nil, target, opts...)
}

// GetGeneTreeMemberByID retrieves the gene tree that contains the gene / transcript / translation stable identifier in the given species.
func (c *Client) GetGeneTreeMemberByID(ctx context.Context, species, id string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"id":      id,
	}
	return c.Call(ctx, "getGeneTreeMemberById", params, nil, target, opts...)
}

// GetAlignmentByRegion retrieves genomic alignments as separate blocks based on a region and species.
func (c *Client) GetAlignmentByRegion(ctx context.Context, species, region string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"region":  region,
	}
	return c.Call(ctx, "getAlignmentByRegion", params, nil, target, opts...)
}

// GetHomologyByID retrieves homology information (orthologs) by species and Ensembl gene id.
func (c *Client) GetHomologyByID(ctx context.Context, species, id string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"id":      id,
	}
	return c.Call(ctx, "getHomologyById", params, nil, target, opts...)
}

// GetHomologyBySymbol retrieves homology information (orthologs) by symbol.
func (c *Client) GetHomologyBySymbol(ctx context.Context, species, symbol string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"symbol":  symbol,
	}
	return c.Call(ctx, "getHomologyBySymbol", params, nil, target, opts...)
}
