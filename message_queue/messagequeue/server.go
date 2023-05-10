package messagequeue

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

var (
	mq                     = newMessageQueue()
	DEFAULT_TOPIC_CAPACITY = 100
	upgrader               = websocket.Upgrader{}
)

type SubcribeQuery struct {
	Topic string `form:"topic" json:"topic" binding:"required"`
}

type ProduceJson struct {
	Topic   string `form:"topic" json:"topic" binding:"required"`
	Message string `form:"message" json:"message" binding:"required"`
}

func subscribe(c *gin.Context) {
	var param SubcribeQuery
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !mq.HasTopic(param.Topic) {
		mq.AddTopic(param.Topic, DEFAULT_TOPIC_CAPACITY)
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()
	for {
		message, err := mq.Get(param.Topic)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			fmt.Printf("%v\n", message.Content)
			err = ws.WriteMessage(websocket.TextMessage, []byte(message.Content))
			if err != nil {
				fmt.Println(err)
				break
			}
		}
	}
}

func produce(c *gin.Context) {
	var jsonData ProduceJson
	if err := c.BindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(jsonData)
	if !mq.HasTopic(jsonData.Topic) {
		mq.AddTopic(jsonData.Topic, DEFAULT_TOPIC_CAPACITY)
	}
	mq.Put(jsonData.Topic, &Message{jsonData.Message})
	c.JSON(http.StatusOK, true)
}

func RunServer() {
	viper.AddConfigPath("./")
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Error reading config.json")
	}

	router := gin.Default()
	router.GET("/subcribe", subscribe)
	router.POST("/produce", produce)

	router.Run(fmt.Sprintf("%v:%v", viper.Get("server.host"), viper.Get("server.port")))
}
