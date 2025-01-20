package main

import "fmt"

type Commander interface {
	Run([]string) error
	ParseArgs(args []string) map[string]string
}

func ParseArgs(args []string) {
	var subC string
	var rest []string
	if len(args) < 1 {
		subC = "--help"
	} else {
		subC = args[0]
		rest = args[1:]
	}

	switch subC {
	case "-h", "--h", "-help", "--help":
		logger.Info("call help")
	case "p", "pulse":
		msg := fmt.Sprintf("remaining args: %v", rest)
		logger.Info(msg)
		Run(args)
	case "i", "import":
		logger.Info("import command")
	case "s", "start", "serve", "server":
		logger.Info("server command")
	default:
		logger.Info("no match, call help")
	}
}
