// SPDX-License-Identifier: BSD-3-Clause

package util

import "git.froth.zone/sam/awl/logawl"

// InitLogger initializes the logawl instance.
func InitLogger(verbosity int) (log *logawl.Logger) {
	log = logawl.New()

	log.SetLevel(logawl.Level(verbosity))

	return log
}
