// SPDX-License-Identifier: BSD-3-Clause

package util

import "git.froth.zone/sam/awl/logawl"

// Initialize the logawl instance.
func InitLogger(verbosity int) (logger *logawl.Logger) {
	logger = logawl.New()

	logger.SetLevel(logawl.Level(verbosity))

	return logger
}
