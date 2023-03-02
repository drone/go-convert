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

package yaml

import "errors"

type Artifacts struct {
	Enabled      bool     `yaml:"enabled,omitempty"`
	Bucket       string   `yaml:"bucket,omitempty"`
	Key          *Secure  `yaml:"key,omitempty"`
	Secret       *Secure  `yaml:"secret,omitempty"`
	Region       string   `yaml:"region,omitempty"`
	Paths        []string `yaml:"paths,omitempty"`
	Branch       string   `yaml:"branch,omitempty"`
	LogFormat    string   `yaml:"log_format,omitempty"`
	TargetPaths  []string `yaml:"target_paths,omitempty"`
	Debug        bool     `yaml:"debug,omitempty"`
	Concurrency  int      `yaml:"concurrency,omitempty"`
	MaxSize      int64    `yaml:"max_size,omitempty"`
	Permissions  string   `yaml:"permissions,omitempty"`
	WorkingDir   string   `yaml:"working_dir,omitempty"`
	CacheControl string   `yaml:"cache_control,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Artifacts) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 bool
	var out2 = struct {
		Enabled      *bool         `yaml:"enabled"`
		Bucket       string        `yaml:"bucket,omitempty"`
		Key          *Secure       `yaml:"key,omitempty"`
		Secret       *Secure       `yaml:"secret,omitempty"`
		Region       string        `yaml:"region,omitempty"`
		Paths        Stringorslice `yaml:"paths,omitempty"`
		Branch       string        `yaml:"branch,omitempty"`
		LogFormat    string        `yaml:"log_format,omitempty"`
		TargetPaths  Stringorslice `yaml:"target_paths,omitempty"`
		Debug        bool          `yaml:"debug,omitempty"`
		Concurrency  int           `yaml:"concurrency,omitempty"`
		MaxSize      int64         `yaml:"max_size,omitempty"`
		Permissions  string        `yaml:"permissions,omitempty"`
		WorkingDir   string        `yaml:"working_dir,omitempty"`
		CacheControl string        `yaml:"cache_control,omitempty"`
		// key aliases
		AWSAccessKeyID *Secure `yaml:"aws_access_key_id,omitempty"`
		AWSAccessKey   *Secure `yaml:"aws_access_key,omitempty"`
		AccessKeyID    *Secure `yaml:"access_key_id,omitempty"`
		AccessKey      *Secure `yaml:"access_key,omitempty"`
		// secret aliases
		AWSSecretAccessKey *Secure `yaml:"aws_secret_access_key,omitempty"`
		AWSSecretKey       *Secure `yaml:"aws_secret_key,omitempty"`
		SecretAccessKey    *Secure `yaml:"secret_access_key,omitempty"`
		SecretKey          *Secure `yaml:"secret_key,omitempty"`
	}{}
	if err := unmarshal(&out1); err == nil {
		v.Enabled = out1
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		v.Enabled = true
		if out2.Enabled != nil {
			v.Enabled = *out2.Enabled
		}
		v.Bucket = out2.Bucket
		v.Key = out2.Key
		v.Secret = out2.Secret
		v.Region = out2.Region
		v.Paths = out2.Paths
		v.Branch = out2.Branch
		v.LogFormat = out2.LogFormat
		v.TargetPaths = out2.TargetPaths
		v.Debug = out2.Debug
		v.Concurrency = out2.Concurrency
		v.MaxSize = out2.MaxSize
		v.Permissions = out2.Permissions
		v.WorkingDir = out2.WorkingDir
		v.CacheControl = out2.CacheControl
		switch {
		case out2.AWSAccessKeyID != nil:
			v.Key = out2.AWSAccessKeyID
		case out2.AWSAccessKey != nil:
			v.Key = out2.AWSAccessKey
		case out2.AccessKeyID != nil:
			v.Key = out2.AccessKeyID
		case out2.AccessKey != nil:
			v.Key = out2.AccessKey
		}
		switch {
		case out2.AWSSecretAccessKey != nil:
			v.Secret = out2.AWSSecretAccessKey
		case out2.AWSSecretKey != nil:
			v.Secret = out2.AWSSecretKey
		case out2.SecretAccessKey != nil:
			v.Secret = out2.SecretAccessKey
		case out2.SecretKey != nil:
			v.Secret = out2.SecretKey
		}
		return nil
	}
	return errors.New("failed to unmarshal artifacts")
}
