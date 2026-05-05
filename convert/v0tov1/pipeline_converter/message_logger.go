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

package pipelineconverter

import "github.com/drone/go-convert/convert/v0tov1/messagelog"

// The structured converter-message logger lives in the leaf package
// convert/v0tov1/messagelog so it can be imported from v0 unmarshallers and
// per-step convert helpers without creating import cycles. The aliases
// below preserve the historical pipelineconverter.X API surface that
// existing tests, the CLI orchestrator, and the service layer rely on.

// Re-exported types.
type (
	Severity       = messagelog.Severity
	Message        = messagelog.Message
	FileMessageLog = messagelog.FileMessageLog
	MessageLogger  = messagelog.MessageLogger
	MessageOption  = messagelog.MessageOption
)

// Re-exported severity constants.
const (
	SeverityInfo    = messagelog.SeverityInfo
	SeverityWarning = messagelog.SeverityWarning
	SeverityError   = messagelog.SeverityError
)

// Re-exported logger and option helpers.
var (
	GetMessageLogger   = messagelog.GetMessageLogger
	ResetMessageLogger = messagelog.ResetMessageLogger
	WithLocation       = messagelog.WithLocation
	WithContext        = messagelog.WithContext
	WithStep           = messagelog.WithStep
	WithStage          = messagelog.WithStage
)
