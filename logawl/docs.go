/*
Package logawl is a package for custom logging needs

LogAwl extends the standard log library with support for log levels
This is _different_ from the syslog package in the standard library because you do not define a file
because awl is a cli utility it writes directly to std err.
*/
//	Use the New() function to init logawl
//
//	logger := logawl.New()
//
// You can call specific logging levels from your new logger using
//
//	logger.Debug("Message to log")
//	logger.Fatal("Message to log")
//	logger.Info("Message to log")
//	logger.Error("Message to log")
//
// You may also set the log level on the fly with
//
//	Logger.SetLevel(3)
// This allows you to change the default level (Info)
// and prevent log messages from being posted at higher verbosity levels
// for example if
//	Logger.SetLevel(3)
// is not called and you call
//	Logger.Debug()
// this runs through
//	IsLevel(level)
// to verify if the debug log should be sent to std.Err or not based on the current expected log level
package logawl
