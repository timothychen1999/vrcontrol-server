package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/timothychen1999/vrcontrol-server/sockets"
)

func ConnectToRoomSocket(c *gin.Context) {
	//roomId := c.Param("roomId")
	deviceId := c.Param("clientId")
	roomId := "Test"
	room, ok := RoomList[roomId]
	if !ok {
		if len(RoomList) > MaxRoomCount {
			log.Println("Room List is full, please try again later.")
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Room List is full, please try again later."})
			return
		}
		room = sockets.NewRoom(roomId)
		RoomList[roomId] = room
		go room.Run()
		log.Println("Room Created: ", roomId)
	}
	conn, err := sockets.SocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error Upgrading Connection: ", err)
		return
	}
	_ = sockets.HandlePlayerConnect(room, conn, deviceId)

}
