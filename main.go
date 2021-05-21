package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/timrcoulson/gromit/agenda"
	"github.com/timrcoulson/gromit/agenda/calendar"
	"github.com/timrcoulson/gromit/agenda/gmail"
	"github.com/timrcoulson/gromit/printer"
	"github.com/timrcoulson/gromit/services/google"
	"github.com/timrcoulson/gromit/services/spotify"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main()  {
    godotenv.Load()

	rand.Seed(time.Now().UnixNano())

	host := os.Getenv("HOST") + ":" + os.Getenv("PORT")
	fmt.Println("Starting up gromit...")

	// Authorise all the services
	sp := spotify.New(host + "/auth/spotify")
	gg := google.New("https://df387d88010c.ngrok.io" + "/auth/google")

	log.Println(gg.LoginUrl())
	// Init other stuff
	calendar.Init()
	gmail.Init()

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`
		<h1>Gromit is running!</h1>
		<a href="%s">Authorize Google</a>
		<a href="%s">Authorize Spotify</a>
		`, gg.LoginUrl(), sp.LoginUrl())))
	})

	wakeup := func() {
		// Print the agenda every day
		//printer.Print(agenda.Today())

		log.Println("time to wake up")

		// Start the morning playlist
		spotify.Play(os.Getenv("MORNING_PLAYLIST"), 20)
	}

	r.HandleFunc("/sleep", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Time for bed"))


		go func() {
			spotify.Play("spotify:playlist:37i9dQZF1DX9uKNf5jGX6m", 15)

			log.Println("song finished")

			sleepSound := exec.Command("omxplayer", "/home/pi/data/sleep.mp3", "--vol", "-1000")

			err := sleepSound.Start()

			if err != nil {
				log.Println(err)
			}
		}()
	})
	r.HandleFunc("/morning", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Time to wake up!"))

		go wakeup()
	})

	r.HandleFunc("/print", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("printing"))
		printer.Print(agenda.Today())
	})

	r.HandleFunc("/agenda", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(agenda.Today()))
	})

	// Register oauth 2 callbacks
	r.HandleFunc("/auth/spotify", sp.Callback)
	r.HandleFunc("/auth/google", gg.Callback)
	//
	http.Handle("/", r)

	l, _ := time.LoadLocation("Europe/London")
	c := cron.New(cron.WithLocation(l))

	c.AddFunc("00 06 * * *", wakeup)
	c.Start()

	log.Print("Gromit running on " + host)
	log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), nil))
}
