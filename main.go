package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
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
	queque     = make([]song, 0)
	admin      = "ADMIN"
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
	discord = s
	s.UpdateStatus(0, "Saint iGNUcius")
	roles, err := s.GuildRoles(guildID)
	pErr(err)
	for _, val := range roles {
		if strings.Contains(strings.ToLower(val.Name), "admin") {
			admin = val.ID
			break
		}
	}
}

func message(s *discordgo.Session, event *discordgo.MessageCreate) {
	chatCh = event.ChannelID
	args := strings.Split(event.Content, " ")
	memb, err := s.GuildMember(guildID, event.Author.ID)
	pErr(err)
	adm := false
	for _, val := range memb.Roles {
		adm = val == admin
	}
	adm = adm || (currSong.owner == event.Author.ID)
	if len(args) > 0 {
		if args[0] == "!play" && len(args) > 1 {
			fmt.Println("Connecting and streaming " + args[1])
			input := args[1]
			url, title, err := GetVideoDownloadURL(input)
			if err == nil {
				sng := song{event.Author.ID, title, url}
				if run != nil {
					sendMsg("Adicionado a fila: \"**" + title + "**\"")
					queque = append(queque, sng)
					return
				}
				sendMsg("Agora tocando \"**" + title + "**\"")
				currSong = sng
				playVideo(dgv, url)
			}
		} else if args[0] == "!skip" {
			if !adm {
				sendMsg("SEM PERMISSAO SEU BADERNEIRO")
				return
			}
			if run == nil {
				sendMsg("Nenhuma música na fila")
				return
			}
			run.Process.Kill()
			sendMsg("Pulado \"**" + currSong.title + "**\"")
		} else if args[0] == "!queque" {
			msg := ""
			for k, v := range queque {
				msg += strconv.Itoa(k+1) + "\"**" + v.title + "**\"\n"
			}
			sendMsg(msg)
		} else if args[0] == "!volume" || args[0] == "!stop" {
			sendMsg("Botão direito no bot pra controlar o volume")
		} else if args[0] == "!song" {
			sendMsg("Música tocando: \"**" + currSong.title + "**\"")
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

func pErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func sendMsg(msg string) {
	_, err := discord.ChannelMessageSend(chatCh, msg)
	pErr(err)
}

type song struct {
	owner string
	title string
	url   string
}
