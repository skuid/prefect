package render

import (
	"strings"
	"testing"
)

var (
	configContent = `apiVersion: v1
kind: ConfigMap
metadata:
  name: webapp.config
  namespace: webapp
data:
  s3.bucket: {{.bucket}}
  webapp.domain: {{.webappdomain}}
  webapp.environment: {{.webappenvironment}}
  newrelic.browserId: {{quote .analyticsID}}
  newrelic.app_name: {{quote .newrelicappname}}
  nginx.conf: |-
	server {
		server_name *.{{.webappdomain}};

		listen 80;

		location / {
			error_page 418 = @proxy_to_app;
			return 418;
		}
		location @proxy_to_app {
			proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

			if ($http_x_forwarded_proto != "https") {
				rewrite ^ https://$host$request_uri? permanent;
			}
			proxy_set_header Host      $host;
			proxy_set_header X-Real-IP $remote_addr;
			proxy_pass http://localhost:3000;
		}
		access_log /dev/stdout;
	}
`

	secretContent = `apiVersion: v1
kind: Secret
metadata:
  name: webapp.secret
  namespace: webapp
type: Opaque
data:
  emailkey: {{b64encode .emailKey}}
  pgpassword: {{b64encode .pgpassword}}
  pguser: {{b64encode .pguser}}
  pgdatabase: {{b64encode .pgdatabase}}
  pghost: {{b64encode .pghost}}
  monitoringLicense: {{b64encode .monitoringLicense}}
  encryptedKey: {{b64encode .encryptedKey}}
`

	values = []ContextValue{
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "test", "region": "us-west-2"},
			Key:    "bucket",
			Value:  "webapp-test-environment",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "test"},
			Key:    "webappdomain",
			Value:  "example.com",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "test"},
			Key:    "webappenvironment",
			Value:  "test",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "test"},
			Key:    "analyticsID",
			Value:  "14001122",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "test", "region": "us-west-2"},
			Key:    "newrelicappname",
			Value:  "WebAppTest,WebAppTestUS",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "test"},
			Key:    "emailKey",
			Value:  "some-value",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "test", "region": "us-west-2"},
			Key:    "pgpassword",
			Value:  "supersecret",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "test", "region": "us-west-2"},
			Key:    "pguser",
			Value:  "webapp",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "test", "region": "us-west-2"},
			Key:    "pgdatabase",
			Value:  "webapp",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "test", "region": "us-west-2"},
			Key:    "pghost",
			Value:  "something.something.us-west-2.rds.amazonaws.com",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "test", "region": "us-west-2"},
			Key:    "encryptedKey",
			Value:  "someSUperlong",
		},
		{
			Labels: map[string]string{"type": "kubernetes"},
			Key:    "monitoringLicense",
			Value:  "somelongsecret",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "prod", "region": "us-west-2"},
			Key:    "bucket",
			Value:  "webapp-prod-us-west-2",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "prod"},
			Key:    "webappdomain",
			Value:  "examplesite.com",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "prod"},
			Key:    "webappenvironment",
			Value:  "prod",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "prod"},
			Key:    "analyticsID",
			Value:  "14001123",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "prod", "region": "us-west-2"},
			Key:    "newrelicappname",
			Value:  "WebApp,WebAppUS",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "prod"},
			Key:    "emailKey",
			Value:  "some-sendgrid-value",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "prod", "region": "us-west-2"},
			Key:    "pgpassword",
			Value:  "supersecret",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "prod", "region": "us-west-2"},
			Key:    "pguser",
			Value:  "webapp",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "prod", "region": "us-west-2"},
			Key:    "pgdatabase",
			Value:  "webapp",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "prod", "region": "us-west-2"},
			Key:    "pghost",
			Value:  "something.production.us-west-2.rds.amazonaws.com",
		},
		{
			Labels: map[string]string{"namespace": "webapp", "type": "kubernetes", "env": "prod", "region": "us-west-2"},
			Key:    "encryptedKey",
			Value:  "someSUperlong",
		},
	}
)

func TestRender(t *testing.T) {
	documents := []Document{
		{
			Content: configContent,
		},
		{
			Content: secretContent,
		},
	}

	targets := []Target{
		{
			Name: "Test",
			Selector: map[string]string{
				"env":       "test",
				"type":      "kubernetes",
				"namespace": "webapp",
				"region":    "us-west-2",
			},
		},
		{
			Name: "Prod",
			Selector: map[string]string{
				"env":       "prod",
				"type":      "kubernetes",
				"namespace": "webapp",
				"region":    "us-west-2",
			},
		},
	}

	for _, doc := range documents {
		for _, target := range targets {
			_, err := doc.Render(target, values)
			if err != nil {
				t.Errorf("Unexpected Error: %s", err.Error())
			}
		}
	}
}

func TestGetTemplateContext(t *testing.T) {
	cvs := ContextValues{
		{
			Labels: map[string]string{
				"type":   "kubernetes",
				"env":    "prod",
				"region": "us-west-2",
			},
			Key:   "appLabel",
			Value: "prodApp",
		},
		{
			Labels: map[string]string{
				"namespace": "default",
				"type":      "kubernetes",
				"region":    "us-west-2",
			},
			Key:   "appLabel",
			Value: "k8sDefault",
		},
	}

	selector := map[string]string{
		"type":      "kubernetes",
		"region":    "us-west-2",
		"namespace": "default",
		"env":       "prod",
	}

	response := cvs.GetTemplateContext(selector)

	want := map[string]string{"appLabel": "k8sDefault"}

	for k, v := range want {
		value, ok := response[k]
		if !(ok && strings.Compare(value, v) == 0) {
			t.Errorf("Expected %s to be %s, not %s", k, v, value)
		}
	}
}

func TestRenderInvalidTemplate(t *testing.T) {
	doc := Document{
		Content: "Invalid template {{ }",
	}

	_, err := doc.Render(Target{}, ContextValues{})

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestRenderError(t *testing.T) {
	doc := Document{
		Content: "Valid template {{.missingValue }}",
	}

	_, err := doc.Render(Target{}, ContextValues{})

	if err == nil {
		t.Error("Expected error, got nil")
	}
}
