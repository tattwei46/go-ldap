package main

import (
	"fmt"

	"github.com/go-ldap/ldap"
)

var ldapURL = "ldap://localhost:10389"
var ldapConn *ldap.Conn

func main() {
	var err error

	// CONNECT LDAP
	ldapConn, err = connectLDAP(ldapURL)
	if err != nil {
		panic(err.Error())
	}

	if ldapConn != nil {
		defer ldapConn.Close()
	}
	fmt.Println("Connection success!")

	// AUTHENTICATE ADMIN
	if err := authenticateAdmin(ldapConn, "uid=admin,ou=system", "secret"); err != nil {
		panic(err.Error())
	}
	fmt.Println("Authentication success!")

	// SEARCH USER
	//if found, err := searchUser(ldapConn, "ou=system", "fooUser"); err != nil {
	//	panic(err.Error())
	//} else {
	//	fmt.Println(found)
	//}
	//fmt.Println("Search success!")

	// ADD NEW USER
	//if err := addUserAndSetPwd(ldapConn, "cn=fooUser,ou=users,ou=system", "fooUser", "1111"); err != nil {
	//	panic(err.Error())
	//}
	//fmt.Println("New User Added!")

	// CHANGE PASSWORD
	//if err := changePassword(ldapConn, "cn=fooUser,ou=users, ou=system", "1111", "22222"); err != nil {
	//	panic(err.Error())
	//}
	//fmt.Println("Password change success!")

	// DELETE REQUEST
	//if err := removeUser(ldapConn, "cn=fooUser,ou=Users,ou=system"); err != nil {
	//	panic(err.Error())
	//}
	//fmt.Printf("User successfully removed")

}

func removeUser(conn *ldap.Conn, dn string) error {
	delReq := ldap.NewDelRequest(dn, []ldap.Control{})

	return conn.Del(delReq)
}

func changePassword(conn *ldap.Conn, dn, oldPwd, newPwd string) error {
	var err error
	if err := conn.Bind(dn, oldPwd); err != nil {
		return err
	}

	passwordModifyRequest := ldap.NewPasswordModifyRequest("", oldPwd, newPwd)
	_, err = conn.PasswordModify(passwordModifyRequest)
	return err
}

func addUserAndSetPwd(conn *ldap.Conn, dn, name, password string) error {
	addReq := ldap.NewAddRequest(dn, []ldap.Control{})
	addReq.Attribute("objectClass", []string{"top", "organizationalPerson", "inetOrgPerson", "person"})
	addReq.Attribute("cn", []string{name})
	addReq.Attribute("sn", []string{name})

	if err := conn.Add(addReq); err != nil {
		return err
	}
	fmt.Printf("User %s Successfully Added!\n", name)

	modReq := ldap.NewModifyRequest(dn, []ldap.Control{})
	modReq.Replace("userPassword", []string{password})

	if err := conn.Modify(modReq); err != nil {
		return err
	}
	fmt.Println("Password Successfully Set!")
	return nil
}

func searchUser(conn *ldap.Conn, baseDN, name string) (string, error) {
	var found string
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(cn=%s))", name),
		[]string{"dn"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return found, err
	}

	if len(sr.Entries) != 1 {
		return found, fmt.Errorf("user does not exist or too many entries returned")
	}

	return sr.Entries[0].DN, nil //cn=fooUser,ou=users,ou=system
}

func connectLDAP(url string) (*ldap.Conn, error) {
	l, err := ldap.DialURL(url)
	return l, err
}

func authenticateAdmin(ldapConn *ldap.Conn, dn, pwd string) error {
	return ldapConn.Bind(dn, pwd)
}
