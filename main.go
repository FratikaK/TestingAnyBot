package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
)

var s *discordgo.Session

func init() {
	flag.Parse()
}

func init() {
	var err error
	s, err = discordgo.New("Bot Token")
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "create",
			Description: "パーティをつくるよ！",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"create": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "パーティ！",
				},
			})
			if err != nil {
				log.Fatalf("Invalid create command parameters: %v", err)
				return
			}
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

var users []string

func contains(array []string, str string) bool {
	for _, v := range array {
		if v == str {
			return true
		}
	}
	return false
}

func voiceHandler(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	if !contains(users, v.UserID) {
		users = append(users, v.UserID)
	}
	if v.ChannelID == "" {
		fmt.Println("削除したよ！")

	}
	fmt.Println("==============各種パラメータ==============")
	fmt.Println(v.ChannelID)
	fmt.Println(v.UserID)
	fmt.Println("Mute", v.Mute, "\n"+
		"SelfMute", v.SelfMute, "\n"+
		"Deaf", v.Deaf, "\n"+
		"Suppress", v.Suppress)
	fmt.Println("=========================================")
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot Start!")
	})
	s.AddHandler(voiceHandler)

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	for _, v := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	defer s.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Bot ShutDowning...")
}
