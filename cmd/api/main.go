package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
    var router *gin.Engine = gin.Default()

    // Do not trust proxy headers (e.g. X-Forwarded-For); use the direct client connection IP.
    router.SetTrustedProxies(nil)

    router.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Todo API is running!",
            "status": "Success",
        })
    })

    router.Run(":8080")
}
