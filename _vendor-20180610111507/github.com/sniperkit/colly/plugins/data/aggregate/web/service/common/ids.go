package common

import (
	"github.com/sniperkit/colly/plugins/data/aggregate/web/service/resource"
)

func AdminChannelID(admin *resource.Admin) string {
	return GetMd5(admin.Id + admin.Org.Id)
}
