package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/a3ylf/geek-bot/video"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}
	var blxckwarningchannel string
	dsToken := os.Getenv("discord")
	ytToken := os.Getenv("ytb")
	sess, err := discordgo.New("Bot " + dsToken)
	if err != nil {
		log.Fatal(err)
	}

	lf := video.NewVideoFetcher(ytToken)
	blxckbro := false

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if strings.HasPrefix(m.Content, "addblxck!") && blxckbro {
			// Remove o prefixo "addblxck!" para capturar a string seguinte
			args := strings.TrimSpace(strings.TrimPrefix(m.Content, "addblxck!"))

			if args == "" {
				s.ChannelMessageSend(m.ChannelID, "kd o id bicho")
				return
			}
			_, err := video.FetchLatestVideo(args, ytToken)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "ESSE ID N EXISTE DOIDAO")
				return
			}
			lf.Ids = append(lf.Ids, args)
            s.ChannelMessageSend(m.ChannelID, "Added bro")
		}
		if m.Content == "blxck!" {
			blxckbro = true
			s.ChannelMessageSend(m.ChannelID, "blxck estará conosco aqui")
			blxckwarningchannel = m.ChannelID

			go func() {
				for {
					videos, err := lf.FetchLatestVideos()
					log.Println("Checando videos novos!")
					if err != nil {
						log.Println("Error fetching videos:", err)
						time.Sleep(1 * time.Minute)
						continue
					}
					if len(videos) > 0 {
						for _, video := range videos {
							s.ChannelMessageSend(blxckwarningchannel, fmt.Sprintf("Lançou musica nova:\n%s\nLink: %s",
								video.Title, video.Link))
						}
					}
					time.Sleep(time.Second * 300)
				}
			}()
		}

		if m.Content == "hello!" {
			s.ChannelMessageSend(m.ChannelID, "shut up nerd")
			time.Sleep(2 * time.Second)
			s.ChannelMessageSend(m.ChannelID, "Just kidding")
		}
		if m.Content == "roll!" {
			s.ChannelMessageSend(m.ChannelID, "Rolling! ...")
			time.Sleep(time.Second / 2)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Number: ", rand.Intn(5)+1))
		}
		if m.Content == "mewing!" {
			s.ChannelMessageSend(m.ChannelID, "https://tenor.com/view/mewing-snowman-mewing-snowman-mewing-streak-gif-8723555851120406070")
		}

		if strings.Index(fmt.Sprint(m.Mentions), fmt.Sprint(s.State.User)) != -1 {
			s.ChannelMessageSend(m.ChannelID, "Marca sua mamãe")
		}
		if m.Content == "baller!" {
			s.ChannelMessageSend(m.ChannelID, "https://tenor.com/view/roblox-baller-baller-roblox-roblox-meme-roblox-memes-gif-27316151")
		}

		if strings.Index(fmt.Sprint(m.Content), "choose!") != -1 {
			if len(m.Mentions) == 1 {
				s.ChannelMessageSend(m.ChannelID, "SÓ TEM UM CARA EU VOU ESCOLHER O QUE PRA VC PQP")
			} else if len(m.Mentions) == 0 {
				s.ChannelMessageSend(m.ChannelID, "Marca alguem imbecil idiota burro esquizofrenico  n sabe como a misera do comando funciona")
			} else {
				s.ChannelMessageSend(m.ChannelID, "The chosen user is: "+m.Mentions[rand.Intn(len(m.Mentions))].Mention())
			}
		}
		if strings.Index(fmt.Sprint(m.Content), "jordana") != -1 || strings.Index(fmt.Sprint(m.Content), "mlr") != -1 {
			s.ChannelMessageSendReply(m.ChannelID, "NÃO FALE DA MAMÃEZINHA DO ADM", (m.Reference()))
		}

	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer sess.Close()

	fmt.Println("Bot is online")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
