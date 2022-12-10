package plugins_test

import (
	"x-crack/models"
	"x-crack/plugins"
	"x-crack/vars"

	"testing"
)

func TestScanImap(t *testing.T) {
	s := models.Service{Ip: "127.0.0.1", Port: 143, Protocol: "imap", Username: vars.USER, Password: vars.PASS}
	t.Log(plugins.ScanImap(s))
}
