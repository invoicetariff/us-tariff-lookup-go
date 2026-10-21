# us-tariff-lookup-go

Offline US import-tariff reference data for Go: HTS code search (English + 中文), the stacked ruleset overview and a machine-readable rule-change feed — **data embedded, zero dependencies, no network calls, no API key**.

This is the Go distribution of the same versioned snapshot that powers [invoicetariff.com](https://invoicetariff.com) and the [`us-tariff-lookup`](https://www.npmjs.com/package/us-tariff-lookup) npm package. Results are deterministic against a pinned version: every answer carries a `rulesetVersion` and a `dataAsOf` date.

## Install

```bash
go get github.com/invoicetariff/us-tariff-lookup-go
# CLI (optional):
go install github.com/invoicetariff/us-tariff-lookup-go/cmd/us-tariff@latest
```

## Use

```go
package main

import (
    "fmt"
    tariff "github.com/invoicetariff/us-tariff-lookup-go"
)

func main() {
    meta := tariff.Meta()
    fmt.Println(meta.RulesetVersion, meta.DataAsOf) // 0.3.4 2026-10-07 (monthly rebuilds)

    // HTS search: code prefix, English description or Chinese keyword.
    // Multi-word queries are AND-tokenized across all surfaces.
    for _, hit := range tariff.SearchHTS("cotton t-shirt", 5) {
        fmt.Println(hit.Entry.HTS8, hit.Entry.General, hit.MatchedBy)
        // 61091000 16.5% en
    }

    // Machine-readable rule changes, newest first (published + upcoming).
    for _, e := range tariff.Changelog() {
        fmt.Println(e.Date, e.Layer, e.Status, e.TitleEn)
    }
}
```

CLI:

```bash
us-tariff meta
us-tariff search "steel screws"
us-tariff search 保温杯
us-tariff changelog 5
```

All commands print JSON with `rulesetVersion`, `dataAsOf` and a `disclaimer` attached.

## What's inside (snapshot v0.3.4, as of 2026-10-07)

| Surface | Count | Notes |
|---|---|---|
| Curated MFN codes | 484 | High-frequency HTS 8-digit codes with official USITC general rates (EN + 中文) |
| China Section 301 lines | 25 | Lists 2/3/4A + 2026 four-year-review additions |
| Global 301 tiers | 60 | The 2026 global Section 301 program tiers |
| Section 232 measures | 8 | steel / aluminium / copper / autos / furniture / pharma / semiconductors / timber |
| Changelog entries | 57 | published + upcoming + historical, each sourced officially |

Need the layers this module doesn't model? `tariff.RulesetSeedRaw()` and `tariff.ChangelogRaw()` return the full bundled JSON (fees, 232 measures, IEEPA history, per-entry source URLs) for you to unmarshal directly.

## Scope

v1 ships the **query surface** (search / meta / changelog). The full landed-cost calculation — stacked duty layers, specific/compound rates, MPF/HMF entry fees — lives in the [npm package](https://www.npmjs.com/package/us-tariff-lookup) and on the [web calculator](https://invoicetariff.com/tariff-calculator).

**Honest limits:** the bundled search covers curated high-frequency codes, not the full 5,689-line HTS — the full table ships in the [hs-codes-database](https://github.com/invoicetariff/hs-codes-database) repo. Most rate lines carry `sample`/`unverified` flags: planning-grade, not certified entry data. Estimates for planning only — not legal advice; your customs broker signs the entry.

## Data & updates

The snapshot is rebuilt monthly (and ad-hoc when major actions land) from the open [us-tariff-rates-2026](https://github.com/invoicetariff/us-tariff-rates-2026) repo. This module is re-released with each snapshot; watch the repo or the [Tariff Radar](https://invoicetariff.com/radar) feed for changes.

## License

MIT — see [LICENSE](LICENSE).
