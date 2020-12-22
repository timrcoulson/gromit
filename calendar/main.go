package calendar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/olekukonko/tablewriter"
	"github.com/timrcoulson/gromit/printing"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type Calendar struct {

}

func (c *Calendar) Output() (output string)  {
	output = "# Calendar \n\n"

	obuf := bytes.NewBufferString("")
	table := tablewriter.NewWriter(obuf)

	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(Bod()).TimeMax(Eod()).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date, err := time.Parse(time.RFC3339, item.Start.DateTime)
			if err != nil {
				continue
			}
			table.Append([]string{date.Format("15:04"), printing.Clean(item.Summary)})
		}
	}

	table.SetHeader([]string{"Time", "Event"})
	table.SetAutoWrapText(false)
	table.Render()
	return output + obuf.String() + "\n"
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	if os.Getenv("GOOGLE_TOKEN") != "" {
		d1 := []byte(os.Getenv("GOOGLE_TOKEN"))
		err := ioutil.WriteFile("token.json", d1, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()

	json.NewEncoder(f).Encode(token)
	file, _ := os.Open(path)
	b, err := ioutil.ReadAll(file)
	fmt.Println(string(b))
}

func Eod() string {
	year, month, day := time.Now().Date()
	return time.Date(year, month, day + 1, 0, 0, 0, 0, time.Now().Location()).Format(time.RFC3339)
}

func Bod() string {
	year, month, day := time.Now().Date()
	return time.Date(year, month, day, 0, 0, 0, 0,  time.Now().Location()).Format(time.RFC3339)
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file")
	}
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON([]byte(os.Getenv("GOOGLE_JSON")), calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err = calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
}

var srv *calendar.Service