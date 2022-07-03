package util

import "git.froth.zone/sam/awl/logawl"

func InitLogger(debug bool) (Logger *logawl.Logger) {
	Logger = logawl.New()

	if debug {
		Logger.SetLevel(3)
	}

	return
}
