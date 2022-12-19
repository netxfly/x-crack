package plugins_test

import (
	"x-crack/models"
	"x-crack/plugins"
	"x-crack/vars"

	"testing"
)

func TestScanRdp(t *testing.T) {
	s := models.Service{Ip: "127.0.0.1", Port: 3389, Username: "administrator", Password: vars.PASS}
	t.Log(plugins.ScanRDP(s))
}