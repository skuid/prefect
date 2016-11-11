package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestSelectorSetSet(t *testing.T) {
	cases := []struct {
		flags    []string
		selector selectorSet
		want     string
	}{
		{
			[]string{"env=test", "type=k8s"},
			selectorSet{},
			"env=test,type=k8s",
		},
		{
			[]string{"env=test", "type=k8s", "type=kubernetes"},
			selectorSet{},
			"env=test,type=k8s,type=kubernetes",
		},
	}

	for _, c := range cases {
		for _, f := range c.flags {
			c.selector.Set(f)
		}
		if strings.Compare(c.selector.String(), c.want) != 0 {
			t.Errorf("Expected %s, got %s", c.want, c.selector.String())
		}
	}
}

func TestSelectorSetType(t *testing.T) {
	want := "string"
	selector := selectorSet{}
	if strings.Compare(selector.Type(), want) != 0 {
		t.Errorf("Expected %s, got %s", want, selector.Type())
	}
}

func TestSelectorSetToMap(t *testing.T) {
	cases := []struct {
		flags    []string
		selector selectorSet
		want     map[string]string
	}{
		{
			[]string{"env=test", "type=k8s"},
			selectorSet{},
			map[string]string{"env": "test", "type": "k8s"},
		},
		{
			[]string{"env=test", "type=k8s", "type=kubernetes"},
			selectorSet{},
			map[string]string{"env": "test", "type": "kubernetes"},
		},
	}

	for _, c := range cases {
		for _, f := range c.flags {
			c.selector.Set(f)
		}
		if !reflect.DeepEqual(c.want, c.selector.ToMap()) {
			t.Errorf("Expected %s, got %s", c.want, c.selector.ToMap())
		}
	}
}
