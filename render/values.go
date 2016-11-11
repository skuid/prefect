package render

import (
	"fmt"
	"strings"
)

// ContextValues is a slice type for ContextValue
type ContextValues []ContextValue

// GetTemplateContext returns a map to be used for a template context
//
// If there are multiple ContextValues that have the same Key, the one
// with the most specific label is chosen. That is, given a selector such as:
//
// 		map[string]string{"region": "us-west-2", "env": "test", "type": "kubernetes"}
//
// and Values contained the following:
//
// 		ContextValues{
// 			{
// 				Labels: `{"env": "test"}`,
// 				Key: "log-level",
// 				Value: "info",
// 			},
// 			{
// 				Labels: `{"env": "test", "region": "us-west-2"}`,
// 				Key: "log-level",
// 				Value: "debug",
// 			},
// 		}
//
// the resulting map would contain the value:
//
// 		map[string]string{"log-level": "debug"}
//
// While both ContextValues matched the selector, the one with the Value of
// "debug" had more matches to the selector, so it was chosen.
//
// In the case of a tie, the last ContextValue evaluated wins.
func (cv ContextValues) GetTemplateContext(selector map[string]string) map[string]string {
	// Get a list of redundant values
	valuesByKey := map[string]ContextValues{}
	for _, value := range cv {
		if matches, ok := valuesByKey[value.Key]; ok {
			matches = append(matches, value)
			valuesByKey[value.Key] = matches
		} else {
			valuesByKey[value.Key] = ContextValues{value}
		}
	}

	// Get the "best" value for each key
	response := map[string]string{}
	for key, matches := range valuesByKey {
		response[key] = getBestMatch(matches, selector).Value
	}

	return response
}

// Get the Value that has the most matches
func getBestMatch(cv ContextValues, selector map[string]string) ContextValue {
	var max float64
	var maxPosition int

	for i := range cv {
		score := cv[i].scoreLabel(selector)
		if score >= max {
			if max > 0 {
				fmt.Println("Warning! Two context values have same score, choosing the last one!")
				fmt.Printf("    %s (%f) matched      %s\n", cv[maxPosition], max, selector)
				fmt.Printf("    %s (%f) also matches %s\n", cv[i], score, selector)
			}
			max = score
			maxPosition = i
		}
	}
	return cv[maxPosition]
}

// ContextValue represents a value to be injected into a template
type ContextValue struct {
	Labels map[string]string `json:"labels" yaml:"labels"`
	Key    string            `json:"key" yaml:"key"`
	Value  string            `json:"value" yaml:"value"`
}

func (cv ContextValue) scoreLabel(selector map[string]string) float64 {
	score := float64(0)
	for k, v := range cv.Labels {
		selectorValue, ok := selector[k]
		if ok && strings.Compare(v, selectorValue) == 0 {
			score += float64(1)
		}
	}
	return score / float64(len(selector))
}
