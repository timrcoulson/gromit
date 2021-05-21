package agenda

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/timrcoulson/gromit/agenda/calendar"
	"github.com/timrcoulson/gromit/agenda/gmail"
	"github.com/timrcoulson/gromit/agenda/money"
	"github.com/timrcoulson/gromit/agenda/news"
	"github.com/timrcoulson/gromit/agenda/trello"
	"log"
	"time"
)

func Today() string  {
	fmt.Println("Starting up gromit...")

	// Register modules
	var modules []Module
	modules = append(modules, &calendar.Calendar{})
	modules = append(modules, &gmail.Gmail{})
	modules = append(modules, &trello.Trello{})
	modules = append(modules, &money.Money{})
	modules = append(modules, &news.News{})

	output := fmt.Sprintf("=== Good Morning, Tim. Today is %v ===\n\n", time.Now().Format("Monday 2 Jan 2006"))
	for _, module := range modules {
		output += module.Output()
	}

	return output
}

type Module interface {
	Output() string
}

func init()  {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
}