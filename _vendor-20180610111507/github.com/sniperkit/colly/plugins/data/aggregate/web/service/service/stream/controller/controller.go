package controller

import (
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/service/stream/controller/room"
)

func InitSession() {
	room.InitSessionManager()
}
