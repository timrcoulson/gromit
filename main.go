package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/robfig/cron/v3"
	"github.com/timrcoulson/gromit/agenda"
	"github.com/timrcoulson/gromit/agenda/calendar"
	"github.com/timrcoulson/gromit/agenda/gmail"
	"github.com/timrcoulson/gromit/printer"
	"github.com/timrcoulson/gromit/services/google"
	"github.com/timrcoulson/gromit/services/spotify"
	"log"
	"net/http"
	"os"
)

func main()  {
	host := os.Getenv("HOST") + ":" + os.Getenv("PORT")
	fmt.Println("Starting up gromit...")

	// Authorise all the services
	sp := spotify.New(host + "/auth/spotify")
	gg := google.New(host + "/auth/google")

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

	agenda.Today()

	r.HandleFunc("/spotify", func(writer http.ResponseWriter, request *http.Request) {
		spotify.Play(os.Getenv("MORNING_PLAYLIST"))
		writer.Write([]byte("playing morning"))
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

	http.Handle("/", r)

	c := cron.New()
	c.AddFunc("00 06 * * *", func() {
		// Print the agenda every day
		printer.Print(agenda.Today())

		// Start the morning playlist
		spotify.Play(os.Getenv("MORNING_PLAYLIST"))
	})
	c.Start()

	log.Print("Gromit running on " + host)
	log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), nil))
}
