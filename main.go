package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	discord, err := discordgo.New(os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Println("Fuck didn't work")
		return
	}
	discord.AddHandler(ready)
	fmt.Println("Hello")

	err = discord.Open()
	if err != nil {
		fmt.Println("Didn't work again, shit: ", err)
	}

	fmt.Println("Bot's now running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, "Maconha")
}
