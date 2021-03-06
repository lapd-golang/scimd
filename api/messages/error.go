package messages

import (
	"encoding/json"

	"github.com/fabbricadigitale/scimd/api"
	"github.com/fabbricadigitale/scimd/schemas/datatype"
)

//ErrorURI error message urn
const ErrorURI = "urn:ietf:params:scim:api:messages:2.0:Error"

//Error is a struct for wrapping scim error
type Error struct {
	Schemas  []string `json:"schemas"`
	Status   int      `json:"status,required"`
	ScimType string   `json:"scimType,omitempty"`
	Detail   string   `json:"detail,omitempty"`
}

// NewError wraps error in a scim Error struct
func NewError(e error) Error {

	var scimError Error
	// NOTE: Here the int codes substitute http status codes
	// to avoid a weird effect when unmarshal is performed.
	switch e.(type) {
	case *json.SyntaxError:
		scimError.Status = 400
		scimError.ScimType = "invalidSyntax"
		break
	case *datatype.InvalidDataTypeError:
		scimError.Status = 400
		scimError.ScimType = "invalidValue"
		break
	case *datatype.InvalidValueError:
		scimError.Status = 400
		scimError.ScimType = "invalidValue"
		break
	case *json.UnmarshalTypeError:
		scimError.Status = 400
		scimError.ScimType = "invalidValue"
		break
	case *api.InvalidPathError:
		scimError.Status = 400
		scimError.ScimType = "invalidPath"
		break
	case *api.InvalidFilterError:
		scimError.Status = 400
		scimError.ScimType = "invalidFilter"
		break
	case *api.ResourceNotFoundError:
		scimError.Status = 404
		break
	case *api.MissingRequiredPropertyError:
		scimError.Status = 400
		scimError.ScimType = "invalidValue"
		break
	case *api.MutabilityError:
		scimError.Status = 400
		scimError.ScimType = "mutability"
		break
	case *api.UniquenessError:
		scimError.Status = 400
		scimError.ScimType = "uniqueness"
		break
	default:
		scimError.Status = 500
	}

	scimError.Schemas = append(scimError.Schemas, ErrorURI)
	scimError.Detail = e.Error()

	return scimError
}

func (e Error) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(b)
}
