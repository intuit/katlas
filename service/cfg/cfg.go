//Package cfg - contains all global configuration that should only be set once at the startup
package cfg

import "flag"

type (
	serverCfg struct {
		ServerType string
		DgraphHost string
	}
)

var (
	//ServerCfg ...All Config var related to Rest API Server
	ServerCfg serverCfg
)

func init() {

	flag.StringVar(&ServerCfg.ServerType, "serverType", "http", "Mode the Rest Service runs in - Secure/Insecure")
	flag.StringVar(&ServerCfg.DgraphHost, "dgraphHost", "127.0.0.1:9080", "Mode the Rest Service runs in - Secure/Insecure")
}
