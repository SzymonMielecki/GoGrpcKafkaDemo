package logs

import (
	"fmt"
	"os"
	"path/filepath"
)

func WriteToLogs(message string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	logsFile := filepath.Join(homeDir, ".chatapp_logs")
	data := []byte(message + "\n")
	file, err := os.OpenFile(logsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error writing to logs file:", err)
		return
	}
	fmt.Println("Message written to logs file:", message)
}
