package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:           "v2node",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Run() error {
	err := command.Execute()
	if err != nil {
		log.WithField("err", err).Error("Execute command failed")
	}
	return err
}
