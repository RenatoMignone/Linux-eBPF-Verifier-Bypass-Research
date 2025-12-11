package logger

import (
	"fmt"
	"os"
	"path/filepath"
)

var verboseMode bool
var logDir string

func Init(verbose bool, exportPath string) {
	verboseMode = verbose
	logDir = exportPath
	os.MkdirAll(logDir, 0755)
}

func LogError(context, msg string) {
	if verboseMode {
		fmt.Printf("[ERROR] %s: %s\n", context, msg)
	}
	SaveLog(context, "[ERROR] "+msg)
}

func SaveLog(filePath, content string) {
	base := filepath.Base(filePath)
	logFile := filepath.Join(logDir, base+".log")
	_ = os.WriteFile(logFile, []byte(content), 0644)
}

