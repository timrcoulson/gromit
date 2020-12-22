package gmail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/olekukonko/tablewriter"
	"github.com/timrcoulson/gromit/printing"
	"google.golang.org/api/gmail/v1"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const FromWidth = 30
const SubjectWidth = 80

type Gmail struct {

}

func (c *Gmail) Output() (output string)  {
	output = "# Emails \n\n"

	user := "me"
	r, err := srv.Users.Messages.List(user).LabelIds("INBOX").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve emails: %v", err)
	}
	if len(r.Messages) == 0 {
		output += "No emails. \n"
		return
	}


	obuf := bytes.NewBufferString("")
	table := tablewriter.NewWriter(obuf)
	table.SetHeader([]string{"From", "Subject"})
	table.SetAutoWrapText(false)
	table.SetColMinWidth(0, FromWidth)
	table.SetColMinWidth(1, SubjectWidth)

	for _, l := range r.Messages {
		msg, err := GetMessage(l.Id)
		if err != nil {
			log.Println(err)
			continue
		}
		header := ParseHeaders(msg.Payload.Headers)
		from :=  printing.Clean(SplitEmail(header["From"]))
		subject := printing.Clean(header["Subject"])

		if len(from) > FromWidth {
			from = from[:FromWidth]
		}
		if len(subject) > SubjectWidth {
			subject = subject[:SubjectWidth]
		}

		table.Append([]string{from, subject})
	}


	table.Render()
	return  output + obuf.String() + "\n"
}

func ParseHeaders(headers []*gmail.MessagePartHeader) map[string]string {
	headersMap := map[string]string{}
	for _, header := range headers {
		headersMap[header.Name] = header.Value
	}
	return headersMap
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	if os.Getenv("GOOGLE_GMAIL_TOKEN") != "" {
		d1 := []byte(os.Getenv("GOOGLE_GMAIL_TOKEN"))
		err := ioutil.WriteFile("token.json", d1, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "gmail-token.json"
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
	config, err := google.ConfigFromJSON([]byte(os.Getenv("GOOGLE_JSON")), gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err = gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
}

var srv *gmail.Service

func GetMessage(messageId string) (*gmail.Message, error) {
	resp, err := srv.Users.Messages.Get("me", messageId).Format("full").Do()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func SplitEmail(from string) string {
	add, _ := mail.ParseAddress(from)

	return add.Name
}