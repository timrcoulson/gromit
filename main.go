package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/robfig/cron/v3"
	"github.com/timrcoulson/gromit/agenda"
	"github.com/timrcoulson/gromit/printer"
	"github.com/timrcoulson/gromit/services/google"
	"github.com/timrcoulson/gromit/services/spotify"
	"log"
	"net/http"
	"os"
)

func main()  {
	host := os.Getenv("HOST")
	fmt.Println("Starting up gromit...")

	// Authorise all the services
	sp := spotify.New(host + "/auth/spotify")
	gg := google.New(host + "/auth/google")

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`
		<h1>Gromit is running!</h1>
		<a href="%s">Authorize Google</a>
		<a href="%s">Authorize Spotify</a>
		
		`, gg.LoginUrl(), sp.LoginUrl())))
	})

	r.HandleFunc("/spotify", func(writer http.ResponseWriter, request *http.Request) {
		devices, _ := spotify.Get().PlayerDevices()

		writer.Write([]byte(devices[0].Name))
	})
	r.HandleFunc("/print", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("printing"))
		printer.Print(agenda.Today())
	})

	// Register oauth 2 callbacks
	r.HandleFunc("/auth/spotify", sp.Callback)
	r.HandleFunc("/auth/google", gg.Callback)

	http.Handle("/", r)

	// Print the agenda every day
	c := cron.New()
	c.AddFunc("00 06 * * *", func() {
		printer.Print(agenda.Today())
	})
	c.Start()

	log.Print("Gromit running on http://localhost:" + os.Getenv("PORT"))
	http.ListenAndServe(":" + os.Getenv("PORT"), nil)
}
