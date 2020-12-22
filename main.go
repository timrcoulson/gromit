package main

import (
	"crypto/subtle"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/timrcoulson/gromit/calendar"
	"github.com/timrcoulson/gromit/gmail"
	"github.com/timrcoulson/gromit/guitar"
	"github.com/timrcoulson/gromit/news"
	"github.com/timrcoulson/gromit/trello"
	"io/ioutil"
	"log"
	"net/http"
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
	modules = append(modules, &guitar.Guitar{})

	daily := func() string {
		output := fmt.Sprintf("=== Good Morning, Tim. Today is %v ===\n\n", time.Now().Format("Monday 2 Jan 2006"))
		for _, module := range modules {
			output += module.Output()
		}

		return output
	}

	c := cron.New()
	c.AddFunc("00 06 * * *", func() {
		print(daily())
	})
	c.Start()

	http.HandleFunc("/", BasicAuth(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(daily()))
		writer.WriteHeader(200)
	}, os.Getenv("USERNAME"), os.Getenv("PASSWORD"), "Please enter your username and password for this site"))

	log.Fatal(http.ListenAndServe(":80", nil))
}

func BasicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}
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
