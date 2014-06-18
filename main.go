package main

import (
	"bytes"
	"github.com/rem7/goradius"
	"log"
)

var ctr = 0

// echo "User-Name=steve,User-Password=testing" | radclient -sx 127.0.0.1:1812 auth secret

func main() {

	log.Printf("Server started")
	server := goradius.RadiusServer{}

	// Add middleware
	server.Use(acceptAll)
	// server.Use(passwordCheck)
	// server.Use(addAttributes)
	server.ListenAndServe("0.0.0.0:1812", "s3cr37")

}

func acceptAll(req, res *goradius.RadiusPacket) error {

	ctr += 1

	if ctr%2000 == 0 {
		log.Printf("%v processed requests", ctr)
	}

	res.Code = goradius.AccessAccept
	return nil

}

func passwordCheck(req, res *goradius.RadiusPacket) error {

	username := req.GetAttribute("User-Name")
	usernameData, _ := username.([]byte)

	password := req.GetAttribute("User-Password")
	passwordData, _ := password.([]byte)

	log.Printf("[radius request] username: [%v] password: [%v]", string(usernameData), string(passwordData))

	if bytes.Equal(passwordData, []byte("testing")) &&
		bytes.Equal(usernameData, []byte("steve")) {
		res.Code = goradius.AccessAccept
	} else {
		res.Code = goradius.AccessReject
	}

	return nil

}

func addAttributes(req, res *goradius.RadiusPacket) error {

	if res.Code == goradius.AccessAccept {
		res.AddAttribute("NAS-Identifier", []byte("rem7"))
		res.AddAttribute("Idle-Timeout", uint32(600))
		res.AddAttribute("Session-Timeout", uint32(10800))
	}

	return nil
}
