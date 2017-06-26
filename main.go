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
	chatCh     = ""
	quequeList = make([]song, 0)
	admin      = "ADMIN"
	discord    *discordgo.Session
	dgv        *discordgo.VoiceConnection
)

type song struct {
	owner string
	title string
	url   string
}

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
	discord = s
	s.UpdateStatus(0, "Saint iGNUcius")
	roles, err := s.GuildRoles(guildID)
	pErr(err)
	for _, val := range roles {
		fmt.Println(val.Name, val.ID)
		if strings.Contains(strings.ToLower(val.Name), "admin") {
			admin = val.ID
			break
		}
	}
}

func message(s *discordgo.Session, event *discordgo.MessageCreate) {
	chatCh = event.ChannelID
	args := strings.Split(event.Content, " ")
	if len(args) <= 0 {
		return
	}
	cmd := args[0]
	if cmd[0] != '!' {
		return
	}
	cmd = cmd[1:]
	memb, err := s.GuildMember(guildID, event.Author.ID)
	pErr(err)
	permiss := false
	for _, val := range memb.Roles {
		permiss = val == admin
	}
	permiss = permiss || (currSong.owner == event.Author.ID)
	switch cmd {
	case "play":
		if len(args) <= 1 {
			break
		}
		play(args[1], event.Author.ID)
	case "skip":
		skip(permiss)
	case "queque":
		printQueque()
	case "volume", "stop":
		sendMsg("Botão direito no bot pra controlar o volume")
	case "song":
		sendMsg("Música tocando: \"**" + currSong.title + "**\"")
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

func pErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func sendMsg(msg string) {
	_, err := discord.ChannelMessageSend(chatCh, msg)
	pErr(err)
}
