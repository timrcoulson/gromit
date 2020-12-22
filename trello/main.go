package trello

import (
	"bytes"
	"github.com/adlio/trello"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
)

type Trello struct {

}

func (t *Trello) Output() (output string) {
	output += "# Tasks\n\n"
	board, err := client.GetBoard(os.Getenv("TRELLO_BOARD_ID"), trello.Defaults())
	if err != nil {
		log.Fatal(err)
	}

	lists, err := board.GetLists(trello.Defaults())

	content := make(map[string][]string)
	var headers []string
	var maxLength int
	for _, list := range lists {
		var listRow []string
		// GetCards makes an API call to /lists/:id/cards using credentials from `client`
		headers = append(headers, list.Name)
		cards, err := list.GetCards(trello.Defaults())

		if err != nil {
			log.Fatal(err)
		}

		for _, card := range cards {
			listRow = append(listRow, card.Name)
		}

		if len(cards) > maxLength {
			maxLength = len(cards)
		}

		content[list.Name] = listRow
	}

	obuf := bytes.NewBufferString("")

	table := tablewriter.NewWriter(obuf)
	table.SetHeader(headers)
	table.SetAutoWrapText(false)

	rows := [][]string{}
	i := 0
	for i <= maxLength {
		row := []string{}
		for _, list := range content {
			if len(list) > i {
				row = append(row, list[i])
			}
		}
		rows = append(rows, row)

		i++
	}
	table.AppendBulk(rows)

	table.Render() // Send output

	output += obuf.String()
	return output + "\n"
}

func init()  {
	client = trello.NewClient(os.Getenv("TRELLO_KEY"), os.Getenv("TRELLO_TOKEN"))

}

var client *trello.Client