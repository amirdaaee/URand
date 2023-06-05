package main

import (
	config "URand/Config"
	supervisor "URand/Supervisor"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func initNS() map[string]*supervisor.NameSpace {
	ns := make(map[string]*supervisor.NameSpace)
	for _, c := range config.Config().NameSpaces {
		_ns, err := supervisor.NewNS(c.Name, c.Length)
		if err != nil {
			logrus.WithError(err).Panic("error initiating namespace")
		}
		ns[c.Name] = _ns
	}
	return ns
}
func main() {
	ns := initNS()
	r := gin.Default()
	for k, v := range ns {
		vv := *v
		r.GET("/"+k, func(c *gin.Context) {
			rnd, err := vv.Get()
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			c.String(http.StatusOK, rnd)
		})
	}
	r.Run(config.Config().Listen)
}
