// +build windows

package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/remote-tech/chirpstack-fuota-server/internal/config"
)

func setSyslog() error {
	if config.C.General.LogToSyslog {
		log.Fatal("syslog logging is not supported on Windows")
	}

	return nil
}
