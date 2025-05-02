package main

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
