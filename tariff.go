// Package tariff provides offline lookup of US import tariff reference data:
// HTS code search (English + 中文), the stacked ruleset overview and a
// machine-readable changelog — the same versioned snapshot that powers
// https://invoicetariff.com and the us-tariff-lookup npm package.
//
// The data ships inside the module (go:embed), so results are deterministic
// and reproducible: every answer carries a ruleset version and an as-of date,
// and nothing calls the network.
//
// This v1 covers search / meta / changelog. The full landed-cost calculation
// (stacked duty layers, specific/compound rates, MPF/HMF entry fees) lives in
// the npm package (github.com/invoicetariff/us-tariff-lookup) and on the web
// calculator; this Go module intentionally ships only the query surface.
package tariff

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed data/ruleset-seed.json
var seedJSON []byte

//go:embed data/changelog.json
var changelogJSON []byte

// MFNEntry is one curated high-frequency HTS 8-digit code with its official
// USITC Column 1 general rate. Sample flags rows not yet verified line-by-line
// against the Federal Register text.
type MFNEntry struct {
	HTS8          string   `json:"hts8"`
	Description   string   `json:"description"`
	DescriptionZh string   `json:"descriptionZh"`
	General       string   `json:"general"`
	KeywordsZh    []string `json:"keywordsZh"`
	Sample        bool     `json:"sample"`
}

type tier struct {
	Rate json.Number `json:"rate"`
}

type global301Program struct {
	Tiers         []tier `json:"tiers"`
	EffectiveFrom string `json:"effectiveFrom"`
}

type ruleset struct {
	Version    string             `json:"version"`
	AsOf       string             `json:"asOf"`
	MFN        []MFNEntry         `json:"mfn"`
	China301   []json.RawMessage  `json:"china301"`
	Global301  global301Program   `json:"global301"`
	Section232 []json.RawMessage  `json:"section232"`
}

var seed = mustParseSeed()

func mustParseSeed() ruleset {
	var r ruleset
	if err := json.Unmarshal(seedJSON, &r); err != nil {
		panic(fmt.Sprintf("tariff: bundled ruleset-seed.json is invalid: %v", err))
	}
	return r
}

// MetaInfo describes the bundled snapshot. Version and AsOf are the two fields
// to surface next to any number derived from this module.
type MetaInfo struct {
	RulesetVersion    string `json:"rulesetVersion"`
	DataAsOf          string `json:"dataAsOf"`
	MFNCuratedCodes   int    `json:"mfnCuratedCodes"`
	China301Lines     int    `json:"china301Lines"`
	Global301Tiers    int    `json:"global301Tiers"`
	Section232Measures int   `json:"section232Measures"`
	ChangelogEntries  int    `json:"changelogEntries"`
	Disclaimer        string `json:"disclaimer"`
	Homepage          string `json:"homepage"`
	Data              DataLinks `json:"data"`
}

// DataLinks points at the open repos the snapshot is rebuilt from.
type DataLinks struct {
	Ruleset      string `json:"ruleset"`
	FullHTSTable string `json:"fullHtsTable"`
}

// Disclaimer is the exact string to echo to end users of anything built on
// this module.
const Disclaimer = "Estimates for planning only — not legal advice. Many rate lines carry sample/unverified flags; your customs broker signs the entry, not this tool."

// Meta reports what's inside the bundled snapshot.
func Meta() MetaInfo {
	return MetaInfo{
		RulesetVersion:     seed.Version,
		DataAsOf:           seed.AsOf,
		MFNCuratedCodes:    len(seed.MFN),
		China301Lines:      len(seed.China301),
		Global301Tiers:     len(seed.Global301.Tiers),
		Section232Measures: len(seed.Section232),
		ChangelogEntries:   len(Changelog()),
		Disclaimer:         Disclaimer,
		Homepage:           "https://invoicetariff.com",
		Data: DataLinks{
			Ruleset:      "https://github.com/invoicetariff/us-tariff-rates-2026",
			FullHTSTable: "https://github.com/invoicetariff/hs-codes-database",
		},
	}
}

// RulesetSeedRaw returns the full bundled ruleset JSON for consumers that need
// layers this package does not model (fees, 232 measures, IEEPA history, …).
func RulesetSeedRaw() []byte { return seedJSON }

// ChangelogRaw returns the full bundled changelog JSON.
func ChangelogRaw() []byte { return changelogJSON }
