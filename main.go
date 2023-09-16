package main

import (
	"opslink/cmd"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.WithField("cmd", "main").Info("[OPSLINK-CLI]")
	cmd.Execute()
}
