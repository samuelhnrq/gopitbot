package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

var (
	discordKey = os.Getenv("DISCORD_TOKEN")
)

func main() {
	if discordKey == "" {
		fmt.Println("You need an API key")
		return
	}

	discord, err := discordgo.New(discordKey)

	if err != nil {
		fmt.Println("Fuck didn't work")
		return
	}

	discord.AddHandler(ready)

	err = discord.Open()
	if err != nil {
		fmt.Println("Didn't work again, shit: ", err)
	}

	fmt.Println("Bot's now running, press CTRL-C to close.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, "Maconha")
}
