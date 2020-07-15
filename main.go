package main

import (
	"fmt"
	"log"

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
	if err := authenticateUser(ldapConn, "uid=admin,ou=system", "secret"); err != nil {
		panic(err.Error())
	}
	fmt.Println("Authentication success!")

	//if err := authenticateUser(ldapConn, "cn=fooUser,ou=users,ou=system", "1111"); err != nil {
	//	panic(err.Error())
	//}
	//fmt.Println("Authentication success!")

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

	// ADD NEW GROUP
	//if err := addGroup(ldapConn, "ou=super_admin_1,ou=system"); err != nil {
	//	panic(err.Error())
	//}
	//fmt.Println("New group added successfully")

	// ADD USER TO GROUP
	//if err := addUserToGroup(ldapConn, "ou=super_admin,ou=system", "cn=fooUser"); err != nil {
	//	panic(err.Error())
	//}
	//fmt.Println("User added to group!")

	// GET USERS IN GROUP
	if result, err := getGroupUsers(ldapConn, "ou=super_admin,ou=system"); err != nil {
		panic(err.Error())
	} else {
		for _, r := range result {
			fmt.Println(r)
		}
	}
	fmt.Println("Get Users in Group Success!")

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

func addUserToGroup(conn *ldap.Conn, dnGroup, dnUser string) error {
	addReq := ldap.NewAddRequest(fmt.Sprintf("%s,%s", dnUser, dnGroup), []ldap.Control{})
	addReq.Attribute("objectClass", []string{"top", "organizationalPerson", "inetOrgPerson", "person"})
	addReq.Attribute("cn", []string{dnUser})
	addReq.Attribute("sn", []string{dnUser})

	return conn.Add(addReq)
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

func addGroup(conn *ldap.Conn, dn string) error {
	addReq := ldap.NewAddRequest(dn, nil)
	addReq.Attribute("objectClass", []string{"top", "organizationalUnit"})
	return conn.Add(addReq)
}

func getGroupUsers(conn *ldap.Conn, baseDN string) ([]string, error) {
	var result = make([]string, 0)

	searchRequest := ldap.NewSearchRequest(
		baseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=organizationalPerson))", // The filter to apply
		[]string{"dn", "cn"},                    // A list attributes to retrieve
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return result, err
	}

	for _, entry := range sr.Entries {
		result = append(result, fmt.Sprintf("%s: %v", entry.DN, entry.GetAttributeValue("cn")))
	}
	return result, err
}

func connectLDAP(url string) (*ldap.Conn, error) {
	l, err := ldap.DialURL(url)
	return l, err
}

func authenticateUser(ldapConn *ldap.Conn, dn, pwd string) error {
	return ldapConn.Bind(dn, pwd)
}

func addPasswordPolicy(conn *ldap.Conn, dn, pwd string) error {

	controls := []ldap.Control{}
	controls = append(controls, ldap.NewControlBeheraPasswordPolicy())
	bindRequest := ldap.NewSimpleBindRequest(dn, pwd, controls)

	r, err := conn.SimpleBind(bindRequest)
	ppolicyControl := ldap.FindControl(r.Controls, ldap.ControlTypeBeheraPasswordPolicy)

	var ppolicy *ldap.ControlBeheraPasswordPolicy
	if ppolicyControl != nil {
		ppolicy = ppolicyControl.(*ldap.ControlBeheraPasswordPolicy)
	} else {
		log.Printf("ppolicyControl response not available.\n")
	}
	if err != nil {
		errStr := "ERROR: Cannot bind: " + err.Error()
		if ppolicy != nil && ppolicy.Error >= 0 {
			errStr += ":" + ppolicy.ErrorString
		}
		log.Print(errStr)
		return err
	} else {
		logStr := "Login Ok"
		if ppolicy != nil {
			if ppolicy.Expire >= 0 {
				logStr += fmt.Sprintf(". Password expires in %d seconds\n", ppolicy.Expire)
			} else if ppolicy.Grace >= 0 {
				logStr += fmt.Sprintf(". Password expired, %d grace logins remain\n", ppolicy.Grace)
			}
		}
		log.Print(logStr)
	}
	return nil
}
