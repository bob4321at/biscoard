package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/AddUser", AddUser)
	r.POST("/GetUser", GetUser)

	r.POST("/MakeChat", MakeChatNetworked)
	r.POST("/GetChatsForUser", func(c *gin.Context) {
		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}

		temp_data := NetworkChat{}

		if err := json.Unmarshal(data, &temp_data); err != nil {
			panic(err)
		}

		previous_data, err := os.ReadFile("./chats.json")
		if err != nil {
			c.Status(http.StatusBadRequest)
		}

		temp_chat_list := []Chat{}

		err = json.Unmarshal(previous_data, &temp_chat_list)
		if err != nil {
			panic(err)
		}

		for _, chat := range temp_chat_list {
			if chat.Name == temp_data.Name {
				c.JSON(http.StatusAccepted, chat)
			}
		}
	})

	r.Run(":5151")
}
