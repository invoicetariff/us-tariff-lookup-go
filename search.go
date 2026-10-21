package tariff

import (
	"strings"
)

// Hit is one search result: the matched MFN entry plus which surface matched.
type Hit struct {
	Entry     MFNEntry `json:"entry"`
	MatchedBy string   `json:"matchedBy"` // "code" | "en" | "zh"
}

// SearchHTS finds codes by HTS prefix/substring, English description or
// Chinese keyword. Multi-word queries use AND semantics — every whitespace-
// separated token must match somewhere on the entry (description, code or
// 中文关键词). Results keep the seed's curation order; limit <= 0 defaults to 10.
func SearchHTS(query string, limit int) []Hit {
	q := strings.TrimSpace(strings.ToLower(query))
	if q == "" {
		return nil
	}
	if limit <= 0 {
		limit = 10
	}
	tokens := strings.Fields(q)

	var hits []Hit
	for _, m := range seed.MFN {
		matched := classify(m, tokens)
		if matched == "" {
			continue
		}
		hits = append(hits, Hit{Entry: m, MatchedBy: matched})
		if len(hits) == limit {
			break
		}
	}
	return hits
}

// classify reports the strongest surface all tokens matched, or "" for no
// match. Priority: code > en > zh.
func classify(m MFNEntry, tokens []string) string {
	code := m.HTS8
	desc := strings.ToLower(m.Description)
	zh := strings.ToLower(m.DescriptionZh + joinKeywords(m.KeywordsZh))

	codeOK, enOK, zhOK := true, true, true
	for _, t := range tokens {
		inCode := strings.Contains(code, t)
		inEn := strings.Contains(desc, t)
		inZh := strings.Contains(zh, t)
		if !inCode {
			codeOK = false
		}
		if !inEn {
			enOK = false
		}
		if !inZh {
			zhOK = false
		}
		if !codeOK && !enOK && !zhOK {
			return ""
		}
	}
	switch {
	case codeOK:
		return "code"
	case enOK:
		return "en"
	case zhOK:
		return "zh"
	default:
		return ""
	}
}

func joinKeywords(kw []string) string {
	if len(kw) == 0 {
		return ""
	}
	return strings.Join(kw, "")
}
