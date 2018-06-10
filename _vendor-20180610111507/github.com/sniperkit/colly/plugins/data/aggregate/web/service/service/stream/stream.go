package stream

import (
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/service/stream/controller"
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/service/stream/router"

	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/service/stream/controller/room"
)

type StreamService struct {
	feedRoomManager *room.FeedRoomManager
}

func (this *StreamService) GetFeedRoomManager() *room.FeedRoomManager {
	return this.feedRoomManager
}

func (this *StreamService) InitRooms() {
}

func (this *StreamService) Init() {
	// should pass controller from here
	// should not creates new one in InitRouters
	router.InitRouters()

	// should pass Feed Room for controllers to have access to CHANNELS
	controller.InitSession()

	this.GetFeedRoomManager().Run()
}

func NewStreamService() *StreamService {
	frm := room.NewFeedRoomManager()

	return &StreamService{frm}
}
