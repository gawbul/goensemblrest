package goensemblrest

import (
	"context"
)

// GetRegulationBindingMatrix returns the specified transcription factor binding matrix.
func (c *Client) GetRegulationBindingMatrix(ctx context.Context, species, bindingMatrix string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species":        species,
		"binding_matrix": bindingMatrix,
	}
	return c.Call(ctx, "getRegulationBindingMatrix", params, nil, target, opts...)
}
