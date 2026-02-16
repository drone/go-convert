package yaml

type VolumeEmptyDir struct {
	MountPath string `json:"mount-path,omitempty"`
	Size  string `json:"size,omitempty"`
	Medium string `json:"medium,omitempty"`
	Target string `json:"target,omitempty"`
}