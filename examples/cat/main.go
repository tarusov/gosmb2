package main

import (
	"fmt"
	"log"
	"os"

	smb "github.com/tarusov/gosmb2"
	"github.com/tarusov/gosmb2/model"
)

func main() {
	/*
		auth := &smb.Auth{
			Type:     smb.AuthTypeNTLM,
			Domain:   "wincp-32",
			Username: "pangeo",
			Password: "Pangeoacess*",
		}
	*/
	auth := &model.Auth{
		Type: model.AuthTypeNegotiate,
		//Domain:   "WORKGROUP",
		//Username: "zaqxsw",
		//Password: "zaqxsw1",
	}

	share, err := smb.Dial("//127.0.0.1/public", auth)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer share.Close()

	for i := 0; i < 3; i++ {
		err := share.Echo()
		if err != nil {
			fmt.Println("echo error", err.Error())
		}
	}

	f, err := share.OpenFile("hello.txt", os.O_RDONLY)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	i, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Seek(5, 0)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1024)
	_, err = f.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(buf))

	fmt.Println(i.ModTime())

	d, err := share.OpenDir("")
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	files, err := d.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name, f.Type, f.Size)
	}
}
