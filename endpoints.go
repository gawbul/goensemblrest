package goensemblrest

import (
	"context"
	"fmt"
	"net/http"
)

// EndpointSpec defines the metadata for an Ensembl REST API endpoint.
type EndpointSpec struct {
	Name           string   `json:"name"`
	Doc            string   `json:"doc"`
	URL            string   `json:"url"`
	Method         string   `json:"method"`
	ContentType    string   `json:"content_type"`
	PostParameters []string `json:"post_parameters,omitempty"`
}

// EndpointsTable contains the full catalog of Ensembl REST API endpoints.
var EndpointsTable = map[string]EndpointSpec{
	// Archive
	"getArchiveById": {
		Name:        "getArchiveById",
		Doc:         "Uses the given identifier to return its latest version",
		URL:         "/archive/id/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getArchiveByMultipleIds": {
		Name:           "getArchiveByMultipleIds",
		Doc:            "Retrieve the latest version for a set of identifiers",
		URL:            "/archive/id",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"id"},
	},

	// Comparative Genomics
	"getCafeGeneTreeById": {
		Name:        "getCafeGeneTreeById",
		Doc:         "Retrieves a cafe tree of the gene tree using the gene tree stable identifier",
		URL:         "/cafe/genetree/id/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getCafeGeneTreeMemberBySymbol": {
		Name:        "getCafeGeneTreeMemberBySymbol",
		Doc:         "Retrieves the cafe tree of the gene tree that contains the gene identified by a symbol",
		URL:         "/cafe/genetree/member/symbol/{{species}}/{{symbol}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getCafeGeneTreeMemberById": {
		Name:        "getCafeGeneTreeMemberById",
		Doc:         "Retrieves the cafe tree of the gene tree that contains the gene / transcript / translation stable identifier in the given species",
		URL:         "/cafe/genetree/member/id/{{species}}/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getGeneTreeById": {
		Name:        "getGeneTreeById",
		Doc:         "Retrieves a gene tree for a gene tree stable identifier",
		URL:         "/genetree/id/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getGeneTreeMemberBySymbol": {
		Name:        "getGeneTreeMemberBySymbol",
		Doc:         "Retrieves the gene tree that contains the gene identified by a symbol",
		URL:         "/genetree/member/symbol/{{species}}/{{symbol}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getGeneTreeMemberById": {
		Name:        "getGeneTreeMemberById",
		Doc:         "Retrieves the gene tree that contains the gene / transcript / translation stable identifier in the given species",
		URL:         "/genetree/member/id/{{species}}/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getAlignmentByRegion": {
		Name:        "getAlignmentByRegion",
		Doc:         "Retrieves genomic alignments as separate blocks based on a region and species",
		URL:         "/alignment/region/{{species}}/{{region}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getHomologyById": {
		Name:        "getHomologyById",
		Doc:         "Retrieves homology information (orthologs) by species and Ensembl gene id",
		URL:         "/homology/id/{{species}}/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getHomologyBySymbol": {
		Name:        "getHomologyBySymbol",
		Doc:         "Retrieves homology information (orthologs) by symbol",
		URL:         "/homology/symbol/{{species}}/{{symbol}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},

	// Cross References
	"getXrefsBySymbol": {
		Name:        "getXrefsBySymbol",
		Doc:         "Looks up an external symbol and returns all Ensembl objects linked to it",
		URL:         "/xrefs/symbol/{{species}}/{{symbol}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getXrefsById": {
		Name:        "getXrefsById",
		Doc:         "Perform lookups of Ensembl Identifiers and retrieve their external references in other databases",
		URL:         "/xrefs/id/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getXrefsByName": {
		Name:        "getXrefsByName",
		Doc:         "Performs a lookup based upon the primary accession or display label of an external reference",
		URL:         "/xrefs/name/{{species}}/{{name}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},

	// Information
	"getInfoAnalysis": {
		Name:        "getInfoAnalysis",
		Doc:         "List the names of analyses involved in generating Ensembl data",
		URL:         "/info/analysis/{{species}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoAssembly": {
		Name:        "getInfoAssembly",
		Doc:         "List the currently available assemblies for a species, along with toplevel sequences, chromosomes and cytogenetic bands",
		URL:         "/info/assembly/{{species}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoAssemblyRegion": {
		Name:        "getInfoAssemblyRegion",
		Doc:         "Returns information about the specified toplevel sequence region for the given species",
		URL:         "/info/assembly/{{species}}/{{region_name}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoBiotypes": {
		Name:        "getInfoBiotypes",
		Doc:         "List the functional classifications of gene models that Ensembl associates with a particular species",
		URL:         "/info/biotypes/{{species}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoBiotypesByGroup": {
		Name:        "getInfoBiotypesByGroup",
		Doc:         "List the properties of biotypes within a group",
		URL:         "/info/biotypes/groups/{{group}}/{{object_type}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoBiotypesByName": {
		Name:        "getInfoBiotypesByName",
		Doc:         "List the properties of biotypes with a given name",
		URL:         "/info/biotypes/name/{{name}}/{{object_type}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoComparaMethods": {
		Name:        "getInfoComparaMethods",
		Doc:         "List all compara analyses available",
		URL:         "/info/compara/methods",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoComparaSpeciesSets": {
		Name:        "getInfoComparaSpeciesSets",
		Doc:         "List all collections of species analysed with the specified compara method",
		URL:         "/info/compara/species_sets/{{methods}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoComparas": {
		Name:        "getInfoComparas",
		Doc:         "Lists all available comparative genomics databases and their data release",
		URL:         "/info/comparas",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoData": {
		Name:        "getInfoData",
		Doc:         "Shows the data releases available on this REST server",
		URL:         "/info/data",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoEgVersion": {
		Name:        "getInfoEgVersion",
		Doc:         "Returns the Ensembl Genomes version of the databases backing this service",
		URL:         "/info/eg_version",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoExternalDbs": {
		Name:        "getInfoExternalDbs",
		Doc:         "Lists all available external sources for a species",
		URL:         "/info/external_dbs/{{species}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoDivisions": {
		Name:        "getInfoDivisions",
		Doc:         "Get list of all Ensembl divisions for which information is available",
		URL:         "/info/divisions",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoGenomesByName": {
		Name:        "getInfoGenomesByName",
		Doc:         "Find information about a given genome",
		URL:         "/info/genomes/{{name}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoGenomesByAccession": {
		Name:        "getInfoGenomesByAccession",
		Doc:         "Find information about genomes containing a specified INSDC accession",
		URL:         "/info/genomes/accession/{{accession}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoGenomesByAssembly": {
		Name:        "getInfoGenomesByAssembly",
		Doc:         "Find information about a genome with a specified assembly",
		URL:         "/info/genomes/assembly/{{assembly_id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoGenomesByDivision": {
		Name:        "getInfoGenomesByDivision",
		Doc:         "Find information about all genomes in a given division",
		URL:         "/info/genomes/division/{{division}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoGenomesByTaxonomy": {
		Name:        "getInfoGenomesByTaxonomy",
		Doc:         "Find information about all genomes beneath a given node of the taxonomy",
		URL:         "/info/genomes/taxonomy/{{taxon_name}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoPing": {
		Name:        "getInfoPing",
		Doc:         "Checks if the service is alive",
		URL:         "/info/ping",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoRest": {
		Name:        "getInfoRest",
		Doc:         "Shows the current version of the Ensembl REST API",
		URL:         "/info/rest",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoSoftware": {
		Name:        "getInfoSoftware",
		Doc:         "Shows the current version of the Ensembl API used by the REST server",
		URL:         "/info/software",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoSpecies": {
		Name:        "getInfoSpecies",
		Doc:         "Lists all available species, their aliases, available adaptor groups and data release",
		URL:         "/info/species",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoVariationBySpecies": {
		Name:        "getInfoVariationBySpecies",
		Doc:         "List the variation sources used in Ensembl for a species",
		URL:         "/info/variation/{{species}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoVariationConsequenceTypes": {
		Name:        "getInfoVariationConsequenceTypes",
		Doc:         "Lists all variant consequence types",
		URL:         "/info/variation/consequence_types",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoVariationPopulationIndividuals": {
		Name:        "getInfoVariationPopulationIndividuals",
		Doc:         "List all individuals for a population from a species",
		URL:         "/info/variation/populations/{{species}}/{{population_name}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getInfoVariationPopulations": {
		Name:        "getInfoVariationPopulations",
		Doc:         "List all populations for a species",
		URL:         "/info/variation/populations/{{species}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},

	// Linkage Disequilibrium
	"getLdId": {
		Name:        "getLdId",
		Doc:         "Computes and returns LD values between the given variant and all other variants in a window centered around the given variant",
		URL:         "/ld/{{species}}/{{id}}/{{population_name}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getLdPairwise": {
		Name:        "getLdPairwise",
		Doc:         "Computes and returns LD values between the given variants",
		URL:         "/ld/{{species}}/pairwise/{{id1}}/{{id2}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getLdRegion": {
		Name:        "getLdRegion",
		Doc:         "Computes and returns LD values between all pairs of variants in the defined region",
		URL:         "/ld/{{species}}/region/{{region}}/{{population_name}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},

	// Lookup
	"getLookupById": {
		Name:        "getLookupById",
		Doc:         "Find the species and database for a single identifier e.g. gene, transcript, protein",
		URL:         "/lookup/id/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getLookupByMultipleIds": {
		Name:           "getLookupByMultipleIds",
		Doc:            "Find the species and database for several identifiers",
		URL:            "/lookup/id",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"ids"},
	},
	"getLookupBySymbol": {
		Name:        "getLookupBySymbol",
		Doc:         "Find the species and database for a symbol in a linked external database",
		URL:         "/lookup/symbol/{{species}}/{{symbol}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getLookupByMultipleSymbols": {
		Name:           "getLookupByMultipleSymbols",
		Doc:            "Find the species and database for a set of symbols in a linked external database",
		URL:            "/lookup/symbol/{{species}}",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"symbols"},
	},

	// Mapping
	"getMapCdnaToRegion": {
		Name:        "getMapCdnaToRegion",
		Doc:         "Convert from cDNA coordinates to genomic coordinates",
		URL:         "/map/cdna/{{id}}/{{region}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getMapCdsToRegion": {
		Name:        "getMapCdsToRegion",
		Doc:         "Convert from CDS coordinates to genomic coordinates",
		URL:         "/map/cds/{{id}}/{{region}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getMapAssemblyOneToTwo": {
		Name:        "getMapAssemblyOneToTwo",
		Doc:         "Convert the coordinates of one assembly to another",
		URL:         "/map/{{species}}/{{asm_one}}/{{region}}/{{asm_two}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getMapTranslationToRegion": {
		Name:        "getMapTranslationToRegion",
		Doc:         "Convert from protein (translation) coordinates to genomic coordinates",
		URL:         "/map/translation/{{id}}/{{region}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},

	// Ontologies and Taxonomy
	"getAncestorsById": {
		Name:        "getAncestorsById",
		Doc:         "Reconstruct the entire ancestry of a term from is_a and part_of relationships",
		URL:         "/ontology/ancestors/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getAncestorsChartById": {
		Name:        "getAncestorsChartById",
		Doc:         "Reconstruct the entire ancestry of a term from is_a and part_of relationships",
		URL:         "/ontology/ancestors/chart/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getDescendantsById": {
		Name:        "getDescendantsById",
		Doc:         "Find all the terms descended from a given term",
		URL:         "/ontology/descendants/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getOntologyById": {
		Name:        "getOntologyById",
		Doc:         "Search for an ontological term by its namespaced identifier",
		URL:         "/ontology/id/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getOntologyByName": {
		Name:        "getOntologyByName",
		Doc:         "Search for a list of ontological terms by their name",
		URL:         "/ontology/name/{{name}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getTaxonomyClassificationById": {
		Name:        "getTaxonomyClassificationById",
		Doc:         "Return the taxonomic classification of a taxon node",
		URL:         "/taxonomy/classification/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getTaxonomyById": {
		Name:        "getTaxonomyById",
		Doc:         "Search for a taxonomic term by its identifier or name",
		URL:         "/taxonomy/id/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getTaxonomyByName": {
		Name:        "getTaxonomyByName",
		Doc:         "Search for a taxonomic id by a non-scientific name",
		URL:         "/taxonomy/name/{{name}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},

	// Overlap
	"getOverlapById": {
		Name:        "getOverlapById",
		Doc:         "Retrieves features that overlap a region defined by the given identifier",
		URL:         "/overlap/id/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getOverlapByRegion": {
		Name:        "getOverlapByRegion",
		Doc:         "Retrieves features that overlap a given region",
		URL:         "/overlap/region/{{species}}/{{region}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getOverlapByTranslation": {
		Name:        "getOverlapByTranslation",
		Doc:         "Retrieve features related to a specific Translation as described by its stable ID",
		URL:         "/overlap/translation/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},

	// Phenotype annotations
	"getPhenotypeByAccession": {
		Name:        "getPhenotypeByAccession",
		Doc:         "Return phenotype annotations for genomic features given a phenotype ontology accession",
		URL:         "/phenotype/accession/{{species}}/{{accession}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getPhenotypeByGene": {
		Name:        "getPhenotypeByGene",
		Doc:         "Return phenotype annotations for a given gene",
		URL:         "/phenotype/gene/{{species}}/{{gene}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getPhenotypeByRegion": {
		Name:        "getPhenotypeByRegion",
		Doc:         "Return phenotype annotations that overlap a given genomic region",
		URL:         "/phenotype/region/{{species}}/{{region}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getPhenotypeByTerm": {
		Name:        "getPhenotypeByTerm",
		Doc:         "Return phenotype annotations for genomic features given a phenotype ontology term",
		URL:         "/phenotype/term/{{species}}/{{term}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},

	// Regulation
	"getRegulationBindingMatrix": {
		Name:        "getRegulationBindingMatrix",
		Doc:         "Return the specified binding matrix",
		URL:         "/species/{{species}}/binding_matrix/{{binding_matrix}}/",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},

	// Sequences
	"getSequenceById": {
		Name:        "getSequenceById",
		Doc:         "Request multiple types of sequence by stable identifier",
		URL:         "/sequence/id/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getSequenceByMultipleIds": {
		Name:           "getSequenceByMultipleIds",
		Doc:            "Request multiple types of sequence by a stable identifier list",
		URL:            "/sequence/id",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"ids"},
	},
	"getSequenceByRegion": {
		Name:        "getSequenceByRegion",
		Doc:         "Returns the genomic sequence of the specified region of the given species",
		URL:         "/sequence/region/{{species}}/{{region}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getSequenceByMultipleRegions": {
		Name:           "getSequenceByMultipleRegions",
		Doc:            "Request multiple types of sequence by a list of regions",
		URL:            "/sequence/region/{{species}}",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"regions"},
	},

	// Transcript Haplotypes
	"getTranscriptHaplotypes": {
		Name:        "getTranscriptHaplotypes",
		Doc:         "Computes observed transcript haplotype sequences based on phased genotype data",
		URL:         "/transcript_haplotypes/{{species}}/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},

	// VEP
	"getVariantConsequencesByHGVSNotation": {
		Name:        "getVariantConsequencesByHGVSNotation",
		Doc:         "Fetch variant consequences based on a HGVS notation",
		URL:         "/vep/{{species}}/hgvs/{{hgvs_notation}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getVariantConsequencesByMultipleHGVSNotations": {
		Name:           "getVariantConsequencesByMultipleHGVSNotations",
		Doc:            "Fetch variant consequences for multiple HGVS notations",
		URL:            "/vep/{{species}}/hgvs/",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"hgvs_notations"},
	},
	"getVariantConsequencesById": {
		Name:        "getVariantConsequencesById",
		Doc:         "Fetch variant consequences based on a variant identifier",
		URL:         "/vep/{{species}}/id/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getVariantConsequencesByMultipleIds": {
		Name:           "getVariantConsequencesByMultipleIds",
		Doc:            "Fetch variant consequences for multiple ids",
		URL:            "/vep/{{species}}/id",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"ids"},
	},
	"getVariantConsequencesByRegion": {
		Name:        "getVariantConsequencesByRegion",
		Doc:         "Fetch variant consequences for a region and allele",
		URL:         "/vep/{{species}}/region/{{region}}/{{allele}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getVariantConsequencesByMultipleRegions": {
		Name:           "getVariantConsequencesByMultipleRegions",
		Doc:            "Fetch variant consequences for multiple regions",
		URL:            "/vep/{{species}}/region",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"variants"},
	},

	// Variation
	"getVariationRecoderById": {
		Name:        "getVariationRecoderById",
		Doc:         "Translate a variant identifier, HGVS notation or genomic SPDI notation to all possible variant IDs, HGVS and genomic SPDI",
		URL:         "/variant_recoder/{{species}}/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getVariationRecoderByMultipleIds": {
		Name:           "getVariationRecoderByMultipleIds",
		Doc:            "Translate a list of variant identifiers, HGVS notations or genomic SPDI notations to all possible variant IDs, HGVS and genomic SPDI",
		URL:            "/variant_recoder/{{species}}",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"ids"},
	},
	"getVariationById": {
		Name:        "getVariationById",
		Doc:         "Uses a variant identifier (e.g. rsID) to return the variation features including optional genotype, phenotype and population data",
		URL:         "/variation/{{species}}/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getVariationByPMCID": {
		Name:        "getVariationByPMCID",
		Doc:         "Uses a variant identifier to return variation features with PMCID data",
		URL:         "/variation/{{species}}/pmcid/{{pmcid}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getVariationByPMID": {
		Name:        "getVariationByPMID",
		Doc:         "Uses a variant identifier to return variation features with PMID data",
		URL:         "/variation/{{species}}/pmid/{{pmid}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getVariationByMultipleIds": {
		Name:           "getVariationByMultipleIds",
		Doc:            "Uses a list of variant identifiers to return variation features",
		URL:            "/variation/{{species}}",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"ids"},
	},

	// GA4GH
	"getGA4GHBeacon": {
		Name:        "getGA4GHBeacon",
		Doc:         "Return Beacon information",
		URL:         "/ga4gh/beacon",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getGA4GHBeaconQuery": {
		Name:        "getGA4GHBeaconQuery",
		Doc:         "Return the Beacon response for allele information",
		URL:         "/ga4gh/beacon/query",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"postGA4GHBeaconQuery": {
		Name:           "postGA4GHBeaconQuery",
		Doc:            "Return the Beacon response for allele information",
		URL:            "/ga4gh/beacon/query",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"alternateBases", "assemblyId", "end", "referenceBases", "referenceName", "start", "variantType"},
	},
	"getGA4GHFeaturesById": {
		Name:        "getGA4GHFeaturesById",
		Doc:         "Return the GA4GH record for a specific sequence feature given its identifier",
		URL:         "/ga4gh/features/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"searchGA4GHFeatures": {
		Name:           "searchGA4GHFeatures",
		Doc:            "Return a list of sequence annotation features in GA4GH format",
		URL:            "/ga4gh/features/search",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"end", "referenceName", "start", "featureSetId", "parentId"},
	},
	"searchGA4GHCallset": {
		Name:           "searchGA4GHCallset",
		Doc:            "Return a list of sets of genotype calls for specific samples in GA4GH format",
		URL:            "/ga4gh/callsets/search",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"variantSetId", "name", "pageToken", "pageSize"},
	},
	"getGA4GHCallsetById": {
		Name:        "getGA4GHCallsetById",
		Doc:         "Return the GA4GH record for a specific CallSet given its identifier",
		URL:         "/ga4gh/callsets/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"searchGA4GHDatasets": {
		Name:           "searchGA4GHDatasets",
		Doc:            "Return a list of datasets in GA4GH format",
		URL:            "/ga4gh/datasets/search",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"pageToken", "pageSize"},
	},
	"getGA4GHDatasetsById": {
		Name:        "getGA4GHDatasetsById",
		Doc:         "Return the GA4GH record for a specific dataset given its identifier",
		URL:         "/ga4gh/datasets/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"searchGA4GHFeaturesets": {
		Name:           "searchGA4GHFeaturesets",
		Doc:            "Return a list of feature sets in GA4GH format",
		URL:            "/ga4gh/featuresets/search",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"datasetId", "pageToken", "pageSize"},
	},
	"getGA4GHFeaturesetsById": {
		Name:        "getGA4GHFeaturesetsById",
		Doc:         "Return the GA4GH record for a specific featureSet given its identifier",
		URL:         "/ga4gh/featuresets/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"getGA4GHVariantsById": {
		Name:        "getGA4GHVariantsById",
		Doc:         "Return the GA4GH record for a specific variant given its identifier",
		URL:         "/ga4gh/variants/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"searchGA4GHVariantAnnotations": {
		Name:           "searchGA4GHVariantAnnotations",
		Doc:            "Return variant annotation information in GA4GH format for a region on a reference sequence",
		URL:            "/ga4gh/variantannotations/search",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"variantAnnotationSetId", "effects", "end", "pageSize", "pageToken", "referenceId", "referenceName", "start"},
	},
	"searchGA4GHVariants": {
		Name:           "searchGA4GHVariants",
		Doc:            "Return variant call information in GA4GH format for a region on a reference sequence",
		URL:            "/ga4gh/variants/search",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"variantSetId", "callSetIds", "referenceName", "start", "end", "pageToken", "pageSize"},
	},
	"searchGA4GHVariantsets": {
		Name:           "searchGA4GHVariantsets",
		Doc:            "Return a list of variant sets in GA4GH format",
		URL:            "/ga4gh/variantsets/search",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"datasetId", "pageToken", "pageSize"},
	},
	"getGA4GHVariantsetsById": {
		Name:        "getGA4GHVariantsetsById",
		Doc:         "Return the GA4GH record for a specific VariantSet given its identifier",
		URL:         "/ga4gh/variantsets/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"searchGA4GHReferences": {
		Name:           "searchGA4GHReferences",
		Doc:            "Return a list of reference sequences in GA4GH format",
		URL:            "/ga4gh/references/search",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"referenceSetId", "md5checksum", "accession", "pageToken", "pageSize"},
	},
	"getGA4GHReferencesById": {
		Name:        "getGA4GHReferencesById",
		Doc:         "Return data for a specific reference in GA4GH format by id",
		URL:         "/ga4gh/references/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"searchGA4GHReferencesets": {
		Name:           "searchGA4GHReferencesets",
		Doc:            "Return a list of reference sets in GA4GH format",
		URL:            "/ga4gh/referencesets/search",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"accession", "pageToken", "pageSize"},
	},
	"getGA4GHReferencesetsById": {
		Name:        "getGA4GHReferencesetsById",
		Doc:         "Return data for a specific reference set in GA4GH format",
		URL:         "/ga4gh/referencesets/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
	"searchGA4GHVariantAnnotationsets": {
		Name:           "searchGA4GHVariantAnnotationsets",
		Doc:            "Return a list of annotation sets in GA4GH format",
		URL:            "/ga4gh/variantannotationsets/search",
		Method:         http.MethodPost,
		ContentType:    DefaultContentType,
		PostParameters: []string{"variantSetId", "pageToken", "pageSize"},
	},
	"getGA4GHVariantAnnotationsetsById": {
		Name:        "getGA4GHVariantAnnotationsetsById",
		Doc:         "Return meta data for a specific annotation set in GA4GH format",
		URL:         "/ga4gh/variantannotationsets/{{id}}",
		Method:      http.MethodGet,
		ContentType: DefaultContentType,
	},
}

// Endpoints returns a map containing the specifications of all available Ensembl REST API endpoints.
func (c *Client) Endpoints() map[string]EndpointSpec {
	table := make(map[string]EndpointSpec, len(EndpointsTable))
	for k, v := range EndpointsTable {
		table[k] = v
	}
	return table
}

// Call executes an API endpoint dynamically by name using the internal endpoint table.
func (c *Client) Call(
	ctx context.Context,
	apiCall string,
	pathParams map[string]string,
	bodyData any,
	target any,
	opts ...RequestOption,
) error {
	spec, ok := EndpointsTable[apiCall]
	if !ok {
		return fmt.Errorf("unknown Ensembl REST API endpoint %q", apiCall)
	}

	reqCfg := applyRequestOptions(spec.ContentType, opts...)
	return c.executeRequest(ctx, spec.Method, spec.URL, pathParams, bodyData, target, reqCfg)
}
