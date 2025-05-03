package main

import (
	"encoding/json"
	"io"
	"net/http"
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

func (c *Chat) SendMessege(username, messege string) {
	c.Messeges = append(c.Messeges, NewMessege(username, messege))
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

func GetChatsForUser(c *gin.Context) {
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
}

func SendMessegeNetworked(c *gin.Context) {
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

	chat_to_operate := &Chat{}

	for chati := range temp_chat_list {
		if temp_data.Name == temp_chat_list[chati].Name {
			chat_to_operate = &temp_chat_list[chati]
		}
	}

	chat_to_operate.SendMessege(temp_data.Username[0], temp_data.Username[1])

	os.Remove("./chats.json")

	f, err := os.Create("./chats.json")
	if err != nil {
		panic(err)
	}

	chat_string_data, err := json.Marshal(temp_chat_list)
	if err != nil {
		panic(err)
	}

	f.Write(chat_string_data)

	f.Close()
}

func AddUserToChatNetworked(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}

	temp_data := NetworkChat{}
	if err := json.Unmarshal(data, &temp_data); err != nil {
		panic(err)
	}

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

	for tci := range temp_chat_list {
		if temp_chat_list[tci].Name == temp_data.Name {
			for _, temp_name := range temp_chat_list[tci].Users {
				if temp_name == temp_data.Username[0] {
					return
				}
			}
			temp_chat_list[tci].Users = append(temp_chat_list[tci].Users, temp_data.Username[0])
		}
	}

	os.Remove("./chats.json")

	f, err = os.Create("./chats.json")
	if err != nil {
		panic(err)
	}

	chat_string_data, err := json.Marshal(temp_chat_list)
	if err != nil {
		panic(err)
	}

	f.Write(chat_string_data)

	f.Close()
}
