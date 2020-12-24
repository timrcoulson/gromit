package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"

	"github.com/zmb3/spotify"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8091/callback"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState)
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func main() {
	err := godotenv.Load()

	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8091", nil)

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)

	devices, err := client.PlayerDevices()

	log.Println("Devices:")
	log.Println(err)
	for _, device := range devices {
		if device.Name == "raspotify" {
			uri := spotify.URI("spotify:playlist:3nzXSd1mD8wCP9tIgZFLzX")
			err := client.PlayOpt(&spotify.PlayOptions{DeviceID: &device.ID, PlaybackContext: &uri})
			log.Println(err)
		}
	}
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)

	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}