package idolmap

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Subject contains normalized six-axis observations.
type Subject map[string]float64

// Dataset is the canonical JSON shape used by the CLI and TUI.
type Dataset struct {
	Axes     map[string]string `json:"axes"`
	Subjects map[string]Subject `json:"subjects"`
}

// LoadJSON loads a dataset from a JSON file.
func LoadJSON(path string) (Dataset, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Dataset{}, err
	}

	var d Dataset
	if err := json.Unmarshal(b, &d); err != nil {
		return Dataset{}, err
	}
	return d, nil
}

// Search performs a simple local RAG-style retrieval over the dataset.
// It scores subjects by matching axis names, axis labels, and subject names.
func (d Dataset) Search(query string, limit int) []string {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return nil
	}

	type hit struct {
		name  string
		score int
	}
	var hits []hit

	for name, subject := range d.Subjects {
		score := 0
		if strings.Contains(strings.ToLower(name), q) {
			score += 10
		}
		for key, label := range d.Axes {
			if strings.Contains(strings.ToLower(key), q) || strings.Contains(strings.ToLower(label), q) {
				if _, ok := subject[key]; ok {
					score += 3
				}
			}
		}
		if score > 0 {
			hits = append(hits, hit{name: name, score: score})
		}
	}

	for i := 0; i < len(hits); i++ {
		for j := i + 1; j < len(hits); j++ {
			if hits[j].score > hits[i].score {
				hits[i], hits[j] = hits[j], hits[i]
			}
		}
	}

	if limit > 0 && len(hits) > limit {
		hits = hits[:limit]
	}

	out := make([]string, len(hits))
	for i, h := range hits {
		out[i] = h.name
	}
	return out
}

// Explain returns the retrieved subject with its six-axis evidence values.
func (d Dataset) Explain(name string) string {
	s, ok := d.Subjects[name]
	if !ok {
		return fmt.Sprintf("not found: %s", name)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "%s\n", name)
	for key, label := range d.Axes {
		fmt.Fprintf(&b, "  %-12s %.2f  %s\n", key, s[key], label)
	}
	return b.String()
}
