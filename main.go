package main

import (
	"github.com/AhmadAboElzahab/bridge/initializers"
	"github.com/gin-gonic/gin"
)
func init(){
	initializers.LoadENV()
	initializers.ConnectDatabase()
}
func main() {
  router := gin.Default()
  router.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })
  router.Run() // listen and serve on 0.0.0.0:8080
}
