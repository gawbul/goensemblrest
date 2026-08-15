package goensemblrest

import (
	"context"
)

// GetInfoAnalysis lists the names of analyses involved in generating Ensembl data.
func (c *Client) GetInfoAnalysis(ctx context.Context, species string, target any, opts ...RequestOption) error {
	params := map[string]string{"species": species}
	return c.Call(ctx, "getInfoAnalysis", params, nil, target, opts...)
}

// GetInfoAssembly lists the currently available assemblies for a species.
func (c *Client) GetInfoAssembly(ctx context.Context, species string, target any, opts ...RequestOption) error {
	params := map[string]string{"species": species}
	return c.Call(ctx, "getInfoAssembly", params, nil, target, opts...)
}

// GetInfoAssemblyRegion returns information about the specified sequence region for the given species.
func (c *Client) GetInfoAssemblyRegion(ctx context.Context, species, regionName string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species":     species,
		"region_name": regionName,
	}
	return c.Call(ctx, "getInfoAssemblyRegion", params, nil, target, opts...)
}

// GetInfoBiotypes lists functional classifications of gene models for a species.
func (c *Client) GetInfoBiotypes(ctx context.Context, species string, target any, opts ...RequestOption) error {
	params := map[string]string{"species": species}
	return c.Call(ctx, "getInfoBiotypes", params, nil, target, opts...)
}

// GetInfoBiotypesByGroup lists the properties of biotypes within a group.
func (c *Client) GetInfoBiotypesByGroup(ctx context.Context, group, objectType string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"group":       group,
		"object_type": objectType,
	}
	return c.Call(ctx, "getInfoBiotypesByGroup", params, nil, target, opts...)
}

// GetInfoBiotypesByName lists the properties of biotypes with a given name.
func (c *Client) GetInfoBiotypesByName(ctx context.Context, name, objectType string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"name":        name,
		"object_type": objectType,
	}
	return c.Call(ctx, "getInfoBiotypesByName", params, nil, target, opts...)
}

// GetInfoComparaMethods lists all compara analyses available.
func (c *Client) GetInfoComparaMethods(ctx context.Context, target any, opts ...RequestOption) error {
	return c.Call(ctx, "getInfoComparaMethods", nil, nil, target, opts...)
}

// GetInfoComparaSpeciesSets lists all collections of species analysed with the specified compara method.
func (c *Client) GetInfoComparaSpeciesSets(ctx context.Context, methods string, target any, opts ...RequestOption) error {
	params := map[string]string{"methods": methods}
	return c.Call(ctx, "getInfoComparaSpeciesSets", params, nil, target, opts...)
}

// GetInfoComparas lists all available comparative genomics databases and their data release.
func (c *Client) GetInfoComparas(ctx context.Context, target any, opts ...RequestOption) error {
	return c.Call(ctx, "getInfoComparas", nil, nil, target, opts...)
}

// GetInfoData shows the data releases available on this REST server.
func (c *Client) GetInfoData(ctx context.Context, target any, opts ...RequestOption) error {
	return c.Call(ctx, "getInfoData", nil, nil, target, opts...)
}

// GetInfoEgVersion returns the Ensembl Genomes version of the databases backing this service.
func (c *Client) GetInfoEgVersion(ctx context.Context, target any, opts ...RequestOption) error {
	return c.Call(ctx, "getInfoEgVersion", nil, nil, target, opts...)
}

// GetInfoExternalDbs lists all available external sources for a species.
func (c *Client) GetInfoExternalDbs(ctx context.Context, species string, target any, opts ...RequestOption) error {
	params := map[string]string{"species": species}
	return c.Call(ctx, "getInfoExternalDbs", params, nil, target, opts...)
}

// GetInfoDivisions gets the list of all Ensembl divisions for which information is available.
func (c *Client) GetInfoDivisions(ctx context.Context, target any, opts ...RequestOption) error {
	return c.Call(ctx, "getInfoDivisions", nil, nil, target, opts...)
}

// GetInfoGenomesByName finds information about a given genome by name.
func (c *Client) GetInfoGenomesByName(ctx context.Context, name string, target any, opts ...RequestOption) error {
	params := map[string]string{"name": name}
	return c.Call(ctx, "getInfoGenomesByName", params, nil, target, opts...)
}

// GetInfoGenomesByAccession finds information about genomes containing a specified INSDC accession.
func (c *Client) GetInfoGenomesByAccession(ctx context.Context, accession string, target any, opts ...RequestOption) error {
	params := map[string]string{"accession": accession}
	return c.Call(ctx, "getInfoGenomesByAccession", params, nil, target, opts...)
}

// GetInfoGenomesByAssembly finds information about a genome with a specified assembly.
func (c *Client) GetInfoGenomesByAssembly(ctx context.Context, assemblyID string, target any, opts ...RequestOption) error {
	params := map[string]string{"assembly_id": assemblyID}
	return c.Call(ctx, "getInfoGenomesByAssembly", params, nil, target, opts...)
}

// GetInfoGenomesByDivision finds information about all genomes in a given division.
func (c *Client) GetInfoGenomesByDivision(ctx context.Context, division string, target any, opts ...RequestOption) error {
	params := map[string]string{"division": division}
	return c.Call(ctx, "getInfoGenomesByDivision", params, nil, target, opts...)
}

// GetInfoGenomesByTaxonomy finds information about all genomes beneath a given node of the taxonomy.
func (c *Client) GetInfoGenomesByTaxonomy(ctx context.Context, taxonName string, target any, opts ...RequestOption) error {
	params := map[string]string{"taxon_name": taxonName}
	return c.Call(ctx, "getInfoGenomesByTaxonomy", params, nil, target, opts...)
}

// GetInfoPing checks if the Ensembl REST service is alive.
func (c *Client) GetInfoPing(ctx context.Context, target any, opts ...RequestOption) error {
	return c.Call(ctx, "getInfoPing", nil, nil, target, opts...)
}

// GetInfoRest shows the current version of the Ensembl REST API.
func (c *Client) GetInfoRest(ctx context.Context, target any, opts ...RequestOption) error {
	return c.Call(ctx, "getInfoRest", nil, nil, target, opts...)
}

// GetInfoSoftware shows the current version of the Ensembl API used by the REST server.
func (c *Client) GetInfoSoftware(ctx context.Context, target any, opts ...RequestOption) error {
	return c.Call(ctx, "getInfoSoftware", nil, nil, target, opts...)
}

// GetInfoSpecies lists all available species, aliases, and adaptor groups.
func (c *Client) GetInfoSpecies(ctx context.Context, target any, opts ...RequestOption) error {
	return c.Call(ctx, "getInfoSpecies", nil, nil, target, opts...)
}

// GetInfoVariationBySpecies lists the variation sources used in Ensembl for a species.
func (c *Client) GetInfoVariationBySpecies(ctx context.Context, species string, target any, opts ...RequestOption) error {
	params := map[string]string{"species": species}
	return c.Call(ctx, "getInfoVariationBySpecies", params, nil, target, opts...)
}

// GetInfoVariationConsequenceTypes lists all variant consequence types.
func (c *Client) GetInfoVariationConsequenceTypes(ctx context.Context, target any, opts ...RequestOption) error {
	return c.Call(ctx, "getInfoVariationConsequenceTypes", nil, nil, target, opts...)
}

// GetInfoVariationPopulationIndividuals lists all individuals for a population from a species.
func (c *Client) GetInfoVariationPopulationIndividuals(ctx context.Context, species, populationName string, target any, opts ...RequestOption) error {
	params := map[string]string{
		"species":         species,
		"population_name": populationName,
	}
	return c.Call(ctx, "getInfoVariationPopulationIndividuals", params, nil, target, opts...)
}

// GetInfoVariationPopulations lists all populations for a species.
func (c *Client) GetInfoVariationPopulations(ctx context.Context, species string, target any, opts ...RequestOption) error {
	params := map[string]string{"species": species}
	return c.Call(ctx, "getInfoVariationPopulations", params, nil, target, opts...)
}
