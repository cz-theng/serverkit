package main

import (
	"fmt"

	"github.com/cz-it/golangutils/log"
)

func main() {
	fmt.Println("Testing log...")

	// Default log 
	// log to console
	log.DEBUG("DEBUG")
	log.INFO("INFO")
	log.WARNING("WARNING")
	log.ERROR("ERROR")
	log.FATAL("FATAL")

	//log to file
	// log.New(logPath, logName string)
	logger, err:= log.NewFileLogger("./test/log", "test")
	if err != nil {
		println("New file logger error")
		fmt.Println(err)
		return 
	}
	logger.SetMaxFileSize(4*1024) // 4M default 500M
	logger.SetLevel(log.LDEBUG)


	logger.Debug("Debug")
	logger.Info("Info")
	logger.Warning("Warning")
	logger.Warning("Warning")
	logger.Error("Error")
	logger.Fatal("Fatal")
	logger.Fatal("Fatal")
}














