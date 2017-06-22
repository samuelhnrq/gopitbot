package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

var (
	discordKey = os.Getenv("DISCORD_TOKEN")
	channelID  = os.Getenv("DISCORD_CHANNEL")
	guildID    = os.Getenv("DISCORD_GUILD")
	discord    *discordgo.Session
	dgv        *discordgo.VoiceConnection
)

func main() {
	if discordKey == "" {
		fmt.Println("You need an API key")
		return
	}

	discord, err := discordgo.New("Bot " + discordKey)
	if err != nil {
		fmt.Println("Fuck didn't work, reason: ", err.Error())
		return
	}

	discord.AddHandler(ready)
	discord.AddHandler(message)

	err = discord.Open()
	if err != nil {
		fmt.Println("Didn't work again, shit: ", err.Error())
	}

	dgv, err = discord.ChannelVoiceJoin(guildID, channelID, false, false)
	if err != nil {
		fmt.Println(err)
		return
	}

	voiceSetup(dgv)

	fmt.Println("Bot's now running, press CTRL-C to close.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dgv.Close()
	discord.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, "Maconha")
}

func message(s *discordgo.Session, event *discordgo.MessageCreate) {
	args := strings.Split(event.Content, " ")
	if len(args) > 0 {
		if args[0] == "!play" {
			if len(args) > 1 {
				fmt.Println("Connecting and streaming " + args[1])
				input := args[1]
				url, title, err := GetVideoDownloadURL(input)
				if err == nil {
					_, err := s.ChannelMessageSend(event.ChannelID, "Playing: "+title)
					if err != nil {
						fmt.Println(err.Error())
					}
					if run != nil {
						run.Process.Kill()
					}
					playVideo(dgv, url)
				}
			}
		}
	}
}

func echo(v *discordgo.VoiceConnection) {
	recv := make(chan *discordgo.Packet, 2)
	go dgvoice.ReceivePCM(v, recv)

	send := make(chan []int16, 2)
	go dgvoice.SendPCM(v, send)

	v.Speaking(true)
	defer v.Speaking(false)

	for {
		p, ok := <-recv
		if !ok {
			return
		}

		send <- p.PCM
	}
}
