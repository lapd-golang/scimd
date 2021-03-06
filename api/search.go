package api

import (
	"encoding/json"
	"strings"

	defaults "github.com/mcuadros/go-defaults"
)

// Attributes represents ...
type Attributes struct {
	Attributes         []string `form:"attributes" json:"attributes,omitempty" validate:"dive,attrpath"`
	ExcludedAttributes []string `form:"excludedAttributes" json:"excludedAttributes,omitempty" validate:"dive,attrpath"`
}

// Explode splits the attributes content by comma.
//
// Notice that this assumes that URNs containing comma/s CANNOT be used.
func (a *Attributes) Explode() {
	attributesAcc := []string{}
	for _, x := range a.Attributes {
		attributesAcc = append(attributesAcc, strings.Split(x, ",")...)

	}

	excludedAttributesAcc := []string{}
	for _, y := range a.ExcludedAttributes {
		excludedAttributesAcc = append(excludedAttributesAcc, strings.Split(y, ",")...)
	}

	a.Attributes = attributesAcc
	a.ExcludedAttributes = excludedAttributesAcc
}

type Filter string

const (
	AscendingOrder  = "ascending"
	DescendingOrder = "descending"
)

type Sorting struct {
	SortBy    string `form:"sortBy" json:"sortBy,omitempty"`
	SortOrder string `form:"sortOrder" json:"sortOrder,omitempty" default:"ascending" validate:"omitempty,eq=ascending|eq=descending"`
}

type Pagination struct {
	StartIndex int `form:"startIndex" json:"startIndex,omitempty" default:"1" validate:"omitempty,gt=0"`
	Count      int `form:"count" json:"count,omitempty" mold:"min=0"`
}

// Search represents the set of parameters of a search query
type Search struct {
	Attributes
	Filter `form:"filter" json:"filter,omitempty"` // (todo)> add validator
	Sorting
	Pagination
}

// NewSearch instantiates a Search instance with defaults
func NewSearch() *Search {
	s := &Search{}
	defaults.SetDefaults(s)
	return s
}

// UnmarshalJSON unmarshals an Attribute taking into account defaults
func (s *Search) UnmarshalJSON(data []byte) error {
	defaults.SetDefaults(s)

	type aliasType Search
	alias := aliasType(*s)
	err := json.Unmarshal(data, &alias)

	*s = Search(alias)
	return err
}
