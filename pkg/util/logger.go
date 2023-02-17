// SPDX-License-Identifier: BSD-3-Clause

package util

import "dns.froth.zone/awl/pkg/logawl"

// InitLogger initializes the logawl instance.
func InitLogger(verbosity int) (log *logawl.Logger) {
	log = logawl.New()

	log.SetLevel(logawl.Level(verbosity))

	return log
}
