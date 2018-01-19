package integration

import (
	"testing"

	"github.com/fabbricadigitale/scimd/schemas/datatype"
	"github.com/stretchr/testify/require"

	"github.com/fabbricadigitale/scimd/api/create"
	"github.com/fabbricadigitale/scimd/schemas/core"
	"github.com/fabbricadigitale/scimd/schemas/resource"
)

func TestCreate(t *testing.T) {

	res := &resource.Resource{
		CommonAttributes: core.CommonAttributes{
			Schemas:    []string{"urn:ietf:params:scim:schemas:core:2.0:User", "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"},
			ExternalID: "5666",
			Meta: core.Meta{
				ResourceType: "User",
				Location:     "something",
			},
		},
	}

	res.SetValues("urn:ietf:params:scim:schemas:core:2.0:User", &datatype.Complex{
		"userName": datatype.String("alelb"),
	})

	res.SetValues("urn:ietf:params:scim:schemas:extension:enterprise:2.0:User", &datatype.Complex{
		"employeeNumber": "701984",
	})

	create.Resource(adapter, resTypeRepo.Get("User"), res)

	retRes, err := adapter.Get(resTypeRepo.Get("User"), res.ID, res.Meta.Version, nil)
	require.Nil(t, err)
	require.NotNil(t, retRes)
	require.Equal(t, res.Meta.Version, retRes.Meta.Version)

}