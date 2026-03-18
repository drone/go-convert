package converthelpers

import (
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// StageConvertedData holds the converted v1 outputs for a single stage.
// New fields can be added here as more cross-stage references are needed.
type StageConvertedData struct {
	Service     *v1.ServiceRef
	Environment *v1.EnvironmentRef
	Runtime     *v1.Runtime
}

// StageConversionContext maintains a registry of converted stage data
// keyed by stage identifier. It is populated after each stage conversion
// and queried when a stage uses useFromStage to reference another.
type StageConversionContext struct {
	stages map[string]*StageConvertedData
}

// NewStageConversionContext creates a new empty context.
func NewStageConversionContext() *StageConversionContext {
	return &StageConversionContext{
		stages: make(map[string]*StageConvertedData),
	}
}

// Set stores converted data for a stage identifier.
func (ctx *StageConversionContext) Set(stageID string, data *StageConvertedData) {
	ctx.stages[stageID] = data
}

// Get retrieves converted data for a stage identifier. Returns nil if not found.
func (ctx *StageConversionContext) Get(stageID string) *StageConvertedData {
	return ctx.stages[stageID]
}

// GetService retrieves the converted service from a previously converted stage.
func (ctx *StageConversionContext) GetService(stageID string) *v1.ServiceRef {
	if data := ctx.Get(stageID); data != nil {
		return data.Service
	}
	return nil
}

// GetEnvironment retrieves the converted environment from a previously converted stage.
func (ctx *StageConversionContext) GetEnvironment(stageID string) *v1.EnvironmentRef {
	if data := ctx.Get(stageID); data != nil {
		return data.Environment
	}
	return nil
}

// GetRuntime retrieves the converted runtime from a previously converted stage.
func (ctx *StageConversionContext) GetRuntime(stageID string) *v1.Runtime {
	if data := ctx.Get(stageID); data != nil {
		return data.Runtime
	}
	return nil
}
