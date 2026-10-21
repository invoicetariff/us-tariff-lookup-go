package tariff

import (
	"encoding/json"
	"testing"
)

func TestMetaGolden(t *testing.T) {
	m := Meta()
	if m.RulesetVersion != "0.3.4" {
		t.Errorf("rulesetVersion = %q, want 0.3.4", m.RulesetVersion)
	}
	if m.DataAsOf != "2026-10-07" {
		t.Errorf("dataAsOf = %q, want 2026-10-07", m.DataAsOf)
	}
	if m.MFNCuratedCodes != 484 {
		t.Errorf("mfnCuratedCodes = %d, want 484", m.MFNCuratedCodes)
	}
	if m.China301Lines != 25 {
		t.Errorf("china301Lines = %d, want 25", m.China301Lines)
	}
	if m.Global301Tiers != 60 {
		t.Errorf("global301Tiers = %d, want 60", m.Global301Tiers)
	}
	if m.Section232Measures != 8 {
		t.Errorf("section232Measures = %d, want 8", m.Section232Measures)
	}
	if m.ChangelogEntries != 57 {
		t.Errorf("changelogEntries = %d, want 57", m.ChangelogEntries)
	}
}

func TestSearchEnglish(t *testing.T) {
	hits := SearchHTS("t-shirt", 10)
	if len(hits) == 0 {
		t.Fatal("t-shirt returned no hits")
	}
	if hits[0].Entry.HTS8 != "61091000" {
		t.Errorf("first hit = %s, want 61091000", hits[0].Entry.HTS8)
	}
	if hits[0].MatchedBy != "en" {
		t.Errorf("matchedBy = %s, want en", hits[0].MatchedBy)
	}
	if hits[0].Entry.General != "16.5%" {
		t.Errorf("general = %s, want 16.5%%", hits[0].Entry.General)
	}
}

func TestSearchChinese(t *testing.T) {
	hits := SearchHTS("保温杯", 10)
	if len(hits) == 0 {
		t.Fatal("保温杯 returned no hits")
	}
	if hits[0].Entry.HTS8 != "96170010" {
		t.Errorf("first hit = %s, want 96170010", hits[0].Entry.HTS8)
	}
	if hits[0].MatchedBy != "zh" {
		t.Errorf("matchedBy = %s, want zh", hits[0].MatchedBy)
	}
}

// The npm 0.1.0 package returns 0 hits for multi-word queries (plain substring
// match); the Go implementation ships AND-tokenized search from day one.
func TestSearchMultiWord(t *testing.T) {
	for _, q := range []string{"cotton t-shirt", "steel screws"} {
		hits := SearchHTS(q, 10)
		if len(hits) == 0 {
			t.Errorf("multi-word query %q returned no hits", q)
		}
	}
}

func TestSearchCodePrefix(t *testing.T) {
	hits := SearchHTS("6109", 10)
	if len(hits) == 0 || hits[0].Entry.HTS8 != "61091000" {
		t.Fatalf("6109 should hit 61091000 first, got %+v", hits)
	}
	if hits[0].MatchedBy != "code" {
		t.Errorf("matchedBy = %s, want code", hits[0].MatchedBy)
	}
}

func TestSearchLimitAndEmpty(t *testing.T) {
	if got := SearchHTS("   ", 10); got != nil {
		t.Errorf("empty query should return nil, got %d hits", len(got))
	}
	hits := SearchHTS("steel", 1)
	if len(hits) > 1 {
		t.Errorf("limit=1 returned %d hits", len(hits))
	}
}

func TestChangelogOrderAndContent(t *testing.T) {
	entries := Changelog()
	if len(entries) != 57 {
		t.Fatalf("changelog entries = %d, want 57", len(entries))
	}
	if entries[0].Date < entries[len(entries)-1].Date {
		t.Errorf("entries not sorted date-desc: first %s last %s", entries[0].Date, entries[len(entries)-1].Date)
	}
	found := false
	for _, e := range entries {
		if e.ID == "2026-07-24-global301-relaunch" {
			found = true
			if e.TitleEn == "" {
				t.Error("global301 relaunch entry missing TitleEn")
			}
		}
	}
	if !found {
		t.Error("changelog missing 2026-07-24-global301-relaunch")
	}
}

func TestRawJSONValid(t *testing.T) {
	var v map[string]any
	if err := json.Unmarshal(RulesetSeedRaw(), &v); err != nil {
		t.Errorf("ruleset raw JSON invalid: %v", err)
	}
	if err := json.Unmarshal(ChangelogRaw(), &v); err != nil {
		t.Errorf("changelog raw JSON invalid: %v", err)
	}
}
