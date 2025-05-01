package main

type User struct {
	Username string
	Password string
}

func NewUser(username, password string) (user User) {
	user.Username = username
	user.Password = password

	return user
}
