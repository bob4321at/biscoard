package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type User struct {
	Username string
	Password string
}

func NewUser(username, password string) (user User) {
	user.Username = username
	user.Password = password

	return user
}

func AddUser(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}

	temp_user := User{}

	if err := json.Unmarshal(data, &temp_user); err != nil {
		panic(err)
	}

	c.Status(http.StatusAccepted)

	if _, err := os.Stat("./users.json"); err != nil {
		f, err := os.Create("./users.json")
		if err != nil {
			panic(err)
		}

		users := []User{}
		users = append(users, temp_user)

		json_file, err := json.Marshal(users)
		if err != nil {
			panic(err)
		}

		f.Write([]byte(json_file))
		f.Close()
	} else {
		f, err := os.Open("./users.json")
		if err != nil {
			panic(err)
		}

		previous_data, err := os.ReadFile(f.Name())
		if err != nil {
			panic(err)
		}

		temp_user_list := []User{}
		err = json.Unmarshal(previous_data, &temp_user_list)
		if err != nil {
			panic(err)
		}

		temp_user_list = append(temp_user_list, temp_user)

		f.Close()
		os.Remove("./users.json")

		f, err = os.Create("./users.json")
		if err != nil {
			panic(err)
		}

		user_list_data, err := json.Marshal(temp_user_list)

		f.Write([]byte(user_list_data))
	}
}

func GetUser(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}

	temp_user := User{}

	if err := json.Unmarshal(data, &temp_user); err != nil {
		panic(err)
	}

	if _, err := os.Stat("./users.json"); err != nil {
		c.Status(http.StatusBadRequest)
	} else {
		previous_data, err := os.ReadFile("./users.json")
		if err != nil {
			panic(err)
		}

		temp_user_list := []User{}
		err = json.Unmarshal(previous_data, &temp_user_list)
		if err != nil {
			panic(err)
		}

		for _, user := range temp_user_list {
			if temp_user.Username == user.Username && temp_user.Password == user.Password {
				c.JSON(http.StatusAccepted, user)
			}
		}

		c.Status(http.StatusBadRequest)
	}
}
