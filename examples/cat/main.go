package main

import (
	"fmt"
	"log"
	"os"

	smb "github.com/tarusov/gosmb2"
	"github.com/tarusov/gosmb2/model"
)

func main() {
	auth := &model.Auth{
		Type: model.AuthTypeNegotiate,
		//Domain:   "WORKGROUP",
		//Username: "zaqxsw",
		//Password: "zaqxsw1",
	}

	share, err := smb.Dial("//192.168.2.8/public", auth)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := share.Close()
		if err != nil {
			log.Printf("close file: %v", err)
		}
	}()

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
	defer func() {
		err := f.Close()
		if err != nil {
			log.Printf("close file: %v", err)
		}
	}()

	// Stat.

	i, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("successful get stat. file created at ", i.ModTime())

	// Seek.

	_, err = f.Seek(8, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("successful seek pos")

	// Read.

	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("successful read. msg: ", string(buf[:n]))

	// List dir.

	d, err := share.OpenDir(".")
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	files, err := d.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Printf("%s\t%s\t%d\n", f.Name, f.Type, f.Size)
	}
}
