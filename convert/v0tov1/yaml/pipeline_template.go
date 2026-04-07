package yaml

type PipelineTemplate struct {
	Uses string                `json:"uses,omitempty" yaml:"uses,omitempty"`
	With *PipelineTemplateWith `json:"with,omitempty" yaml:"with,omitempty"`
}

// PipelineTemplateWith represents the 'with' field for pipeline templates
type PipelineTemplateWith struct {
	Overlay *PipelineTemplateOverlay `json:"overlay,omitempty" yaml:"overlay,omitempty"`
}

// PipelineTemplateOverlay represents the overlay configuration for pipeline templates
type PipelineTemplateOverlay struct {
	Pipeline *Pipeline `json:"pipeline,omitempty" yaml:"pipeline,omitempty"`
}
