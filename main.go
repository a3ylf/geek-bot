package main

import (
	"context"
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
	"github.com/chromedp/chromedp"
	"github.com/joho/godotenv"
)

func blackingtime() (string, string) {
	// Configure o contexto com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Configure o Chrome (use o caminho correto se necessário)
	ctx, cancel = chromedp.NewExecAllocator(ctx, chromedp.DefaultExecAllocatorOptions[:]...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// URL do canal do YouTube (substitua pelo canal desejado)
	url := "https://www.youtube.com/@blxckoficial/videos"

	// Variáveis para armazenar os resultados
	var videoTitles []string
	var videoLinks []string

	// Execute o scraping
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),            // Acessa a página do canal
		chromedp.WaitVisible(`#contents`), // Aguarda o carregamento dos vídeos

		// Extrai os títulos dos vídeos
		chromedp.Evaluate(`Array.from(document.querySelectorAll("#video-title")).map(e => e.textContent.trim())`, &videoTitles),

		// Extrai os links dos vídeos
		chromedp.Evaluate(`Array.from(document.querySelectorAll("#video-title-link")).map(e => "https://www.youtube.com" + e.getAttribute("href"))`, &videoLinks),
	)
	if err != nil {
		log.Fatal("Failed to scrape YouTube:", err)
	}

	// Exibe os títulos e links

	return videoTitles[0], videoLinks[0]
}

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

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if m.Content == "blxck!" {
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
