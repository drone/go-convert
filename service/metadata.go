package service

import (
	"context"

	"gopkg.in/yaml.v3"
)

// entityMetadata holds identifying fields parsed from an entity's YAML for
// logging. Any field may be empty when it is absent or null in the YAML.
type entityMetadata struct {
	Account    string
	Org        string
	Project    string
	EntityType string
	EntityID   string
}

// metaFields captures the identifier fields common to every entity wrapper.
// Pointers/strings default to empty when the key is absent or null.
type metaFields struct {
	Identifier        string `yaml:"identifier"`
	AccountIdentifier string `yaml:"accountIdentifier"`
	OrgIdentifier     string `yaml:"orgIdentifier"`
	ProjectIdentifier string `yaml:"projectIdentifier"`
}

// metaEnvelope mirrors the top-level entity keys so a single lightweight parse
// can extract metadata for any entity type without depending on the full v0
// schema (which would fail on otherwise-recoverable YAML).
type metaEnvelope struct {
	Pipeline *metaFields `yaml:"pipeline"`
	Template *metaFields `yaml:"template"`
	InputSet *metaFields `yaml:"inputSet"`
	Trigger  *metaFields `yaml:"trigger"`
}

// extractEntityMetadata best-effort parses yamlStr and returns the account,
// org, project and entity identifiers for the given entity type. It never
// fails: on a parse error or absent fields it returns whatever could be
// determined (at minimum the entity type).
func extractEntityMetadata(entityType, yamlStr string) entityMetadata {
	meta := entityMetadata{EntityType: entityType}

	var env metaEnvelope
	if err := yaml.Unmarshal([]byte(yamlStr), &env); err != nil {
		return meta
	}

	var f *metaFields
	switch entityType {
	case entityPipeline:
		f = env.Pipeline
	case entityTemplate:
		f = env.Template
	case entityInputSet:
		f = env.InputSet
	case entityTrigger:
		f = env.Trigger
	}

	// Fallback: if the requested type wasn't present, use whichever wrapper
	// the YAML actually contains so metadata is still captured.
	if f == nil {
		switch {
		case env.Pipeline != nil:
			f = env.Pipeline
		case env.Template != nil:
			f = env.Template
		case env.InputSet != nil:
			f = env.InputSet
		case env.Trigger != nil:
			f = env.Trigger
		}
	}

	if f != nil {
		meta.Account = f.AccountIdentifier
		meta.Org = f.OrgIdentifier
		meta.Project = f.ProjectIdentifier
		meta.EntityID = f.Identifier
	}
	return meta
}

// logAttrs returns the metadata as slog key/value pairs. Empty values are
// still emitted so log consumers can rely on a stable set of keys.
func (m entityMetadata) logAttrs() []any {
	return []any{
		"entity_type", m.EntityType,
		"account", m.Account,
		"org", m.Org,
		"project", m.Project,
		"entity_id", m.EntityID,
	}
}

// metadataCtxKey is the context key under which a *entityMetadata pointer is
// stored by the logging middleware/interceptor and populated by the handler.
const metadataCtxKey contextKey = "entityMetadata"

// withMetadataSlot stores an empty metadata pointer in ctx for the handler to
// populate, returning the new context and the pointer.
func withMetadataSlot(ctx context.Context) (context.Context, *entityMetadata) {
	m := &entityMetadata{}
	return context.WithValue(ctx, metadataCtxKey, m), m
}

// setRequestMetadata populates the metadata slot in ctx (if present) with the
// fields extracted from yamlStr for entityType.
func setRequestMetadata(ctx context.Context, entityType, yamlStr string) {
	m, ok := ctx.Value(metadataCtxKey).(*entityMetadata)
	if !ok || m == nil {
		return
	}
	*m = extractEntityMetadata(entityType, yamlStr)
}
