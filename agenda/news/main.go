package news

import (
	"bytes"
	"github.com/olekukonko/tablewriter"
	"github.com/timrcoulson/gromit/agenda/printing"
	"log"
	"net/http"
	"os"
	"github.com/robtec/newsapi/api"
	"strings"
)

type News struct {

}

const HeadlineWidth = 110


func (c *News) Output() (output string)  {
	httpClient := &http.Client{}
	key := os.Getenv("NEWS_API_KEY")
	url := "https://newsapi.org"

	// Create a client, passing in the above
	client, err := api.New(httpClient, key, url)
	if err != nil {
		log.Fatal(err)
	}

	// Create options for Ireland and Business'
	opts := api.Options{Country: "gb", PageSize: 10, Page: 1, SortBy: "popularity"}

	// Get Top Headlines with options from above
	topHeadlines, err := client.TopHeadlines(opts)
	if err != nil {
		log.Fatal(err)
	}

	obuf := bytes.NewBufferString("")
	table := tablewriter.NewWriter(obuf)

	table.SetColMinWidth(0, HeadlineWidth)
	table.SetHeader([]string{"Title"})
	table.SetAutoWrapText(false)

	for _, headline := range topHeadlines.Articles[:10] {
		title :=  printing.Clean(headline.Title)

		if shouldFilter(headline.Title) {
			continue
		}

		if len(title) > HeadlineWidth {
			title = title[:HeadlineWidth]
		}

		table.Append([]string{title})
	}

	table.Render()
	// Get Everything with options from above
	return "# News\n\n" + obuf.String() + "\n"
}

func shouldFilter(headline string) bool {
	// Filter out football ha
	for _, footballTerm := range []string{"Liverpool", "Man City", "Football", "Tottenham", "Arsenal"} {
		if strings.Contains(headline, footballTerm) {
			return true
		}
	}
	return false
}
