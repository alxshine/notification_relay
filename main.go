package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type notification struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title" binding:"required"`
	Message    string    `json:"message" binding:"required"`
	Tag        string    `json:"tag" binding:"required"`
	Expiration time.Time `json:"exp"`
}

// TODO: add coroutine that goes through and removes expired notifications
var notifications = []notification{
	{ID: uuid.New(), Title: "Test", Message: "This is a test message :)", Tag: "testing", Expiration: time.Now()},
	{ID: uuid.New(), Title: "Test2", Message: "This is another message", Tag: "testing", Expiration: time.Now()},
}

func getNotifications(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, notifications)
}

func postNotification(c *gin.Context) {
	var newNotification notification

	if err := c.BindJSON(&newNotification); err != nil {
		return
	}
	newNotification.ID = uuid.New() // override ID to avoid shenanigans and predictability

	// set default expiration if not set manually
	if newNotification.Expiration == (time.Time{}) {
		newNotification.Expiration = time.Now().Add(10 * time.Minute)
	}

	notifications = append(notifications, newNotification)
	c.IndentedJSON(http.StatusCreated, newNotification)
}

func cleanupExpiredNotifications() {
	notExpired := []notification{}

	now := time.Now()
	for _, n := range notifications {
		if n.Expiration.Local().After(now) {
			notExpired = append(notExpired, n)
		}

	}

	notifications = notExpired // I think because this is an assignment we don't have to use mutexes :)

}

func main() {
	go func() {
		for {
			time.Sleep(20 * time.Second)
			println("Cleaning up notifications")
			cleanupExpiredNotifications()
		}
	}()

	router := gin.Default()
	router.GET("/notifications", getNotifications)
	router.POST("/notification", postNotification)

	router.Run("localhost:8080")
}
