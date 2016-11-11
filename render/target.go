package render

import (
	"os"
	"strings"
)

// Targets is a slice type of Target
type Targets []Target

// Target defines a filter for filtering the correct ContextValues for a template
type Target struct {
	Name      string            `json:"name" yaml:"name"`
	Directory string            `json:"directory" yaml:"directory"`
	Selector  map[string]string `json:"selector" yaml:"selector"`
}

// MakedirAll creates the directory that the target specifies
func (t Target) MakedirAll() error {
	return os.MkdirAll(t.Directory, os.ModeDir+0755)
}

func jsonContains(superset, subset map[string]string) bool {
	for k, v := range subset {
		subsetValue, ok := superset[k]
		if ok && strings.Compare(subsetValue, v) == 0 {
			continue
		} else {
			return false
		}
	}
	return true
}

// GetMatchingValues filters for values that match the target's selector
func (t *Target) GetMatchingValues(values ContextValues) ContextValues {
	results := ContextValues{}
	for _, contextValue := range values {
		if jsonContains(t.Selector, contextValue.Labels) {
			results = append(results, contextValue)
		}
	}
	return results
}
