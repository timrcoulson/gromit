package printer

import (
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode"
)

const PrinterName = "default"

func Print(outputs string) {
	// Add the printer
	output, err := exec.Command("lpadmin", "-p", PrinterName, "-E", "-v", os.Getenv("PRINTER_TCP")).Output()
	if err != nil {
		log.Println(string(output))
		panic(err)
	}

	clean := strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, outputs)

	ioutil.WriteFile("/etc/gromit/file.txt", []byte(clean), 0644)

	cmd := exec.Command("enscript", "--no-header", "-fCourier7", "/etc/gromit/file.txt","--pages", "1", "--non-printable-format=space", "-d", "default", "-DDuplex:true")

	cmd.Stdin = strings.NewReader(strings.Replace(outputs, "\n", "\r\n", -1))
	o, err := cmd.Output()

	if err != nil {
		panic(err)
	}
	log.Println(string(o))

	time.Sleep(5 * time.Second)
}

func init()  {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
}
