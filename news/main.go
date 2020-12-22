package news

import (
	"bytes"
	"github.com/olekukonko/tablewriter"
	"github.com/timrcoulson/gromit/printing"
	"log"
	"net/http"
	"os"
	"github.com/robtec/newsapi/api"
)

type News struct {

}

const HeadlineWidth = 90


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
	opts := api.Options{Country: "gb", PageSize: 3, Page: 1, SortBy: "popularity"}

	// Get Top Headlines with options from above
	topHeadlines, err := client.TopHeadlines(opts)
	if err != nil {
		log.Fatal(err)
	}

	obuf := bytes.NewBufferString("")
	table := tablewriter.NewWriter(obuf)

	table.SetColMinWidth(0, HeadlineWidth)
	table.SetHeader([]string{"Title", "Published At"})
	table.SetAutoWrapText(false)

	for _, headline := range topHeadlines.Articles[:3] {
		title :=  printing.Clean(headline.Title)

		if len(title) > HeadlineWidth {
			title = title[:HeadlineWidth]
		}

		table.Append([]string{title, headline.PublishedAt})
	}

	table.Render()
	// Get Everything with options from above
	return "# News\n\n" + obuf.String()
}
