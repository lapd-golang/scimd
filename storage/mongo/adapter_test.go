package mongo

import (
	"fmt"
	"testing"

	"github.com/fabbricadigitale/scimd/api/filter"
)

var filters = []string{
	`emails[type eq "work" and value co "@example.com"]`,
	`userName eq "bjensen" and name.familyName sw "J"`,
	`userType ne "Employee" and not (emails co "example.com" or emails.value co "example.org")`,
	`userName eq "bjensen"`,
	`emails.type ne true`,
	`name.familyName co "O'Malley"`,
	`userName sw "J"`,
	`meta.lastModified gt "2011-05-13T04:42:34Z"`,
	`meta.lastModified ge "2011-05-13T04:42:34Z"`,
	`meta.lastModified lt "2011-05-13T04:42:34Z"`,
	`meta.lastModified le "2011-05-13T04:42:34Z"`,
	`emails[type eq "work" and value co "@example.com"] or ims[type eq "xmpp" and value co "@foo.com"]`,
	`userType eq "Employee" and (emails co "example.com" or emails.value co "example.org")`,
	`userType eq "Employee" and emails[type eq "work" and value co "@example.com"]`,
	`userType eq "Employee" and (emails.type eq "work")`,
	`title pr`,
	`not (userName eq "strings")`,
	`not (userName.Child eq "strings")`,
}

func TestConvertToMongoQuery(t *testing.T) {

	for i, str := range filters {

		ft, err := filter.CompileString(str)
		if err != nil {
			t.Log(err)
		}
		m, err := convertToMongoQuery(ft)
		if err != nil {
			t.Log(err)
		}
		fmt.Printf("%d ----------------\n", i)
		fmt.Printf("Filter %s\n", str)
		fmt.Printf("Converted: %+v\n\n\n", m)
	}

}
