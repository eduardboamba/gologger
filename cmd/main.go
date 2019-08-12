package main

import "github.com/eduardboamba/gologger/pkg/util/logger"

func main() {
	logger.Info("this app is running great")

	logger.Debug("got some debugging stuff logged here")

	logger.Error("something might have gone wrong")

	logger.Fatal("oops, fatality...")
}
