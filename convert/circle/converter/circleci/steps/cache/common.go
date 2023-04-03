package cache

import (
	"github.com/drone/go-convert/convert/circle/commons"
)

func getKey(key *string, keys []string) string {
	if key != nil && *key != "" {
		return *key
	}

	if len(keys) > 0 {
		return keys[0]
	}
	return "replace-cache-key"
}

func getBucket(opts commons.Opts) string {
	bucket := ""
	if opts.GCS != nil {
		bucket = opts.GCS.Bucket
	} else if opts.S3 != nil {
		bucket = opts.S3.Bucket
	}
	if bucket == "" {
		bucket = "replace-bucket"
	}
	return bucket
}

func getRegion(opts commons.Opts) string {
	backend := getBackend(opts)

	region := ""
	if backend == "s3" {
		region = opts.S3.Region
	}
	if region == "" {
		region = "replace-region"
	}
	return region
}

func getBackend(opts commons.Opts) string {
	if opts.GCS != nil {
		return "gcs"
	} else if opts.S3 != nil {
		return "s3"
	}
	return "gcs"
}

func getGCSJSONKey(opts commons.Opts) string {
	if opts.GCS != nil {
		return string(opts.GCS.JSONKey)
	}
	return ""
}

func getS3AccessKey(opts commons.Opts) string {
	if opts.S3 != nil {
		return string(opts.S3.AccessKey)
	}
	return ""
}

func getS3SecretKey(opts commons.Opts) string {
	if opts.S3 != nil {
		return string(opts.S3.SecretKey)
	}
	return ""
}
