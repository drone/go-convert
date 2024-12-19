// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// package yaml provides definitions for the Yaml schema.
package yaml

// Resource enums.
const (
	KindDeployment = "deployment"
	KindPipeline   = "pipeline"
	KindSecret     = "secret"
	KindSignature  = "signature"
)

type (
	// Pipeline defines a pipeline resource.
	Pipeline struct {
		Version     string            `json:"version,omitempty"`
		Kind        string            `json:"kind,omitempty"`
		Type        string            `json:"type,omitempty"`
		Name        string            `json:"name,omitempty"`
		Hmac        string            `json:"hmac,omitempty"`
		Deps        []string          `json:"deps,omitempty" yaml:"depends_on"`
		Node        map[string]string `json:"node,omitempty"`
		Concurrency Concurrency       `json:"concurrency,omitempty"`
		Platform    Platform          `json:"platform,omitempty"`
		Data        Secret            `json:"secret,omitempty"`
		Clone       Clone             `json:"clone,omitempty"`
		Trigger     Conditions        `json:"conditions,omitempty"`
		Environment map[string]string `json:"environment,omitempty"`
		Services    []*Step           `json:"services,omitempty"`
		Steps       []*Step           `json:"steps,omitempty"`
		Volumes     []*Volume         `json:"volumes,omitempty"`
		PullSecrets []string          `json:"image_pull_secrets,omitempty" yaml:"image_pull_secrets"`
		Workspace   Workspace         `json:"workspace,omitempty"`

		// Kubernetes Runner
		DnsConfig      DnsConfig         `json:"dns_config,omitempty"   yaml:"dns_config"`
		HostAliases    []HostAlias       `json:"host_aliases,omitempty" yaml:"host_aliases"`
		Metadata       Metadata          `json:"metadata,omitempty"`
		NodeName       string            `json:"node_name,omitempty"            yaml:"node_name"`
		NodeSelector   map[string]string `json:"node_selector,omitempty"        yaml:"node_selector"`
		ServiceAccount string            `json:"service_account_name,omitempty" yaml:"service_account_name"`
		Tolerations    []Toleration      `json:"tolerations,omitempty"`
		Resource       Resources         `json:"resource,omitempty"`
	}

	// Resources configures resource limits.
	Resources struct {
		Limits   ResourceLimits `json:"limits,omitempty" yaml:"limits"`
		Requests ResourceLimits `json:"requests,omitempty" yaml:"requests"`
	}

	// ResourceLimits configures resource limits.
	ResourceLimits struct {
		CPU    int64     `json:"cpu" yaml:"cpu"`
		Memory BytesSize `json:"memory"`
	}

	// Clone configures the git clone.
	Clone struct {
		Disable    bool `json:"disable,omitempty"`
		Depth      int  `json:"depth,omitempty"`
		Retries    int  `json:"retries,omitempty"`
		SkipVerify bool `json:"skip_verify,omitempty" yaml:"skip_verify"`
		Trace      bool `json:"trace,omitempty"`
	}

	// Concurrency limits pipeline concurrency.
	Concurrency struct {
		Limit int `json:"limit,omitempty"`
	}

	// Platform defines the target platform.
	Platform struct {
		OS      string `json:"os,omitempty"`
		Arch    string `json:"arch,omitempty"`
		Variant string `json:"variant,omitempty"`
		Version string `json:"version,omitempty"`
	}

	// Secret defines an external secret.
	Secret struct {
		Name string `json:"name,omitempty"`
		Path string `json:"path,omitempty"`
	}

	// Step defines a Pipeline step.
	Step struct {
		Command      []string              `json:"command,omitempty"`
		Commands     []string              `json:"commands,omitempty"`
		Detach       bool                  `json:"detach,omitempty"`
		DependsOn    []string              `json:"depends_on,omitempty" yaml:"depends_on"`
		Devices      []*VolumeDevice       `json:"devices,omitempty"`
		DNS          []string              `json:"dns,omitempty"`
		DNSSearch    []string              `json:"dns_search,omitempty" yaml:"dns_search"`
		Entrypoint   []string              `json:"entrypoint,omitempty"`
		Environment  map[string]*Variable  `json:"environment,omitempty"`
		ExtraHosts   []string              `json:"extra_hosts,omitempty" yaml:"extra_hosts"`
		Failure      string                `json:"failure,omitempty"`
		Image        string                `json:"image,omitempty"`
		MemLimit     BytesSize             `json:"mem_limit,omitempty" yaml:"mem_limit"`
		MemSwapLimit BytesSize             `json:"memswap_limit,omitempty" yaml:"memswap_limit"`
		Network      string                `json:"network_mode,omitempty" yaml:"network_mode"`
		Name         string                `json:"name,omitempty"`
		Privileged   bool                  `json:"privileged,omitempty"`
		Pull         string                `json:"pull,omitempty"`
		Resource     Resources             `json:"resource,omitempty"`
		Settings     map[string]*Parameter `json:"settings,omitempty"`
		Shell        string                `json:"shell,omitempty"`
		ShmSize      BytesSize             `json:"shm_size,omitempty" yaml:"shm_size"`
		User         string                `json:"user,omitempty"`
		Volumes      []*VolumeMount        `json:"volumes,omitempty"`
		When         Conditions            `json:"when,omitempty"`
		WorkingDir   string                `json:"working_dir,omitempty" yaml:"working_dir"`
	}

	// Volume that can be mounted by containers.
	Volume struct {
		Name     string          `json:"name,omitempty"`
		EmptyDir *VolumeEmptyDir `json:"temp,omitempty" yaml:"temp"`
		HostPath *VolumeHostPath `json:"host,omitempty" yaml:"host"`
	}

	// VolumeDevice describes a mapping of a raw block
	// device within a container.
	VolumeDevice struct {
		Name       string `json:"name,omitempty"`
		DevicePath string `json:"path,omitempty" yaml:"path"`
	}

	// VolumeMount describes a mounting of a Volume
	// within a container.
	VolumeMount struct {
		Name      string `json:"name,omitempty"`
		MountPath string `json:"path,omitempty" yaml:"path"`
	}

	// VolumeEmptyDir mounts a temporary directory from the
	// host node's filesystem into the container. This can
	// be used as a shared scratch space.
	VolumeEmptyDir struct {
		Medium    string    `json:"medium,omitempty"`
		SizeLimit BytesSize `json:"size_limit,omitempty" yaml:"size_limit"`
	}

	// VolumeHostPath mounts a file or directory from the
	// host node's filesystem into your container.
	VolumeHostPath struct {
		Path string `json:"path,omitempty"`
	}

	// Workspace represents the pipeline workspace configuration.
	Workspace struct {
		Base string `json:"base,omitempty"`
		Path string `json:"path,omitempty"`
	}
)

// Kubernetes specific
type (
	// Metadata defines Kubernetes pod metadata
	Metadata struct {
		Namespace   string            `json:"namespace,omitempty"`
		Annotations map[string]string `json:"annotations,omitempty"`
		Labels      map[string]string `json:"labels,omitempty"`
	}
	// DnsConfig defines Kubernetes pod dnsConfig
	DnsConfig struct {
		Nameservers []string           `json:"nameservers,omitempty"`
		Searches    []string           `json:"searches,omitempty"`
		Options     []DNSConfigOptions `json:"options,omitempty"`
	}
	// DNSConfigOptions dns config option
	DNSConfigOptions struct {
		Name  string  `json:"name,omitempty"`
		Value *string `json:"value,omitempty" yaml:"value"`
	}

	// HostAlias hosts
	HostAlias struct {
		IP        string   `json:"ip,omitempty"`
		Hostnames []string `json:"hostnames,omitempty"`
	}

	// Toleration defines Kubernetes pod tolerations
	Toleration struct {
		Effect            string `json:"effect,omitempty"`
		Key               string `json:"key,omitempty"`
		Operator          string `json:"operator,omitempty"`
		TolerationSeconds *int   `json:"toleration_seconds,omitempty"`
		Value             string `json:"value,omitempty"`
	}
)
