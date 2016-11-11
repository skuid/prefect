# Examples

To render templates using targets, run:

```bash
prefect -t targets.yaml -c kv.yaml config.yaml
prefect -t targets.yaml -c kv.yaml secret.yaml
```

To render templates to STDOUT using selectors, run:

```bash
prefect -s namespace=webapp -s type=kubernetes -s region=us-west-2 -s env=test -c kv.yaml config.yaml
prefect -s namespace=webapp -s type=kubernetes -s region=us-west-2 -s env=prod -c kv.yaml config.yaml
```
