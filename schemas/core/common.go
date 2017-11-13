package core

import (
	"time"
)

// A ScimError is a description of a SCIM error.
type ScimError struct {
	Msg string // description of error
}

func (e *ScimError) Error() string { return e.Msg }

// Meta ...
type Meta struct {
	Location     string     `json:"location"`
	ResourceType string     `json:"resourceType"`
	Created      *time.Time `json:"created"`
	LastModified *time.Time `json:"lastModified"`
	Version      string     `json:"version"`
}

// Common ...
type Common struct {
	Schemas []string `json:"schemas"`

	// Common attributes
	ID         string `json:"id"`
	ExternalID string `json:"externaId,omitempty"`
	Meta       Meta   `json:"meta"`
}

func (c *Common) GetCommon() *Common {
	return c
}

func (c *Common) GetResourceType() *ResourceType {
	return GetResourceTypeRepository().Get(c.Meta.ResourceType)
}

func (c *Common) GetSchema() *Schema {
	if rt := c.GetResourceType(); rt != nil {
		return GetSchemaRepository().Get(rt.Schema)
	}
	return nil
}

func (c *Common) GetSchemaExtensions() map[string]*Schema {
	repo := GetSchemaRepository()
	schExts := c.GetResourceType().SchemaExtensions
	schemas := map[string]*Schema{}
	for _, ext := range schExts {
		schemas[ext.Schema] = repo.Get(ext.Schema)
	}
	return schemas
}
