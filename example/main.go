package main

import (
	"log"

	"github.com/tarusov/gosmb2"
)

func main() {
	opts := &gosmb2.SessionOptions{
		Server:       "127.0.0.1",
		Port:         445,
		User:         "admin",
		SecurityMode: "enabled",
		Path:         "/share/readme.txt",
	}

	s, err := gosmb2.NewSession(opts)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Close(); err != nil {
		log.Fatal(err)
	}

	log.Print("everything is ok")
}
