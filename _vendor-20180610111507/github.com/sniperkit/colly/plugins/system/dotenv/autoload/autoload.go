package autoload

/*
	You can just read the .env file on import just by doing
		import _ "github.com/sniperkit/colly/plugins/system/dotenv/autoload"
	And bob's your mother's brother
*/

import (
	"github.com/sniperkit/colly/plugins/system/dotenv"
)

func init() {
	dotenv.Load()
}
