package tariff

import (
	"encoding/json"
	"time"
)

// ChangelogEntry is one rule-change record: what moved, when, on which layer,
// and whether it is published, upcoming or historical.
type ChangelogEntry struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Layer   string `json:"layer"`
	Status  string `json:"status"` // "published" | "upcoming" | "historical"
	Title   string `json:"title"`   // 中文标题
	TitleEn string `json:"titleEn"` // English title
	Summary string `json:"summary"` // 中文摘要，含 FR 编号与生效细节
}

type changelogFile struct {
	Version string           `json:"version"`
	Entries []ChangelogEntry `json:"entries"`
}

var changelog = mustParseChangelog()

func mustParseChangelog() changelogFile {
	var c changelogFile
	if err := json.Unmarshal(changelogJSON, &c); err != nil {
		panic("tariff: bundled changelog.json is invalid: " + err.Error())
	}
	return c
}

// Changelog returns every bundled rule-change entry, newest date first.
// "published" entries are in force, "upcoming" entries have an effective date
// in the future, "historical" entries are superseded layers kept for refund
// lookups. Entries also carry their official source in the raw JSON (see
// ChangelogRaw).
func Changelog() []ChangelogEntry {
	out := make([]ChangelogEntry, len(changelog.Entries))
	copy(out, changelog.Entries)
	sortByDateDesc(out)
	return out
}

func sortByDateDesc(entries []ChangelogEntry) {
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && less(entries[j].Date, entries[j-1].Date); j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}
}

// less compares ISO dates; equal dates keep a stable order.
func less(a, b string) bool {
	if a == b {
		return false
	}
	ta, errA := time.Parse("2006-01-02", a)
	tb, errB := time.Parse("2006-01-02", b)
	if errA != nil || errB != nil {
		return a > b
	}
	return ta.After(tb)
}
