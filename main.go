package main

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap"
)

func main() {
	ldapURL := "ldap://localhost:10389"
	l, err := ldap.DialURL(ldapURL)

	if err != nil {
		fmt.Println(err.Error())
	}

	if l != nil {
		defer l.Close()
	}

	fmt.Println("Connection success!")

	err = l.Bind("uid=admin,ou=system", "secret")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Authentication success!")

	// SEARCH USER
	//searchRequest := ldap.NewSearchRequest(
	//	"ou=system",
	//	ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
	//	fmt.Sprintf("(&(objectClass=organizationalPerson)(cn=%s))", "fooUser"),
	//	[]string{"dn"},
	//	nil,
	//)
	//
	//sr, err := l.Search(searchRequest)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//if len(sr.Entries) != 1 {
	//	log.Fatal("User does not exist or too many entries returned")
	//}
	//
	//userdn := sr.Entries[0].DN
	//fmt.Println(userdn) //cn=fooUser,ou=users,ou=system

	// ADD NEW USER
	//addReq := ldap.NewAddRequest("cn=fooUser,ou=users,ou=system", []ldap.Control{})
	//addReq.Attribute("objectClass", []string{"top", "organizationalPerson", "inetOrgPerson", "person"})
	//addReq.Attribute("cn", []string{"fooUser"})
	//addReq.Attribute("sn", []string{"fooUser"})
	//
	//if err := l.Add(addReq); err != nil {
	//	panic(fmt.Sprintf("error adding service, %s", err))
	//}
	//
	//fmt.Println("Record Added!")

	// CHANGE PASSWORD
	//err = l.Bind("cn=fooUser,ou=users, ou=system", "1111")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//passwordModifyRequest := ldap.NewPasswordModifyRequest("", "1111", "22222")
	//_, err = l.PasswordModify(passwordModifyRequest)
	//
	//if err != nil {
	//	log.Fatalf("Password could not be changed: %s", err.Error())
	//}

	// DELETE REQUEST
	delReq := ldap.NewDelRequest("cn=fooUser,ou=Users,ou=system", []ldap.Control{})

	if err := l.Del(delReq); err != nil {
		log.Fatalf("Error deleting service: %v", err)
	}

	fmt.Println("Deleted success!")

}
