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

import (
	"encoding/json"
	"github.com/drone/go-convert/internal/flexible"
)

type (
	// Infrastructure provides pipeline infrastructure.
	Infrastructure struct {
		Type string      `json:"type,omitempty"          yaml:"type,omitempty"`
		From string      `json:"useFromStage,omitempty"  yaml:"useFromStage,omitempty"`
		Spec interface{} `json:"spec,omitempty"          yaml:"spec,omitempty"`
	}

	// InfrastructureKubernetesDirectSpec describes Kubernetes Direct infrastructure
	InfrastructureKubernetesDirectSpec struct {
		Conn                         string                              `json:"connectorRef,omitempty"                 yaml:"connectorRef,omitempty"`
		Namespace                    string                              `json:"namespace,omitempty"                    yaml:"namespace,omitempty"`
		Annotations                  *flexible.Field[map[string]string]  `json:"annotations,omitempty"                  yaml:"annotations,omitempty"`
		Labels                       *flexible.Field[map[string]string]  `json:"labels,omitempty"                       yaml:"labels,omitempty"`
		ServiceAccountName           string                              `json:"serviceAccountName,omitempty"           yaml:"serviceAccountName,omitempty"`
		InitTimeout                  string                              `json:"initTimeout,omitempty"                  yaml:"initTimeout,omitempty"`
		NodeSelector                 *flexible.Field[map[string]string]  `json:"nodeSelector,omitempty"                 yaml:"nodeSelector,omitempty"`
		HostNames                    *flexible.Field[[]string]           `json:"hostNames,omitempty"                    yaml:"hostNames,omitempty"`
		Tolerations                  *flexible.Field[[]*Toleration]      `json:"tolerations,omitempty"                  yaml:"tolerations,omitempty"`
		Volumes                      []*Volume                           `json:"volumes,omitempty"                      yaml:"volumes,omitempty"`
		AutomountServiceAccountToken *flexible.Field[bool]               `json:"automountServiceAccountToken,omitempty" yaml:"automountServiceAccountToken,omitempty"`
		ContainerSecurityContext     *flexible.Field[*SecurityContext]   `json:"containerSecurityContext,omitempty"     yaml:"containerSecurityContext,omitempty"`
		PriorityClassName            string                              `json:"priorityClassName,omitempty"            yaml:"priorityClassName,omitempty"`
		OS                           string                              `json:"os,omitempty"                           yaml:"os,omitempty"`
		HarnessImageConnectorRef     string                              `json:"harnessImageConnectorRef,omitempty"     yaml:"harnessImageConnectorRef,omitempty"`
		ImagePullPolicy              string                              `json:"imagePullPolicy,omitempty"              yaml:"imagePullPolicy,omitempty"`
		PodSpecOverlay               string                              `json:"podSpecOverlay,omitempty"               yaml:"podSpecOverlay,omitempty"`
	}

	// InfrastructureVMSpec describes VM infrastructure
	InfrastructureVMSpec struct {
		Type string                `json:"type,omitempty" yaml:"type,omitempty"`
		Spec *InfrastructureVMPool `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	// InfrastructureVMPool describes VM pool configuration
	InfrastructureVMPool struct {
		PoolName string `json:"poolName,omitempty" yaml:"poolName,omitempty"`
		OS       string `json:"os,omitempty"       yaml:"os,omitempty"`
		Identifier string `json:"identifier,omitempty" yaml:"identifier,omitempty"`
	}

	// Toleration defines Kubernetes toleration configuration
	Toleration struct {
		Effect            string               `json:"effect,omitempty"            yaml:"effect,omitempty"`
		Key               string               `json:"key,omitempty"               yaml:"key,omitempty"`
		Operator          string               `json:"operator,omitempty"          yaml:"operator,omitempty"`
		TolerationSeconds *flexible.Field[int] `json:"tolerationSeconds,omitempty" yaml:"tolerationSeconds,omitempty"`
		Value             string               `json:"value,omitempty"             yaml:"value,omitempty"`
	}

	// SecurityContext defines container security context
	SecurityContext struct {
		AllowPrivilegeEscalation *flexible.Field[bool]          `json:"allowPrivilegeEscalation,omitempty" yaml:"allowPrivilegeEscalation,omitempty"`
		ProcMount                string                         `json:"procMount,omitempty"                yaml:"procMount,omitempty"`
		Privileged               *flexible.Field[bool]          `json:"privileged,omitempty"               yaml:"privileged,omitempty"`
		ReadOnlyRootFilesystem   *flexible.Field[bool]          `json:"readOnlyRootFilesystem,omitempty"   yaml:"readOnlyRootFilesystem,omitempty"`
		RunAsNonRoot             *flexible.Field[bool]          `json:"runAsNonRoot,omitempty"             yaml:"runAsNonRoot,omitempty"`
		RunAsGroup               *flexible.Field[int]           `json:"runAsGroup,omitempty"               yaml:"runAsGroup,omitempty"`
		RunAsUser                *flexible.Field[int]           `json:"runAsUser,omitempty"                yaml:"runAsUser,omitempty"`
		Capabilities             *flexible.Field[*Capabilities] `json:"capabilities,omitempty"             yaml:"capabilities,omitempty"`
	}

	// Capabilities defines Linux capabilities
	Capabilities struct {
		Add  *flexible.Field[[]string] `json:"add,omitempty"  yaml:"add,omitempty"`
		Drop *flexible.Field[[]string] `json:"drop,omitempty" yaml:"drop,omitempty"`
	}

	// Volume represents a Kubernetes volume with type-specific spec
	Volume struct {
		Type      string      `json:"type,omitempty"           yaml:"type,omitempty"`
		MountPath string      `json:"mountPath,omitempty"      yaml:"mountPath,omitempty"`
		Spec      interface{} `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	// EmptyDirVolumeSpec defines emptyDir volume spec
	EmptyDirVolumeSpec struct {
		Medium string `json:"medium,omitempty" yaml:"medium,omitempty"`
		Size   string `json:"size,omitempty"   yaml:"size,omitempty"`
	}

	// PersistentVolumeClaimVolumeSpec defines PVC volume spec
	PersistentVolumeClaimVolumeSpec struct {
		ClaimName string                `json:"claimName,omitempty"          yaml:"claimName,omitempty"`
		ReadOnly  *flexible.Field[bool] `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	}

	// HostPathVolumeSpec defines hostPath volume spec
	HostPathVolumeSpec struct {
		Path string `json:"path,omitempty"           yaml:"path,omitempty"`
		Type string `json:"type,omitempty" yaml:"type,omitempty"`
	}

	// ConfigMapVolumeSpec defines configMap volume spec
	ConfigMapVolumeSpec struct {
		Name     string                `json:"name,omitempty"               yaml:"name,omitempty"`
		Optional *flexible.Field[bool] `json:"optional,omitempty" yaml:"optional,omitempty"`
	}

	// SecretVolumeSpec defines secret volume spec
	SecretVolumeSpec struct {
		Name     string                `json:"name,omitempty"               yaml:"name,omitempty"`
		Optional *flexible.Field[bool] `json:"optional,omitempty" yaml:"optional,omitempty"`
	}
)

// UnmarshalJSON implements custom unmarshalling for Infrastructure based on type
func (i *Infrastructure) UnmarshalJSON(data []byte) error {
	type Alias Infrastructure
	aux := &struct {
		Type string          `json:"type"`
		Spec json.RawMessage `json:"spec"`
		*Alias
	}{
		Alias: (*Alias)(i),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	i.Type = aux.Type
	i.From = aux.From

	if len(aux.Spec) == 0 {
		i.Spec = nil
		return nil
	}

	// Unmarshal spec based on type
	switch aux.Type {
	case "KubernetesDirect":
		var spec InfrastructureKubernetesDirectSpec
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return err
		}
		i.Spec = &spec
	case "VM":
		var spec InfrastructureVMSpec
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return err
		}
		i.Spec = &spec
	default:
		// For unknown types, keep as interface{}
		var spec interface{}
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return err
		}
		i.Spec = spec
	}

	return nil
}

func (v *Volume) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type      string          `json:"type,omitempty"`
		MountPath string          `json:"mountPath,omitempty"`
		Spec      json.RawMessage `json:"spec,omitempty"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	v.Type = aux.Type
	v.MountPath = aux.MountPath

	if len(aux.Spec) == 0 {
		v.Spec = nil
		return nil
	}

	switch v.Type {
	case "EmptyDir":
		var spec EmptyDirVolumeSpec
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return err
		}
		v.Spec = spec
	case "PersistentVolumeClaim":
		var spec PersistentVolumeClaimVolumeSpec
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return err
		}
		v.Spec = spec
	case "HostPath":
		var spec HostPathVolumeSpec
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return err
		}
		v.Spec = spec
	case "ConfigMap":
		var spec ConfigMapVolumeSpec
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return err
		}
		v.Spec = spec
	case "Secret":
		var spec SecretVolumeSpec
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return err
		}
		v.Spec = spec
	default:
		v.Spec = nil
	}

	return nil
}