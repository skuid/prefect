# Prefect

A template rendering tool based on labeled context values and selectors

Prefect takes a [go template file](https://golang.org/pkg/text/template/) and applies a set of
context values to it, filtering in the context values based a list of selectors.

## Etymology

**prefect** |ˈprēˌfekt|

noun

1. a chief officer, magistrate, or regional governor in certain countries: the prefect of police.
    * a senior magistrate or governor in the ancient Roman world: _Avitus was prefect of Gaul from AD 439._
1. chiefly Brit. in some schools, a senior student authorized to enforce discipline.

## Example

Given the template `config.yaml`:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: webapp.config
  namespace: default
data:
  s3.bucket: {{.bucket}}
  nginx.domain: {{.domain}}
```

And given a value file `context.yaml` (json or yaml):

```yaml
- labels:
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
  value: skuid.com
```

To render for a test environment, you could run:

```
$ prefect -s env=test -s type=kubernetes -s namespace=default -c context.yaml  config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: webapp.config
  namespace: default
data:
  s3.bucket: s3://test-bucket
  nginx.domain: skuid.com

```

And for a prod environment, you could run:

```
$ prefect -s env=prod -s type=kubernetes -s namespace=default -c context.yaml  config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: webapp.config
  namespace: default
data:
  s3.bucket: s3://prod-bucket
  nginx.domain: skuid.com
```

To automatically render for multiple targets, you can specify a Target file (in JSON or YAML format)

```yaml
- name: K8sTestUsWebapp
  directory: "kubernetes/test/us-east-2/webapp"
  selector:
    type: kubernetes
    env: test
    namespace: webapp
    region: us-east-2
- name: K8sProdUsWebapp
  directory: "kubernetes/prod/us-east-2/webapp"
  selector:
    type: kubernetes
    env: prod
    namespace: webapp
    region: us-east-2
```

and run:

```
$ prefect -t targets.yaml -c kv.yaml config.yaml
Rendered file at kubernetes/test/us-east-2/webapp/config.yaml
Rendered file at kubernetes/prod/us-east-2/webapp/config.yaml
```

Prefect will automatically create the files in the directory specified by each target.

## Functions

There are several custom functions you can call from your template. Currently these are:

- `b64encode` - Base64 encode a string into the template
- `quote` - Surround the rendered variable with quotes

### Example

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: MyConfig
data:
  some-key: {{b64encode .key}}
  numeric-string: {{quote .count}}
```

would render to

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: MyConfig
data:
  some-key: c29tZS1rZXk=
  numeric-string: "3"
```

## Installation

### Prerequisite

You must have go installed

```bash
brew install go
mkdir -p ~/go/{src,bin,pkg}
export GOPATH=~/go
export PATH="$PATH:~/go/bin"

# Append GOPATH to profile
echo 'export GOPATH=~/go' | tee -a ~/.profile
echo 'export PATH="$PATH:~/go/bin"' | tee -a ~/.profile
```

Install prefect

```bash
mkdir -p $GOPATH/src/github.com/skuid/prefect
git clone git@github.com:skuid/prefect.git $GOPATH/src/github.com/skuid/prefect
cd $GOPATH/src/github.com/skuid/prefect
go install .
```

## Usage

```
prefect takes a template and injects context from a given context file.

Usage:
prefect [options] <template>

  <template>
    	The template to read in
  -c, --context string
    	The context values file to use
  -s, --selector value
    	The selectors to use. Each selector should have the format "k=v".
    	Can be specified multiple times, or a comma-separated list
```

## Development

Always, always, always run `go fmt ./...` before committing!

### Running the tests

```bash
go get golang.org/x/tools/cmd/cover

go test -coverprofile=coverage.out
go tool cover -func=coverage.out
go test -coverprofile=coverage.out ./render/
go tool cover -func=coverage.out
```

See the html output of the coverage information

```
go tool cover -html=coverage.out
```

### Updating dependencies

```
go get -u github.com/kardianos/govendor

govendor add +external
```

### Linting

Perfect linting is not required, but it is helpful for new people coming to the code.

```
go get -u github.com/golang/lint/golint

golint ./
golint ./render
```

## License

MIT License (see [LICENSE](/LICENSE))
