package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AssignSequence assigns a specific sequence to a client, instead of automatically assigning it based on the order of connection.

func AssignSequence(c *gin.Context) {
	// This function is a placeholder for assigning a sequence to a client.
	// The actual implementation will depend on the specific requirements of the application.
	r := RoomList[c.Param("roomId")]
	p := r.GetPlayerByDeviceID(c.Param("clientId"))
	if p == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Player not found",
		})
		return
	}

	// For now, it simply returns a success message.
	c.JSON(http.StatusOK, gin.H{
		"message": "Sequence assigned successfully",
	})
}
