package cmd

import (
	log "github.com/sirupsen/logrus"

	_ "github.com/InazumaV/V2bX/core/hy2"
	_ "github.com/InazumaV/V2bX/core/imports"
	_ "github.com/InazumaV/V2bX/core/sing"
	_ "github.com/InazumaV/V2bX/core/xray"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use: "V2bX",
}

func Run() {
	err := command.Execute()
	if err != nil {
		log.WithField("err", err).Error("Execute command failed")
	}
}
