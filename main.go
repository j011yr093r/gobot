package main

import (
	"./pkg"
	"flag"
	"github.com/Syfaro/telegram-bot-api"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var conf *config.Config
var mess config.MsgDB

func RandomizeAnswers(Collection []string) string { // Returns the random value from an array
	s1 := rand.NewSource(time.Now().UnixNano())
	AnsIndex := rand.New(s1).Intn(len(Collection))
	return Collection[AnsIndex]
}

func init() {
	const (
		defaultConfigFilename = "./configs/conf.yaml"
		defaultMessagesFilename = "./configs/messages.yaml"
	)
	var err error
	var configFilename string
	flag.StringVar(&configFilename, "config", defaultConfigFilename, "the config filename") // Optional app [configfile] launch option for custom configuration
	flag.Parse()
	if err = initConf(configFilename); err != nil {
		log.Panic("Can't init config: %s", err.Error())
	}
	if err = initMsgDB(defaultMessagesFilename); err != nil {
                log.Panic("Can't init messages storage: %s", err.Error())
        }
}
func initConf(filename string) (err error) { // Importing config structure
	if conf, err = config.NewConf(filename); err != nil {
		return
	}
	return
}

func initMsgDB(filename string) (err error) { // Importing messages hash table
        if mess, err = config.NewMsgDB(filename); err != nil {
                return
        }
        return
}

func main() {
	WorkAnswers := []string{"Ты хотел сказать \"чиллить\"?", "Ну я лично скоро на Охотном чилить буду", "Опять на работу пиздос"}
	FactAnswers := []string{"То что в рашке нет конституции это факт", "Факт это когда ты неправ короче", "Если тебе нужен факт - чекни лс"}
	KFCAnswers := []string{"3870 чекай кстати, нидораха", "Ну я в канал абузы кинул, так что чекай", "Ну и что что говно, #затонидораха"}
	proxyUrl, err := url.Parse("socks5://127.0.0.1:9050") // Proxy pass
	if err != nil {
		log.Panic(err)
	}
	transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	client := &http.Client{Transport: transport}
	// Auth section
	bot, err := tgbotapi.NewBotAPIWithClient(conf.Conf.Token, client)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account: %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "about":
				msg.Text = "Меня зовут Гореси. Если вы считаете, что я имею отношение к любому реальному человеку, или являюсь какой-либо отсылкой, то вы, наверное, долбоёб"
			case "abuse":
				msg.Text = "Список актуальных абуз KFC: 3870 - боксмастер, картошка, напиток 300мл - 184р (ТОП); 3869 - любой кофе 200мл - 45р; 3878 - твистер, картошка, напиток 400мл - 169р "
			default:
				msg.Text = "Error: ты пидор, нет такой команды"
			}
			bot.Send(msg)
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		switch { // TODO: refactor this section. Use map[string][]string instead of this crappy code below
		case strings.Contains(msg.Text, "работать"):
			msg.ReplyToMessageID = update.Message.MessageID
			msg.Text = RandomizeAnswers(WorkAnswers)
			bot.Send(msg)
		case strings.Contains(msg.Text, "финик"):
			msg.ReplyToMessageID = update.Message.MessageID
			msg.Text = "Ну вообще-то финик заебал меня хейтить"
			bot.Send(msg)
		case strings.Contains(msg.Text, "залупин"):
			msg.ReplyToMessageID = update.Message.MessageID
			msg.Text = "Залупин бы сначала научился микротик настраивать"
			bot.Send(msg)
		case strings.Contains(msg.Text, "справедливо"):
			msg.ReplyToMessageID = update.Message.MessageID
			msg.Text = "имхо нет"
			bot.Send(msg)
		case strings.Contains(msg.Text, "хорошо плохо"):
			msg.ReplyToMessageID = update.Message.MessageID
			msg.Text = "Ну они с фиником вообще-то меня тупо хейтят постоянно"
			bot.Send(msg)
		case strings.Contains(msg.Text, "факт"):
			msg.ReplyToMessageID = update.Message.MessageID
			msg.Text = RandomizeAnswers(FactAnswers)
			bot.Send(msg)
		case strings.Contains(msg.Text, "KFC"):
			msg.ReplyToMessageID = update.Message.MessageID
			msg.Text = RandomizeAnswers(KFCAnswers)
			bot.Send(msg)
		}
	}
}
