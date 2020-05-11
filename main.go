/*package main

import (
	"log"

	"gopkg.in/telegram-bot-api.v4"
)

func  main () {

//
	bot, err  :=  tgbotapi.NewBotAPI("1198461761:AAEb97q2Bjm474YaH3oBlrZiiSRcqq8sNEI")
	if  err  !=  nil {
		log.Panic (err)
	}

	bot . Debug  =  true

	log . Printf ( "Authorized on account% s" , bot . Self . UserName )

	u  :=  tgbotapi . NewUpdate ( 0 )
	u . Timeout  =  60

	updates , err  :=  bot . GetUpdatesChan ( u )

	for  update  :=  range  updates {
		if  update . Message  ==  nil {
			continue
		}

		log . Printf ( "[% s]% s" , update . Message . From . UserName , update . Message . Text )

		msg  :=  tgbotapi . NewMessage ( update . Message . Chat . ID , update . Message . Text )
		msg . ReplyToMessageID  =  update . Message . MessageID

		bot . Send ( msg )
	}
}*/


package main

import (
    "encoding/json"
    //"io"
    "io/ioutil"
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "gopkg.in/telegram-bot-api.v4"

    _ "github.com/heroku/x/hmetrics/onload"
    _ "github.com/lib/pq"
)

var (
    bot      *tgbotapi.BotAPI
  /*  botToken = "1198461761:AAEb97q2Bjm474YaH3oBlrZiiSRcqq8sNEI"
   // baseURL  = "https://<YOUR-APP-NAME>.herokuapp.com/"
    baseURL  = "https://may010app.herokuapp.com/"  */
)

func initTelegram() {
    var err error
	
   // bot, err = tgbotapi.NewBotAPI(botToken)
    bot, err = tgbotapi.NewBotAPI(os.Getenv("botToken"))
	
    if err != nil {
        log.Println(err)
        return
    }

    // this perhaps should be conditional on GetWebhookInfo()
    // only set webhook if it is not set properly
   
//   url := baseURL + bot.Token
    url := os.Getenv("baseURL") + bot.Token
    _, err = bot.SetWebhook(tgbotapi.NewWebhook(url))
    if err != nil {
        log.Println(err)
    } 
}

func webhookHandler(c *gin.Context) {
    defer c.Request.Body.Close()

    bytes, err := ioutil.ReadAll(c.Request.Body)
    if err != nil {
        log.Println(err)
        return
    }

    var update tgbotapi.Update
    err = json.Unmarshal(bytes, &update)
    if err != nil {
        log.Println(err)
        return
    }

    // to monitor changes run: heroku logs --tail
    log.Printf("From: %+v Text: %+v\n", update.Message.From, update.Message.Text)
}

func main() {
    port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set")
    }

    // gin router
    router := gin.New()
    router.Use(gin.Logger())

    // telegram
    initTelegram()
    router.POST("/" + bot.Token, webhookHandler)

    err := router.Run(":" + port)
    if err != nil {
        log.Println(err)
    }
}