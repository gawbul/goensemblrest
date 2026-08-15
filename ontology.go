package goensemblrest

import (
	"context"
)

// GetAncestorsByID reconstructs the entire ancestry of an ontology term from is_a and part_of relationships.
func (c *Client) GetAncestorsByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getAncestorsById", params, nil, target, opts...)
}

// GetAncestorsChartByID reconstructs the entire ancestry chart of a term.
func (c *Client) GetAncestorsChartByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getAncestorsChartById", params, nil, target, opts...)
}

// GetDescendantsByID finds all terms descended from a given term.
func (c *Client) GetDescendantsByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getDescendantsById", params, nil, target, opts...)
}

// GetOntologyByID searches for an ontological term by its namespaced identifier.
func (c *Client) GetOntologyByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getOntologyById", params, nil, target, opts...)
}

// GetOntologyByName searches for a list of ontological terms by their name.
func (c *Client) GetOntologyByName(ctx context.Context, name string, target any, opts ...RequestOption) error {
	params := map[string]string{"name": name}
	return c.Call(ctx, "getOntologyByName", params, nil, target, opts...)
}

// GetTaxonomyClassificationByID returns the taxonomic classification of a taxon node.
func (c *Client) GetTaxonomyClassificationByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getTaxonomyClassificationById", params, nil, target, opts...)
}

// GetTaxonomyByID searches for a taxonomic term by its identifier or name.
func (c *Client) GetTaxonomyByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getTaxonomyById", params, nil, target, opts...)
}

// GetTaxonomyByName searches for a taxonomic id by a non-scientific name.
func (c *Client) GetTaxonomyByName(ctx context.Context, name string, target any, opts ...RequestOption) error {
	params := map[string]string{"name": name}
	return c.Call(ctx, "getTaxonomyByName", params, nil, target, opts...)
}
