// SPDX-License-Identifier: BSD-3-Clause

package util

import "git.froth.zone/sam/awl/logawl"

// Initialize the logawl instance.
func InitLogger(verbosity int) (Logger *logawl.Logger) {
	Logger = logawl.New()

	Logger.SetLevel(logawl.Level(verbosity))

	return
}
