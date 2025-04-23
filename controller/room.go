package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/timothychen1999/vrcontrol-server/sockets"
)

const MaxRoomCount = 10

var RoomList map[string]*sockets.Room = make(map[string]*sockets.Room)

func GetRoomList(c *gin.Context) {
	lis := make([]string, len(RoomList))
	i := 0
	for k := range RoomList {
		lis[i] = k
		i++
	}

}
