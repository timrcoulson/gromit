package money

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"net/http"
	"os"
)

type Money struct {

}

func (m *Money) Output() string {
	c := http.DefaultClient
	rt := WithHeader(c.Transport)
	rt.Set("Authorization", "Bearer " + os.Getenv("STARLING_PERSONAL_ACCESS_TOKEN"))
	c.Transport = rt

	type Account struct {
		Id string `json:"accountUid"`
		Name string `json:"name"`
	}
	type Accounts struct {
		Accounts []Account `json:"accounts"`
	}
	type Money struct {
		Currency string `json:"currency"`
		MinorUnits int `json:"minorUnits"`
	}
	type Balance struct {
		ClearedBalance Money `json:"clearedBalance"`
		EffectiveBalance Money `json:"effectiveBalance"`
		Amount Money `json:"amount"`
	}
	
	resp, err := c.Get("https://api.starlingbank.com/api/v2/accounts")

	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	var accounts Accounts
	json.Unmarshal(b, &accounts)

	obuf := bytes.NewBufferString("")

	table := tablewriter.NewWriter(obuf)
	table.SetHeader([]string{"Account", "Amount"})
	table.SetAutoWrapText(false)

	for _, a := range accounts.Accounts {
		resp, err = c.Get(fmt.Sprintf("https://api.starlingbank.com/api/v2/accounts/%s/balance", a.Id))
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var balance Balance
		json.Unmarshal(b, &balance)

		table.Append([]string{a.Name, fmt.Sprintf("%.2f", float32(balance.Amount.MinorUnits) / 100.0)})
	}

	table.Render()

	return "# Money\n\n" + obuf.String() + "\n"
}

type withHeader struct {
	http.Header
	rt http.RoundTripper
}

func WithHeader(rt http.RoundTripper) withHeader {
	if rt == nil {
		rt = http.DefaultTransport
	}

	return withHeader{Header: make(http.Header), rt: rt}
}

func (h withHeader) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range h.Header {
		req.Header[k] = v
	}

	return h.rt.RoundTrip(req)
}
