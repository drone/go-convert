package commons

type SecretID string

type Opts struct {
	GCS *GCS `embed:"" prefix:"gcs."`
	S3  *S3  `embed:"" prefix:"s3."`
}

type GCS struct {
	Bucket  string   `name:"bucket" help:"GCS bucket."`
	JSONKey SecretID `name:"json-key" help:"Harness secret identifier for GCP JSON key secret id."`
}

type S3 struct {
	Bucket    string   `name:"bucket" help:"S3 bucket."`
	Region    string   `name:"region" help:"S3 region."`
	AccessKey SecretID `name:"access-key" help:"Harness secret identifier for S3 access key."`
	SecretKey SecretID `name:"secret-key" help:"Harness secret identifier for S3 secret key."`
}
