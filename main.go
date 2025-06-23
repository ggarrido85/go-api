package main

import (
	"flag"
	"log"
	"log/slog"

	"github.com/ggarrido85/api-backend/cmd"
	"github.com/ggarrido85/api-backend/utils"
)

// Static variable set at compilation-time through linker flags.
//
//	$ go build -ldflags '-X main.apiVersion=v0.10.0 -X ...' .
var (
	apiVersion      string = "dev"
	segmentWriteKey string = ""
)

var compiledConfig = cmd.CompiledConfig{
	Version:         apiVersion,
	SegmentWriteKey: segmentWriteKey,
}

func main() {
	shouldRunServer := flag.Bool("server", false, "Run server")

	flag.Parse()
	logger := utils.NewLogger("text")
	logger.Info("Flags",
		slog.Bool("shouldRunServer", *shouldRunServer),
	)

	if *shouldRunServer {
		err := cmd.RunServer(compiledConfig)
		if err != nil {
			log.Fatal(err)
		}
	}
}
