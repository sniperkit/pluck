package gin

import (
	"github.com/devopsfaith/krakend/config"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"

	"github.com/devopsfaith/krakend-httpsecure"
)

// Register registers the secure middleware into the gin engine
func Register(cfg config.ExtraConfig, engine *gin.Engine) error {
	engine.Use(NewSecureMw(cfg))
	return nil
}

// NewSecureMw creates a secured middleware for the gin engine
func NewSecureMw(cfg config.ExtraConfig) gin.HandlerFunc {
	opt, ok := httpsecure.ConfigGetter(cfg).(secure.Options)
	if !ok {
		return func(c *gin.Context) {}
	}

	secureMiddleware := secure.New(opt)

	return func(c *gin.Context) {
		err := secureMiddleware.Process(c.Writer, c.Request)

		if err != nil {
			c.Abort()
			return
		}

		if status := c.Writer.Status(); status > 300 && status < 399 {
			c.Abort()
		}
	}
}
