package main

import (
	"fmt"
	"strconv"
)

func play(url, author string) {
	url, title, err := GetVideoDownloadURL(url)
	if err == nil {
		sng := song{author, title, url}
		if run != nil {
			sendMsg("Adicionado a fila: \"**" + title + "**\"")
			quequeList = append(quequeList, sng)
			return
		}
		sendMsg("Agora tocando \"**" + title + "**\"")
		currSong = sng
		fmt.Println("Connecting and streaming " + url)
		playVideo(dgv, url)
	} else {
		sendMsg("URL inválida")
	}
}

func skip(adm bool) {
	if !adm {
		sendMsg("A música não é sua cuzão!")
		return
	}
	if run == nil {
		sendMsg("Nenhuma música na fila")
		return
	}
	run.Process.Kill()
	sendMsg("Pulado \"**" + currSong.title + "**\"")
}

func printQueque() {
	msg := ""
	for k, v := range quequeList {
		msg += strconv.Itoa(k+1) + "\"**" + v.title + "**\"\n"
	}
	sendMsg(msg)
}
