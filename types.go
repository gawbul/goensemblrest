package goensemblrest

import (
	"encoding/json"
	"net/url"
)

// RequestOption allows customizing individual API calls (e.g. query params, content type, headers).
type RequestOption func(*RequestConfig)

// RequestConfig holds request-specific configuration.
type RequestConfig struct {
	QueryParams url.Values
	Headers     map[string]string
	ContentType string
}

// WithQuery adds a single query parameter to the request.
func WithQuery(key, value string) RequestOption {
	return func(rc *RequestConfig) {
		if rc.QueryParams == nil {
			rc.QueryParams = make(url.Values)
		}
		rc.QueryParams.Add(key, value)
	}
}

// WithQueryParams sets multiple query parameters from a map.
func WithQueryParams(params map[string]string) RequestOption {
	return func(rc *RequestConfig) {
		if rc.QueryParams == nil {
			rc.QueryParams = make(url.Values)
		}
		for k, v := range params {
			rc.QueryParams.Set(k, v)
		}
	}
}

// WithURLValues sets query parameters directly from url.Values.
func WithURLValues(values url.Values) RequestOption {
	return func(rc *RequestConfig) {
		if rc.QueryParams == nil {
			rc.QueryParams = make(url.Values)
		}
		for k, vs := range values {
			for _, v := range vs {
				rc.QueryParams.Add(k, v)
			}
		}
	}
}

// WithRequestContentType overrides the Content-Type / Accept header for the request.
func WithRequestContentType(contentType string) RequestOption {
	return func(rc *RequestConfig) {
		rc.ContentType = contentType
	}
}

// WithRequestHeader sets a custom header on the individual request.
func WithRequestHeader(key, value string) RequestOption {
	return func(rc *RequestConfig) {
		if rc.Headers == nil {
			rc.Headers = make(map[string]string)
		}
		rc.Headers[key] = value
	}
}

// applyRequestOptions creates a RequestConfig with default Content-Type and applies all options.
func applyRequestOptions(defaultContentType string, opts ...RequestOption) RequestConfig {
	cfg := RequestConfig{
		QueryParams: make(url.Values),
		Headers:     make(map[string]string),
		ContentType: defaultContentType,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return cfg
}

// --- Domain Models ---

// PingResponse represents the response from /info/ping.
type PingResponse struct {
	Ping int `json:"ping"`
}

// ArchiveRecord represents an archived identifier entry.
type ArchiveRecord struct {
	ID        string   `json:"id"`
	Latest    string   `json:"latest,omitempty"`
	Version   int      `json:"version,omitempty"`
	Release   int      `json:"release,omitempty"`
	Assembly  string   `json:"assembly,omitempty"`
	Peptide   string   `json:"peptide,omitempty"`
	Type      string   `json:"type,omitempty"`
	Possible  []string `json:"possible,omitempty"`
	IsCurrent int      `json:"is_current,omitempty"`
}

// LookupRecord represents a genomic feature returned by /lookup endpoints.
type LookupRecord struct {
	ID           string          `json:"id"`
	ObjectType   string          `json:"object_type,omitempty"`
	Species      string          `json:"species,omitempty"`
	DisplayName  string          `json:"display_name,omitempty"`
	Description  string          `json:"description,omitempty"`
	Biotype      string          `json:"biotype,omitempty"`
	SeqRegion    string          `json:"seq_region_name,omitempty"`
	Start        int             `json:"start,omitempty"`
	End          int             `json:"end,omitempty"`
	Strand       int             `json:"strand,omitempty"`
	Source       string          `json:"source,omitempty"`
	Version      int             `json:"version,omitempty"`
	DBType       string          `json:"db_type,omitempty"`
	Assembly     string          `json:"assembly_name,omitempty"`
	CanonicalTx  string          `json:"canonical_transcript,omitempty"`
	Transcripts  []LookupRecord  `json:"Transcript,omitempty"`
	Translations []LookupRecord  `json:"Translation,omitempty"`
	Exons        []LookupRecord  `json:"Exon,omitempty"`
	Extra        json.RawMessage `json:"-"`
}

// SequenceRecord represents sequence data returned by /sequence endpoints.
type SequenceRecord struct {
	ID       string `json:"id,omitempty"`
	Seq      string `json:"seq"`
	Desc     string `json:"desc,omitempty"`
	Query    string `json:"query,omitempty"`
	Version  int    `json:"version,omitempty"`
	Molecule string `json:"molecule,omitempty"`
}

// SpeciesRecord represents species metadata from /info/species.
type SpeciesRecord struct {
	Name                string   `json:"name"`
	DisplayName         string   `json:"display_name"`
	TaxonID             string   `json:"taxon_id"`
	Species             string   `json:"species"`
	CommonName          string   `json:"common_name"`
	Division            string   `json:"division"`
	Release             int      `json:"release"`
	Assembly            string   `json:"assembly"`
	Accession           string   `json:"accession"`
	Aliases             []string `json:"aliases"`
	Groups              []string `json:"groups"`
	Strain              string   `json:"strain,omitempty"`
	StrainType          string   `json:"strain_type,omitempty"`
	IsReference         int      `json:"is_reference,omitempty"`
	HasPanCompara       int      `json:"has_pan_compara,omitempty"`
	HasVariations       int      `json:"has_variations,omitempty"`
	HasPeptideCompara   int      `json:"has_peptide_compara,omitempty"`
	HasGenomeAlignments int      `json:"has_genome_alignments,omitempty"`
	HasSynteny          int      `json:"has_synteny,omitempty"`
}

// SpeciesResponse represents the response containing species records.
type SpeciesResponse struct {
	Species []SpeciesRecord `json:"species"`
}

// AssemblyInfo represents assembly information from /info/assembly.
type AssemblyInfo struct {
	AssemblyDate       string          `json:"assembly_date,omitempty"`
	AssemblyLevel      string          `json:"assembly_level,omitempty"`
	AssemblyAcc        string          `json:"assembly_accession,omitempty"`
	AssemblyName       string          `json:"assembly_name,omitempty"`
	DefaultCoordSystem string          `json:"default_coord_system_version,omitempty"`
	GenebuildDate      string          `json:"genebuild_date,omitempty"`
	GenebuildMethod    string          `json:"genebuild_method,omitempty"`
	GenebuildStart     string          `json:"genebuild_start_date,omitempty"`
	GenebuildVersion   string          `json:"genebuild_version,omitempty"`
	TopLevelRegions    []string        `json:"top_level_region,omitempty"`
	Karyotype          []string        `json:"karyotype,omitempty"`
	CoordinateSystems  json.RawMessage `json:"coord_system_versions,omitempty"`
}

// AssemblyRegionInfo represents region metadata from /info/assembly/:species/:region_name.
type AssemblyRegionInfo struct {
	Length        int    `json:"length"`
	CoordSystem   string `json:"coord_system"`
	IsChromosome  int    `json:"is_chromosome"`
	IsCircular    int    `json:"is_circular"`
	Assembly      string `json:"assembly"`
	SeqRegionName string `json:"seq_region_name"`
}

// HomologyRecord represents homologous relationship data.
type HomologyRecord struct {
	Type     string          `json:"type"`
	Method   string          `json:"method_link_type"`
	Taxonomy string          `json:"taxonomy_level"`
	Source   json.RawMessage `json:"source"`
	Target   json.RawMessage `json:"target"`
	DN       *float64        `json:"dn,omitempty"`
	DS       *float64        `json:"ds,omitempty"`
}

// HomologyResponse represents the response structure for homology endpoints.
type HomologyResponse struct {
	Data []struct {
		ID         string           `json:"id"`
		Homologies []HomologyRecord `json:"homologies"`
		LastCommon string           `json:"last_common_ancestor,omitempty"`
	} `json:"data"`
}

// XrefRecord represents an external reference cross-reference.
type XrefRecord struct {
	ID          string   `json:"id"`
	DBName      string   `json:"dbname"`
	DisplayID   string   `json:"display_id"`
	PrimaryID   string   `json:"primary_id"`
	Description string   `json:"description,omitempty"`
	Synonyms    []string `json:"synonyms,omitempty"`
	InfoType    string   `json:"info_type,omitempty"`
	InfoText    string   `json:"info_text,omitempty"`
}

// VariationRecord represents a genetic variant from /variation endpoints.
type VariationRecord struct {
	Name                 string          `json:"name"`
	Source               string          `json:"source"`
	VarClass             string          `json:"var_class"`
	MostSevere           string          `json:"most_severe_consequence"`
	MAF                  *float64        `json:"MAF,omitempty"`
	MinorAllele          string          `json:"minor_allele,omitempty"`
	Evidence             []string        `json:"evidence,omitempty"`
	Synonyms             []string        `json:"synonyms,omitempty"`
	ClinicalSignificance []string        `json:"clinical_significance,omitempty"`
	Mappings             json.RawMessage `json:"mappings,omitempty"`
	Genotypes            json.RawMessage `json:"genotypes,omitempty"`
	Phenotypes           json.RawMessage `json:"phenotypes,omitempty"`
}

// LDRecord represents Linkage Disequilibrium statistics.
type LDRecord struct {
	Variation1 string   `json:"variation1"`
	Variation2 string   `json:"variation2"`
	DPrime     *float64 `json:"d_prime"`
	R2         *float64 `json:"r2"`
}

// MappingRecord represents genomic coordinate mapping result.
type MappingRecord struct {
	Mappings []struct {
		Original struct {
			SeqRegionName string `json:"seq_region_name"`
			Start         int    `json:"start"`
			End           int    `json:"end"`
			Strand        int    `json:"strand"`
			Assembly      string `json:"assembly,omitempty"`
		} `json:"original"`
		Mapped struct {
			SeqRegionName string `json:"seq_region_name"`
			Start         int    `json:"start"`
			End           int    `json:"end"`
			Strand        int    `json:"strand"`
			Assembly      string `json:"assembly,omitempty"`
			CoordSystem   string `json:"coord_system,omitempty"`
		} `json:"mapped"`
	} `json:"mappings"`
}

// VEPRecord represents a Variant Effect Predictor output.
type VEPRecord struct {
	ID                     string          `json:"id"`
	Input                  string          `json:"input"`
	Assembly               string          `json:"assembly_name"`
	SeqRegion              string          `json:"seq_region_name"`
	Start                  int             `json:"start"`
	End                    int             `json:"end"`
	Strand                 int             `json:"strand"`
	MostSevereConsequence  string          `json:"most_severe_consequence"`
	TranscriptConsequences json.RawMessage `json:"transcript_consequences,omitempty"`
	IntergenicConsequences json.RawMessage `json:"intergenic_consequences,omitempty"`
	ColocatedVariants      json.RawMessage `json:"colocated_variants,omitempty"`
}

// BeaconResponse represents GA4GH Beacon information.
type BeaconResponse struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	APIVersion   string          `json:"apiVersion"`
	Organization json.RawMessage `json:"organization"`
	Description  string          `json:"description,omitempty"`
	Version      string          `json:"version,omitempty"`
	WelcomeURL   string          `json:"welcomeUrl,omitempty"`
	Datasets     json.RawMessage `json:"datasets,omitempty"`
}

// BeaconQueryResponse represents GA4GH Beacon query response.
type BeaconQueryResponse struct {
	BeaconID               string          `json:"beaconId"`
	Exists                 *bool           `json:"exists"`
	AlleleRequest          json.RawMessage `json:"alleleRequest,omitempty"`
	DatasetAlleleResponses json.RawMessage `json:"datasetAlleleResponses,omitempty"`
	Error                  json.RawMessage `json:"error,omitempty"`
}
