// Copyright 2022 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package circle

// Option configures a Converter option.
type Option func(*Converter)

// WithDockerhub returns an option to set the default
// dockerhub registry connector.
func WithDockerhub(connector string) Option {
	return func(d *Converter) {
		d.dockerhubConn = connector
	}
}

// WithKubernetes returns an option to set the default
// runtime to Kubernetes.
func WithKubernetes(namespace, connector string) Option {
	return func(d *Converter) {
		d.kubeNamespace = namespace
		d.kubeConnector = connector
	}
}

// WithGCS returns an option to the set the default
// storage provider to GCS.
func WithGCS(bucket, token string) Option {
	return func(d *Converter) {
		d.gcsBucket = bucket
		d.gcsToken = token
	}
}

// WithS3 returns an option to the set the default
// storage provider to S3.
func WithS3(access, secret, region, bucket string) Option {
	return func(d *Converter) {
		d.s3AccessKey = access
		d.s3SecretKey = secret
		d.s3Region = region
		d.s3Bucket = bucket
	}
}
