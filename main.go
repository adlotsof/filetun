package main


import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/adlotsof/filetun/config"
	"github.com/adlotsof/filetun/cli"
)


func main() {
	ctx := kong.Parse(&config.CLI)
	fmt.Println(ctx.Args)
	cli.Run()

}
