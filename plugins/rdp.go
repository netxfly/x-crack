package plugins

import (
	"fmt"
	"github.com/hi-unc1e/grdp"
	"x-crack/models"
)

func ScanRDP(service models.Service) (err error, result models.ScanResult) {
	result.Service = service
	IpPort := fmt.Sprintf("%v:%d", service.Ip, service.Port)
	err = grdp.LoginForRDP(IpPort, "DESKTOP-Q1Test", service.Username, service.Password)
	if err != nil {
		fmt.Println(fmt.Errorf("[x224 connect err] %v", err))
	} else {
		result.Result = true
	}
	return err, result
}
