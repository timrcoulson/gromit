package gmail

import (
	"bytes"
	"github.com/olekukonko/tablewriter"
	"github.com/timrcoulson/gromit/agenda/printing"
	"github.com/timrcoulson/gromit/services/google"
	"google.golang.org/api/gmail/v1"
	"log"
	"net/mail"
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

func init() {
	client := google.Get()
	var err error
	srv, err = gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail client: %v", err)
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