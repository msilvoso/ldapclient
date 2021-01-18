# ldapclient
Basic ldap client funcitonality

Adapted from https://github.com/jtblin/go-ldap-client

Sample code:
``` go
    client := ldapclient.NewLDAPClient(*host, *port, *bindDN, bindPassword)
    err = client.Connect()
    if err != nil {
    	log.Fatalf("LDAP connect error %s\n", err.Error())
    }
    err = client.Authenticate()
    if err != nil {
    	log.Fatalf("LDAP auth error %s\n", err.Error())
    }
    defer client.Close()
    search := '(cn=*example*)'
    attributes := []string{"sn"} // if empty all attributes are fetched -> !performance 
    searchResult, err := client.Search("o=Things", search, attributes)
    if err != nil {
        log.Fatalf("Search error %s\n", err.Error())
    }
    for _, entry := range searchResult.Entries {
        client.Replace(entry.DN, "sn", []string{"replaced by this"})
    }
    return
```