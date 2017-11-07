package messages

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/fabbricadigitale/scimd/schemas"
	"github.com/fabbricadigitale/scimd/schemas/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalResource(t *testing.T) {

	repo := schemas.Repository()
	if err := repo.LoadResourceType("../../schemas/core/testdata/user.json"); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := repo.LoadSchema("../../schemas/core/testdata/user_schema.json"); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := repo.LoadSchema("../../schemas/core/testdata/enterprise_user_schema.json"); err != nil {
		t.Log(err)
		t.Fail()
	}

	// Non-normative of SCIM user resource type [https://tools.ietf.org/html/rfc7643#section-8.2]
	dat, err := ioutil.ReadFile("testdata/enterprise_user_resource.json")

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	require.NotNil(t, dat)
	require.Nil(t, err)

	res := Resource{}
	err = json.Unmarshal(dat, &res)

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	equalities := []struct {
		value interface{}
		field interface{}
	}{
		{"2819c223-7f76-453a-919d-413861904646", res.ID},
		{"2819c223-7f76-453a-919d-413861904646", res.Common.ID},
		{[]string{"urn:ietf:params:scim:schemas:core:2.0:User", "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"}, res.Schemas},
		{[]string{"urn:ietf:params:scim:schemas:core:2.0:User", "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"}, res.Common.Schemas},
		{"https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646", res.Meta.Location},
		{"https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646", res.Common.Meta.Location},
	}

	for _, row := range equalities {
		assert.Equal(t, row.value, row.field)
	}

	baseAttr, baseAttrOk := res.Attributes["urn:ietf:params:scim:schemas:core:2.0:User"]
	assert.Equal(t, true, baseAttrOk)

	extAttr, extAttrOk := res.Attributes["urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"]
	assert.Equal(t, true, extAttrOk)

	attrEqualities := []struct {
		value interface{}
		field interface{}
	}{
		{"bjensen@example.com", baseAttr["userName"]},
		{"Babs Jensen", baseAttr["displayName"]},
		{true, baseAttr["active"]},
		{"Ms. Barbara J Jensen, III", baseAttr["name"].(core.Complex)["formatted"]},

		{"701984", extAttr["employeeNumber"]},
		{"4130", extAttr["costCenter"]},
	}

	for _, row := range attrEqualities {
		assert.Equal(t, row.value, row.field)
	}

}