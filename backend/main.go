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

	r.POST("/GetUser", func(c *gin.Context) {
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
	})

	r.Run(":5151")
}
