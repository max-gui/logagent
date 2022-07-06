package logsets

import (
	"flag"
)

var (
	Appname, Appenv, Appdc, Apppath, Port *string
	Jsonlog                               *bool
)

func init() {
	// Appname = "charon"
	Apppath = flag.String("apppath", "/Users/max/Downloads/regagent", "app root path")
	Appname = flag.String("appname", "noname", "application's name")
	Appenv = flag.String("appenv", "prod", "this application's working env")
	Appdc = flag.String("appdc", "LFB", "this application's working dc")
	Port = flag.String("port", "7979", "this app's port")
	Jsonlog = flag.Bool("jsonlog", false, "jsonlog or not")

}

// var Reppath = func() string {
// 	return *Apppath + PthSep + *Repopathname + PthSep
// }
