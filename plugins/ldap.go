package plugins

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"log"
	"x-crack/models"
)

const (
	ScopeBaseObject   = 0
	ScopeSingleLevel  = 1
	ScopeWholeSubtree = 2
)

const (
	NeverDerefAliases   = 0
	DerefInSearching    = 1
	DerefFindingBaseObj = 2
	DerefAlways         = 3
)

func ScanLdap(s models.Service) (err error, result models.ScanResult) {
	result.Service = s

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	/*
		LDAP匿名查询
		https://pkg.go.dev/gopkg.in/ldap.v3#pkg-index
	*/
	searchRequest := ldap.NewSearchRequest(
		"DC=server,DC=local", // The base dn to search
		ScopeWholeSubtree, NeverDerefAliases, 0, 0, false,
		"(&(objectClass=organizationalPerson))", // The filter to apply
		[]string{"dn", "cn"},                    // A list attributes to retrieve
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	DNs := ""
	for i, entry := range sr.Entries {
		msg := fmt.Sprintf("DN[%v]: %v\n", i, entry.DN)
		DNs += msg
	}
	log.Printf(DNs)

	// LDAP认证，但fapro并不支持ldap认证，so skip it
	cn := fmt.Sprintf("cn=%v,DC=server,DC=local", s.Username)
	err = l.Bind(cn, s.Password)
	if err != nil {
		log.Fatal(err)
	} else {
		result.Result = true
	}
	return err, result
}
