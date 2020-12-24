package spotify

import (
	"encoding/json"
	"github.com/timrcoulson/gromit/data"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"strings"
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
		panic(err)
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

func Play(uri string)  {
	devices, _ := client.PlayerDevices()
	for _, d := range devices {
		if d.Name == os.Getenv("SPOTIFY_DEVICE") {
			playlist := spotify.URI(uri)
			client.PlayOpt(&spotify.PlayOptions{DeviceID: &d.ID, PlaybackContext: &playlist})
		}
	}
}
