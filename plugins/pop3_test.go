package plugins_test

import (
	"x-crack/models"
	"x-crack/plugins"
	"x-crack/vars"

	"testing"
)

func TestScanPop3(t *testing.T) {
	s := models.Service{Ip: "127.0.0.1", Port: 110, Protocol: "pop3", Username: vars.USER, Password: vars.PASS}
	t.Log(plugins.ScanPop3(s))
}
