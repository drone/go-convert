# v0 to v1 Pipeline Converter

## Single Pipeline Conversion
```bash 
go run ./convert/v0tov1/cmd/v0tov1/main.go --file_path path/to/pipeline.yaml
```
Outputs to: `path/to/pipeline_v1.yaml`

## Batch Conversion
Converts all pipelines from `base_dir/v0/` to `base_dir/v1/`

```bash 
go run ./convert/v0tov1/cmd/v0tov1/main.go --base_dir path/to/base_dir
```

**Example:**
```bash
go run ./convert/v0tov1/cmd/v0tov1/main.go --base_dir convert/v0tov1/test_pipelines/
```
Converts all pipelines from `convert/v0tov1/test_pipelines/v0/` → `convert/v0tov1/test_pipelines/v1/`