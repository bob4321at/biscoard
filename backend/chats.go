package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

type Messege struct {
	Username string
	Messege  string
}

func NewMessege(username, messege_content string) (messege Messege) {
	messege.Username = username
	messege.Messege = messege_content

	return messege
}

type Chat struct {
	Name     string
	Users    []string
	Messeges []Messege
}

type NetworkChat struct {
	Name     string
	Username []string
}

func NewChat(name string, start_users []string) (chat Chat) {
	chat.Name = name
	chat.Users = start_users

	return chat
}

func (chat *Chat) AddMessege(user, messege string) {
	for _, username := range chat.Users {
		if username == user {
			chat.Messeges = append(chat.Messeges, NewMessege(user, messege))
		}
	}
}

func MakeChatNetworked(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}

	temp_data := NetworkChat{}
	if err := json.Unmarshal(data, &temp_data); err != nil {
		panic(err)
	}

	if _, err := os.Stat("./chats.json"); err != nil {
		f, err := os.Create("./chats.json")
		if err != nil {
			panic(err)
		}

		chats := []Chat{}
		chats = append(chats, NewChat(temp_data.Name, temp_data.Username))

		temp_data, err := json.Marshal(chats)
		if err != nil {
			panic(err)
		}

		f.Write(temp_data)

		f.Close()
	} else {
		f, err := os.Open("./chats.json")
		if err != nil {
			panic(err)
		}

		previous_data, err := os.ReadFile(f.Name())
		if err != nil {
			panic(err)
		}

		temp_chat_list := []Chat{}

		err = json.Unmarshal(previous_data, &temp_chat_list)
		if err != nil {
			panic(err)
		}

		temp_chat_list = append(temp_chat_list, NewChat(temp_data.Name, temp_data.Username))

		temp_string_data, err := json.Marshal(temp_chat_list)
		if err != nil {
			panic(err)
		}

		f.Close()
		os.Remove("./chats.json")

		f, err = os.Create("./chats.json")
		if err != nil {
			panic(err)
		}

		f.Write(temp_string_data)

		f.Close()
	}
}
