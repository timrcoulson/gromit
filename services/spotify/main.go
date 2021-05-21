package spotify

import (
	"encoding/json"
	"github.com/timrcoulson/gromit/data"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	auth *spotify.Authenticator
	client *spotify.Client
	state string
)

func New(redirectUri string) Spotify {
	state = "abc123"
	authV := spotify.NewAuthenticator(redirectUri, spotify.ScopeUserReadPrivate, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState)
	auth = &authV

	token, err := data.Get("spotify-token")
	if err == nil {
		tok := &oauth2.Token{}
		err = json.NewDecoder(strings.NewReader(token)).Decode(tok)

		if err != nil {
			panic(err)
		}

		clientV := auth.NewClient(tok)
		client = &clientV
	} else {
		log.Println(err)
	}


	return Spotify{}
}

type Spotify struct {}

func (s *Spotify) LoginUrl() string  {
	url := auth.AuthURL(state)
	return url
}

func (s *Spotify) Callback(w http.ResponseWriter, r *http.Request)  {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}

	out, _ := json.Marshal(tok)
	data.Set("spotify-token", string(out))

	// use the token to get an authenticated client
	clientV := auth.NewClient(tok)
	client = &clientV

	w.Write([]byte("Authorised"))
}

func Get() *spotify.Client {
	return client
}

func Play(uri string, mins int64)  {
	go exec.Command("killall", "-s", "9", "omxplayer.bin").Run()

	devices, _ := client.PlayerDevices()
	for _, d := range devices {
		if d.Name == os.Getenv("SPOTIFY_DEVICE") {
			playlist := spotify.URI(uri)

			splitUri := strings.Split(uri, ":")
			id := splitUri[len(splitUri)- 1]

			log.Printf("id %v", id)
			pl, err := client.GetPlaylist(spotify.ID(splitUri[len(splitUri)- 1]))

			if err != nil {
				log.Println(err)
			}

			log.Printf("%v tracks", pl.Tracks.Total)
			offset := rand.Intn(pl.Tracks.Total)

			var pbo *spotify.PlaybackOffset
			if offset != 0 {
				pbo = &spotify.PlaybackOffset{Position: offset}
			}

			log.Printf("playing %v", strconv.Itoa(offset))
			err = client.PlayOpt(&spotify.PlayOptions{DeviceID: &d.ID, PlaybackContext: &playlist, PlaybackOffset: pbo})

			if err != nil {
				log.Println(err)
			}
		}
	}

	t := time.Tick(time.Duration(mins) * time.Minute)
	<- t

	client.Pause()
}


func PlaySingle(uri string)  {
	go exec.Command("killall", "-s", "9", "omxplayer.bin").Run()

	devices, _ := client.PlayerDevices()
	for _, d := range devices {
		if d.Name == os.Getenv("SPOTIFY_DEVICE") {
			song := spotify.URI(uri)
			client.PlayOpt(&spotify.PlayOptions{DeviceID: &d.ID, URIs: []spotify.URI{song}})
		}
	}

	t := time.Tick(5 * time.Second)

	for {
		<- t

		state, _ := client.PlayerState()
		if !state.Playing {
			return
		}
	}
}
