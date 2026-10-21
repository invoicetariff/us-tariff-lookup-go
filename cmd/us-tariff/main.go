// us-tariff is the CLI for the offline US tariff reference data bundled in
// github.com/invoicetariff/us-tariff-lookup-go. All commands print JSON to
// stdout with the ruleset version, data-as-of date and disclaimer attached.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	tariff "github.com/invoicetariff/us-tariff-lookup-go"
)

func usage() {
	fmt.Fprint(os.Stderr, `us-tariff — offline US tariff lookup (data bundled, no network)

Usage:
  us-tariff meta                     what's inside / which data version
  us-tariff search <query>           find HTS codes (EN / 中文 / code prefix)
  us-tariff changelog [limit]        recent rule changes (newest first)

Examples:
  us-tariff search "cotton t-shirt"
  us-tariff search 保温杯
  us-tariff changelog 5
`)
}

type envelope struct {
	RulesetVersion string      `json:"rulesetVersion"`
	DataAsOf       string      `json:"dataAsOf"`
	Disclaimer     string      `json:"disclaimer"`
	Result         interface{} `json:"result"`
}

func printEnvelope(w *os.File, result interface{}) {
	meta := tariff.Meta()
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	_ = enc.Encode(envelope{
		RulesetVersion: meta.RulesetVersion,
		DataAsOf:       meta.DataAsOf,
		Disclaimer:     meta.Disclaimer,
		Result:         result,
	})
}

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		usage()
		os.Exit(2)
	}

	switch args[0] {
	case "meta":
		printEnvelope(os.Stdout, tariff.Meta())
	case "search":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "usage: us-tariff search <query>")
			os.Exit(2)
		}
		limit := 10
		if len(args) >= 3 {
			if n, err := strconv.Atoi(args[2]); err == nil {
				limit = n
			}
		}
		hits := tariff.SearchHTS(strings.Join(args[1:], " "), limit)
		printEnvelope(os.Stdout, map[string]interface{}{"count": len(hits), "hits": hits})
	case "changelog":
		limit := 10
		if len(args) >= 2 {
			if n, err := strconv.Atoi(args[1]); err == nil {
				limit = n
			}
		}
		entries := tariff.Changelog()
		if limit > len(entries) {
			limit = len(entries)
		}
		printEnvelope(os.Stdout, map[string]interface{}{"count": len(entries), "entries": entries[:limit]})
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", args[0])
		usage()
		os.Exit(2)
	}
}
