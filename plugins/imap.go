package plugins

import (
	"fmt"
	_ "github.com/emersion/go-imap"
	"log"
	"x-crack/models"

	"github.com/emersion/go-imap/client"
)

func ScanImap(s models.Service) (err error, result models.ScanResult) {
	result.Service = s

	IpPort := fmt.Sprintf("%v:%d", s.Ip, s.Port)
	c, err := client.Dial(IpPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login(s.Username, s.Password); err != nil {
		log.Fatal(err)
	}
	result.Result = true

	log.Println("Done!")
	return err, result
}
