package plugins

import (
	"github.com/knadh/go-pop3"
	"log"
	"x-crack/models"
)

func ScanPop3(service models.Service) (err error, result models.ScanResult) {
	result.Service = service

	// Initialize the client.
	p := pop3.New(pop3.Opt{
		Host:       service.Ip,
		Port:       service.Port,
		TLSEnabled: false,
	})

	// Create a new connection. POP3 connections are stateful and should end
	// with a Quit() once the opreations are done.
	c, err := p.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()

	// Authenticate.
	if err := c.Auth(service.Username, service.Password); err != nil {
		log.Fatal(err)
	}

	result.Result = true
	return err, result
}
