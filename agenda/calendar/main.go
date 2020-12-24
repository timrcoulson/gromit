package calendar

import (
	"bytes"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/timrcoulson/gromit/agenda/printing"
	"github.com/timrcoulson/gromit/services/google"
	"log"
	"time"

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

func Eod() string {
	year, month, day := time.Now().Date()
	return time.Date(year, month, day + 1, 0, 0, 0, 0, time.Now().Location()).Format(time.RFC3339)
}

func Bod() string {
	year, month, day := time.Now().Date()
	return time.Date(year, month, day, 0, 0, 0, 0,  time.Now().Location()).Format(time.RFC3339)
}

func init() {
	client := google.Get()
	var err error
	srv, err = calendar.New(client)
	if err != nil {
		panic(err)
	}
}

var srv *calendar.Service