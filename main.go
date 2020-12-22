package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/timrcoulson/gromit/calendar"
	"github.com/timrcoulson/gromit/gmail"
	"github.com/timrcoulson/gromit/news"
	"github.com/timrcoulson/gromit/trello"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode"
)

const PrinterName = "default"

func init()  {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main()  {
	fmt.Println("Starting up gromit...")

	// Add the printer
	output, err := exec.Command("lpadmin", "-p", PrinterName, "-E", "-v", os.Getenv("PRINTER_TCP")).Output()
	if err != nil {
		log.Println(string(output))
		panic(err)
	}

	// Register modules
	var modules []Module
	modules = append(modules, &calendar.Calendar{})
	modules = append(modules, &gmail.Gmail{})
	modules = append(modules, &trello.Trello{})
	modules = append(modules, &news.News{})

	daily := func() {
		output := fmt.Sprintf("=== Good Morning, Tim. Today is %v ===\n\n", time.Now().Format("Monday 2 Jan 2006"))
		for _, module := range modules {
			output += module.Output()
		}

		fmt.Println(output)
		print(output)
	}

	daily()

	fmt.Println("Gromit shutting down")
}

func print(outputs string)  {
	clean := strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, outputs)
	// Send "Hello, world!" to the printer via a pipe
	ioutil.WriteFile("/tmp/daily.txt", []byte(clean), 0644)
	cmd := exec.Command("enscript", "--no-header", "-fCourier7", "/tmp/daily.txt","--pages", "1", "--non-printable-format=space")

	cmd.Stdin = strings.NewReader(strings.Replace(outputs, "\n", "\r\n", -1))

	cmd.Output()

	time.Sleep(5 * time.Second)
}

type Module interface {
	Output() string
}
