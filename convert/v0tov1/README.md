# v0 to v1 Pipeline Converter

## Single Pipeline Conversion
```bash 
go run ./convert/v0tov1/cmd/v0tov1/main.go --file_path path/to/pipeline.yaml
```
Output: `path/to/pipeline_v1.yaml`

## Batch Conversion
Directory structure required:
- v0 pipelines: `path/to/base_dir/v0/`
- v1 output: `path/to/base_dir/v1/`

```bash 
go run ./convert/v0tov1/cmd/v0tov1/main.go --base_dir path/to/base_dir
```