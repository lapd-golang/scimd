package validation

import (
	"fmt"
	"reflect"
	"regexp"

	validator "gopkg.in/go-playground/validator.v9"
)

// AttrNameExpr is the source text used to compile a regular expression macching a SCIM attribute name
const AttrNameExpr = `[A-Za-z][\-$_0-9A-Za-z]*`

// AttrNameRegexp is the compiled Regexp built from AttrNameExpr
var AttrNameRegexp = regexp.MustCompile("^" + AttrNameExpr + "$")

var attrName = func(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		str := field.String()
		return AttrNameRegexp.MatchString(str)
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}
