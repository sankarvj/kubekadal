package main

import (
	"github.com/codegangsta/negroni"
	"github.com/sankarvj/kubekadal/app/router"
)

func main() {
	n := negroni.Classic()
	n.UseHandler(router.InitRouter())
	n.Run(":3011")
}
