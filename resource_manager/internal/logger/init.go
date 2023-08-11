package logger

import (
	"flag"
	"os"

	"k8s.io/klog"
)

func Init() {
	klog.InitFlags(nil)

	logToConsole := os.Getenv("LOG_TO_CONSOLE")
	flag.Set("logtostderr", logToConsole)

	flag.Set("log_file", os.Getenv("LOG_FILE_DIR"))

	// remove unnecessary header in logs
	flag.Set("skip_log_headers", "true")

	flag.Parse()
}
