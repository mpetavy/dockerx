package main

import (
	"embed"
	"flag"
	"github.com/docker/docker/client"
	"github.com/mpetavy/common"
	"strings"
)

var (
	container = flag.Bool("container", false, "container")
	image     = flag.Bool("image", false, "image")
	execute   = flag.String("execute", "", "execute")

	list  = flag.Bool("list", false, "list")
	kill  = flag.Bool("kill", false, "kill")
	query = flag.String("q", "", "query")

	filter = flag.String("f", "", "filter")
)

//go:embed go.mod
var resources embed.FS

func init() {
	common.Init("", "", "", "", "Extended Docker interaction", "", "", "", &resources, nil, nil, run, 0)
}

func TrimApostroph(str string) string {
	s := 0
	e := len(str)
	if strings.HasPrefix(str, "\"") {
		s++
	}
	if strings.HasSuffix(str, "\"") {
		e--
	}

	return str[s:e]
}

func run() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if common.Error(err) {
		return err
	}

	switch {
	case *container:
		err = containers(cli)
		if common.Error(err) {
			return err
		}
	case *image:
		err = images(cli)
		if common.Error(err) {
			return err
		}
	case *execute != "":
		err = executeIt(cli)
		if common.Error(err) {
			return err
		}
	}

	return nil
}

func main() {
	common.Run([]string{"image|container"})
}
