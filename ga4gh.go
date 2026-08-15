package goensemblrest

import (
	"context"
	"fmt"
)

// GetGA4GHBeacon returns Beacon metadata.
func (c *Client) GetGA4GHBeacon(ctx context.Context, target any, opts ...RequestOption) error {
	return c.Call(ctx, "getGA4GHBeacon", nil, nil, target, opts...)
}

// GetGA4GHBeaconQuery executes a Beacon query for allele information via GET.
func (c *Client) GetGA4GHBeaconQuery(
	ctx context.Context,
	alternateBases, assemblyID, referenceBases, referenceName string,
	start int,
	target any,
	opts ...RequestOption,
) error {
	queryOpts := []RequestOption{
		WithQuery("alternateBases", alternateBases),
		WithQuery("assemblyId", assemblyID),
		WithQuery("referenceBases", referenceBases),
		WithQuery("referenceName", referenceName),
		WithQuery("start", fmt.Sprintf("%d", start)),
	}
	allOpts := append(queryOpts, opts...)
	return c.Call(ctx, "getGA4GHBeaconQuery", nil, nil, target, allOpts...)
}

// PostGA4GHBeaconQuery executes a Beacon query via POST.
func (c *Client) PostGA4GHBeaconQuery(ctx context.Context, bodyData any, target any, opts ...RequestOption) error {
	return c.Call(ctx, "postGA4GHBeaconQuery", nil, bodyData, target, opts...)
}

// GetGA4GHFeaturesByID returns the GA4GH record for a specific sequence feature.
func (c *Client) GetGA4GHFeaturesByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getGA4GHFeaturesById", params, nil, target, opts...)
}

// SearchGA4GHFeatures searches for sequence annotation features in GA4GH format.
func (c *Client) SearchGA4GHFeatures(ctx context.Context, bodyData any, target any, opts ...RequestOption) error {
	return c.Call(ctx, "searchGA4GHFeatures", nil, bodyData, target, opts...)
}

// SearchGA4GHCallset returns sets of genotype calls for specific samples in GA4GH format.
func (c *Client) SearchGA4GHCallset(ctx context.Context, bodyData any, target any, opts ...RequestOption) error {
	return c.Call(ctx, "searchGA4GHCallset", nil, bodyData, target, opts...)
}

// GetGA4GHCallsetByID returns the GA4GH record for a CallSet by ID.
func (c *Client) GetGA4GHCallsetByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getGA4GHCallsetById", params, nil, target, opts...)
}

// SearchGA4GHDatasets searches for datasets in GA4GH format.
func (c *Client) SearchGA4GHDatasets(ctx context.Context, bodyData any, target any, opts ...RequestOption) error {
	return c.Call(ctx, "searchGA4GHDatasets", nil, bodyData, target, opts...)
}

// GetGA4GHDatasetsByID returns a dataset in GA4GH format by ID.
func (c *Client) GetGA4GHDatasetsByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getGA4GHDatasetsById", params, nil, target, opts...)
}

// SearchGA4GHFeaturesets searches for feature sets in GA4GH format.
func (c *Client) SearchGA4GHFeaturesets(ctx context.Context, bodyData any, target any, opts ...RequestOption) error {
	return c.Call(ctx, "searchGA4GHFeaturesets", nil, bodyData, target, opts...)
}

// GetGA4GHFeaturesetsByID returns a feature set by ID.
func (c *Client) GetGA4GHFeaturesetsByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getGA4GHFeaturesetsById", params, nil, target, opts...)
}

// GetGA4GHVariantsByID returns a specific variant by ID.
func (c *Client) GetGA4GHVariantsByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getGA4GHVariantsById", params, nil, target, opts...)
}

// SearchGA4GHVariantAnnotations searches for variant annotations in GA4GH format.
func (c *Client) SearchGA4GHVariantAnnotations(ctx context.Context, bodyData any, target any, opts ...RequestOption) error {
	return c.Call(ctx, "searchGA4GHVariantAnnotations", nil, bodyData, target, opts...)
}

// SearchGA4GHVariants searches for variant calls in GA4GH format.
func (c *Client) SearchGA4GHVariants(ctx context.Context, bodyData any, target any, opts ...RequestOption) error {
	return c.Call(ctx, "searchGA4GHVariants", nil, bodyData, target, opts...)
}

// SearchGA4GHVariantsets searches for variant sets in GA4GH format.
func (c *Client) SearchGA4GHVariantsets(ctx context.Context, bodyData any, target any, opts ...RequestOption) error {
	return c.Call(ctx, "searchGA4GHVariantsets", nil, bodyData, target, opts...)
}

// GetGA4GHVariantsetsByID returns a variant set by ID.
func (c *Client) GetGA4GHVariantsetsByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getGA4GHVariantsetsById", params, nil, target, opts...)
}

// SearchGA4GHReferences searches for reference sequences in GA4GH format.
func (c *Client) SearchGA4GHReferences(ctx context.Context, bodyData any, target any, opts ...RequestOption) error {
	return c.Call(ctx, "searchGA4GHReferences", nil, bodyData, target, opts...)
}

// GetGA4GHReferencesByID returns reference sequence data by ID.
func (c *Client) GetGA4GHReferencesByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getGA4GHReferencesById", params, nil, target, opts...)
}

// SearchGA4GHReferencesets searches for reference sets in GA4GH format.
func (c *Client) SearchGA4GHReferencesets(ctx context.Context, bodyData any, target any, opts ...RequestOption) error {
	return c.Call(ctx, "searchGA4GHReferencesets", nil, bodyData, target, opts...)
}

// GetGA4GHReferencesetsByID returns a reference set by ID.
func (c *Client) GetGA4GHReferencesetsByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getGA4GHReferencesetsById", params, nil, target, opts...)
}

// SearchGA4GHVariantAnnotationsets searches for annotation sets in GA4GH format.
func (c *Client) SearchGA4GHVariantAnnotationsets(ctx context.Context, bodyData any, target any, opts ...RequestOption) error {
	return c.Call(ctx, "searchGA4GHVariantAnnotationsets", nil, bodyData, target, opts...)
}

// GetGA4GHVariantAnnotationsetsByID returns metadata for an annotation set by ID.
func (c *Client) GetGA4GHVariantAnnotationsetsByID(ctx context.Context, id string, target any, opts ...RequestOption) error {
	params := map[string]string{"id": id}
	return c.Call(ctx, "getGA4GHVariantAnnotationsetsById", params, nil, target, opts...)
}
