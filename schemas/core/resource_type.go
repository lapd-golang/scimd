package core

// SchemaExtension ...
type SchemaExtension struct {
	Schema   string `json:"schema" validate:"urn,required" mold:"normurn"`
	Required bool   `json:"required"`
}

// ResourceType is a structured resource for "urn:ietf:params:scim:schemas:core:2.0:ResourceType"
type ResourceType struct {
	CommonAttributes
	Name             string            `json:"name" validate:"required"`
	Endpoint         string            `json:"endpoint" validate:"startswith=/,required"`
	Description      string            `json:"description,omitempty"`
	Schema           string            `json:"schema" validate:"urn,required" mold:"normurn"`
	SchemaExtensions []SchemaExtension `json:"schemaExtensions,omitempty" validate:"dive"`
}

// ResourceTypeURI is the Resource Type Configuration schema used by ResourceType
const ResourceTypeURI = "urn:ietf:params:scim:schemas:core:2.0:ResourceType"

// NewResourceType returns a ResourceType filled with min values set which identify a particular schema and resourceType (eg. User)
func NewResourceType(schema, resourceType string) *ResourceType {
	return &ResourceType{
		CommonAttributes: *NewCommon(ResourceTypeURI, "ResourceType", resourceType),
		Schema:           schema,
		Name:             resourceType,
	}
}

var _ ResourceTyper = (*ResourceType)(nil)

// GetIdentifier ...
func (rt ResourceType) GetIdentifier() string {
	return rt.Name
}

// GetSchema returns the resource Schema, if any.
func (rt ResourceType) GetSchema() *Schema {
	return GetSchemaRepository().Pull(rt.Schema)
}

// GetSchemaExtensions returns a map of resource's extensions Schema(s) indexed by URN
func (rt ResourceType) GetSchemaExtensions() map[string]*Schema {
	repo := GetSchemaRepository()
	schExts := rt.SchemaExtensions
	schemas := map[string]*Schema{}
	for _, ext := range schExts {
		schemas[ext.Schema] = repo.Pull(ext.Schema)
	}
	return schemas
}

// GetSchemas returns a map containing all the Schema(s) indexed by URN
func (rt ResourceType) GetSchemas() (map[string]*Schema, error) {
	s := rt.GetSchema()
	if s == nil {
		return nil, ScimError{
			Msg: "ResourceType does not have a schema",
		}
	}
	m := rt.GetSchemaExtensions()
	m[s.GetIdentifier()] = s

	return m, nil
}
