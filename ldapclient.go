package ldapclient

import (
	"crypto/tls"
	"fmt"
	"github.com/go-ldap/ldap/v3"
)

type LDAPClient struct {
	BindDN             string
	BindPassword       string
	Host               string
	Port               int
	ServerName         string
	Conn               *ldap.Conn
	InsecureSkipVerify bool
	UseTLS             bool
	StartTLS           bool
	ClientCertificates []tls.Certificate // Adding client certificates
}

func NewLDAPClient(host string, port int, dn, password string) LDAPClient {
	lc := LDAPClient{
		BindDN:       dn,
		BindPassword: password,
		Host:         host,
		Port:         port,
	}
	switch port {
	case 389:
		lc.StartTLS = true
		lc.UseTLS = false
	default:
		lc.StartTLS = false
		lc.UseTLS = true
	}
	return lc
}

// Connect connects to the ldap backend.
func (lc *LDAPClient) Connect() (err error) {
	if lc.Conn != nil {
		return nil
	}
	address := fmt.Sprintf("%s:%d", lc.Host, lc.Port)
	// No SSL, plaintext or starttls
	if !lc.UseTLS || lc.StartTLS {
		lc.Conn, err = ldap.Dial("tcp", address)
		if err != nil {
			return
		}
		// Reconnect with TLS
		if lc.StartTLS {
			err = lc.Conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
		}
		return
	}
	// SSL
	if lc.ServerName == "" {
		lc.ServerName = lc.Host
	}
	config := &tls.Config{
		InsecureSkipVerify: lc.InsecureSkipVerify,
		ServerName:         lc.ServerName,
	}
	if len(lc.ClientCertificates) > 0 {
		config.Certificates = lc.ClientCertificates
	}
	lc.Conn, err = ldap.DialTLS("tcp", address, config)
	return
}

// Close closes the ldap backend connection.
func (lc *LDAPClient) Close() {
	if lc.Conn != nil {
		lc.Conn.Close()
		lc.Conn = nil
	}
}

// Authenticate authenticates the user against the ldap backend.
func (lc *LDAPClient) Authenticate() (err error) {
	err = lc.Connect()
	if err != nil {
		return
	}

	err = lc.Conn.Bind(lc.BindDN, lc.BindPassword)
	return
}
