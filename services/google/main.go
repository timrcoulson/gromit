package google

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/timrcoulson/gromit/data"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	state string
	oauthConfig *oauth2.Config
	client *http.Client
	tokenSource *oauth2.TokenSource
)
var tok *oauth2.Token


func New(redirectUri string) Google {
	oauthConfig = &oauth2.Config{
		RedirectURL:  redirectUri,
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{calendar.CalendarReadonlyScope, gmail.GmailReadonlyScope},
		Endpoint:     google.Endpoint,
	}

	state = "abc123"

	return Google{}
}

type Google struct {}

func (s *Google) LoginUrl() string  {
	url := oauthConfig.AuthCodeURL(state)
	return url
}

func (s *Google) Callback(w http.ResponseWriter, r *http.Request)  {
	tok, err := getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	out, _ := json.Marshal(tok)
	data.Set("google-token", string(out))

	// use the token to get an authenticated client
	client = oauthConfig.Client(context.Background(), tok)

	w.Write([]byte("Authorised"))
}


func getUserDataFromGoogle(code string) (*oauth2.Token, error) {
	// Use code to get token and get user info from Google.
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	return token, err
}

func init()  {
	token, err := data.Get("google-token")
	if err == nil {
		err = json.NewDecoder(strings.NewReader(token)).Decode(&tok)

		if err != nil {
			panic(err)
		}

		if !tok.Valid() {
			tokenSource := oauthConfig.TokenSource(oauth2.NoContext, tok)
			client = oauth2.NewClient(oauth2.NoContext, tokenSource)

			//tok, err = tokenSource.Token()
		} else {

			client = oauthConfig.Client(context.Background(), tok)
		}

	} else {
		panic(err)
	}
}

func Get() (*http.Client, oauth2.TokenSource){
	return client, oauthConfig.TokenSource(context.Background(), tok)
}