package plugins_test

import (
	"testing"
	"x-crack/models"
	"x-crack/plugins"
	"x-crack/vars"
)

func TestScanLdap(t *testing.T) {
	s := models.Service{Ip: "127.0.0.1", Port: 389, Protocol: "ldap", Username: vars.USER, Password: vars.PASS}
	t.Log(plugins.ScanLdap(s))
}

func test_query() {
}
