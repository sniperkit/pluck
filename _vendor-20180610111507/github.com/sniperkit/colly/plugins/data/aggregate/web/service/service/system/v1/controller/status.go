package controller

import (
	"github.com/feedlabs/feedify"
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/service/stream/controller/room"
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/service/stream/model"
	"os"
	"runtime"
	"strconv"
)

type StatusController struct {
	feedify.Controller
}

func (this *StatusController) Get() {
	var memstats runtime.MemStats

	runtime.ReadMemStats(&memstats)
	hostname, _ := os.Hostname()

	this.Data["json"] = map[string]interface{}{
		"system": map[string]interface{}{
			"hostname": hostname,
			"pid":      strconv.Itoa(os.Getpid()),
			"cpus":     strconv.Itoa(runtime.NumCPU()),
		},
		"go": map[string]interface{}{
			"routines": strconv.Itoa(runtime.NumGoroutine()),
		},
		"mem": map[string]interface{}{
			"mem_alloc":       strconv.Itoa(int(memstats.Alloc)),
			"mem_alloc_heap":  strconv.Itoa(int(memstats.HeapAlloc)),
			"mem_alloc_total": strconv.Itoa(int(memstats.TotalAlloc)),
			"mem_sys":         strconv.Itoa(int(memstats.Sys)),
		},
		"stream": map[string]interface{}{
			"subscribers":    strconv.Itoa(room.FeedRoom.Subscribers.Len()),
			"waitinglist":    strconv.Itoa(room.FeedRoom.WaitingList.Len()),
			"archived_queue": strconv.Itoa(model.Archive.Len()),
		},
	}

	this.Controller.ServeJson()
}
