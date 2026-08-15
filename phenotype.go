package goensemblrest

import (
	"context"
)

// GetPhenotypeByAccession returns phenotype annotations given a phenotype ontology accession.
func (c *Client) GetPhenotypeByAccession(ctx context.Context, species, accession string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species":   species,
		"accession": accession,
	}
	return c.Call(ctx, "getPhenotypeByAccession", params, nil, target, opts...)
}

// GetPhenotypeByGene returns phenotype annotations for a given gene.
func (c *Client) GetPhenotypeByGene(ctx context.Context, species, gene string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"gene":    gene,
	}
	return c.Call(ctx, "getPhenotypeByGene", params, nil, target, opts...)
}

// GetPhenotypeByRegion returns phenotype annotations that overlap a given genomic region.
func (c *Client) GetPhenotypeByRegion(ctx context.Context, species, region string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"region":  region,
	}
	return c.Call(ctx, "getPhenotypeByRegion", params, nil, target, opts...)
}

// GetPhenotypeByTerm returns phenotype annotations given a phenotype ontology term.
func (c *Client) GetPhenotypeByTerm(ctx context.Context, species, term string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species": species,
		"term":    term,
	}
	return c.Call(ctx, "getPhenotypeByTerm", params, nil, target, opts...)
}
