package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/AddUser", AddUser)
	r.POST("/GetUser", GetUser)

	r.POST("/MakeChat", MakeChatNetworked)
	r.POST("/GetChatsForUser", GetChatsForUser)
	r.POST("/SendMessege", SendMessegeNetworked)

	r.Run(":5151")
}
