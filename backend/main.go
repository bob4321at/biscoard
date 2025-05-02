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

		chats := []Chat{}

		for _, chat := range temp_chat_list {
			for _, user := range chat.Users {
				if user == temp_data.Name {
					chats = append(chats, chat)
				}
			}
		}

		c.JSON(http.StatusAccepted, chats)
	})

	r.Run(":5151")
}
