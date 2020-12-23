package main

import (
	"fmt"
	"log"
	"os"

	smb "github.com/tarusov/gosmb2"
)

func main() {
	share, err := smb.Dial("//127.0.0.1:445/sharex", &smb.Auth{
		Type:     smb.AuthTypeNTLM,
		Domain:   "WORKGROUP",
		Username: "zaqxsw",
		Password: "zaqxsw1",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer share.Close()

	f, err := share.OpenFile("rd", os.O_RDONLY)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	i, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(i.ModTime())
}
