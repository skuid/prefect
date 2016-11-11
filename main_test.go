package main

import (
	"github.com/skuid/prefect/render"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	exampleDoc := []byte(`apiVersion: v1
kind: ConfigMap
metadata:
  name: webapp.config
  namespace: default
data:
  s3.bucket: {{.bucket}}
  nginx.domain: {{.domain}}`)
	exampleValue := []byte(`- labels:
    env: test
    namespace: default
    type: kubernetes
  key: bucket
  value: s3://test-bucket
- labels:
    env: prod
    namespace: default
    type: kubernetes
  key: bucket
  value: s3://prod-bucket
- labels:
    namespace: default
    type: kubernetes
  key: domain
  value: skuid.com`)

	want := `apiVersion: v1
kind: ConfigMap
metadata:
  name: webapp.config
  namespace: default
data:
  s3.bucket: s3://test-bucket
  nginx.domain: skuid.com`

	selector := render.Target{
		Selector: selectorSet{"env=test", "namespace=default", "type=kubernetes"}.ToMap(),
	}

	got, err := execute("config", exampleDoc, exampleValue, selector)

	if err != nil {
		t.Error(err.Error())
	}
	if strings.Compare(want, got) != 0 {
		t.Errorf("Didn't get desired output! Expected:\n %s\nGot:\n%s", want, got)
	}
}

func TestExecuteFail(t *testing.T) {
	exampleDoc := []byte(`{.something}`)
	exampleValue := []byte(`{{`)

	_, err := execute("config", exampleDoc, exampleValue, render.Target{})

	if err == nil {
		t.Error("Expected error, got nil")
	}
}
