package logsets

import (
	"flag"
)

var (
	Appname, Appenv, Apppath *string
	Jsonlog                  *bool
)

func init() {
	// Appname = "charon"
	Appname = flag.String("appname", "noname", "application's name")
	Appenv = flag.String("appenv", "prod", "this application's working env")
	Jsonlog = flag.Bool("jsonlog", false, "jsonlog or not")

}

// var Reppath = func() string {
// 	return *Apppath + PthSep + *Repopathname + PthSep
// }
