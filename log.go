package main

// This trick is learnt from a post by Rob Pike
// https://groups.google.com/d/msg/golang-nuts/gU7oQGoCkmg/j3nNxuS2O_sJ

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type infoLogging bool
type debugLogging bool
type errorLogging bool
type requestLogging bool
type responseLogging bool

var (
	info   infoLogging
	debug  debugLogging
	errl   errorLogging
	dbgRq  requestLogging
	dbgRep responseLogging

	logFile io.Writer
	logBuf  *bufio.Writer // only set if output is not stdio

	debugLog, errorLog, requestLog, responseLog *log.Logger
)

var (
	verbose  bool
	colorize bool
)

func init() {
	flag.BoolVar((*bool)(&info), "info", true, "info log")
	flag.BoolVar((*bool)(&debug), "debug", false, "debug log")
	flag.BoolVar((*bool)(&errl), "err", true, "error log")
	flag.BoolVar((*bool)(&dbgRq), "request", false, "request log")
	flag.BoolVar((*bool)(&dbgRep), "reply", false, "reply log")
	flag.BoolVar(&verbose, "v", false, "More info in request/response logging")
	flag.BoolVar(&colorize, "color", false, "Colorize log output")
}

func initLog() {
	logFile = os.Stdout
	if config.logFile != "" {
		if config.logFile[0] == '~' {
			config.logFile = homeDir + config.logFile[1:]
		}

		if f, err := os.OpenFile(config.logFile,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600); err != nil {
			fmt.Printf("Can't open log file, logging to stdout: %v\n", err)
		} else {
			logBuf = bufio.NewWriter(f)
			logFile = logBuf
		}
	}
	log.SetOutput(logFile)
	if colorize {
		errorLog = log.New(logFile, "\033[31m[Error]\033[0m ", log.LstdFlags)
		debugLog = log.New(logFile, "\033[34m[Debug]\033[0m ", log.LstdFlags)
		requestLog = log.New(logFile, "\033[32m[>>>>>]\033[0m ", log.LstdFlags)
		responseLog = log.New(logFile, "\033[33m[<<<<<]\033[0m ", log.LstdFlags)
	} else {
		errorLog = log.New(logFile, "[ERROR] ", log.LstdFlags)
		debugLog = log.New(logFile, "[DEBUG] ", log.LstdFlags)
		requestLog = log.New(logFile, "[>>>>>] ", log.LstdFlags)
		responseLog = log.New(logFile, "[<<<<<] ", log.LstdFlags)
	}
}

func (d infoLogging) Printf(format string, args ...interface{}) {
	if d {
		log.Printf(format, args...)
	}
}

func (d infoLogging) Println(args ...interface{}) {
	if d {
		log.Println(args...)
	}
}

func (d debugLogging) Printf(format string, args ...interface{}) {
	if d {
		debugLog.Printf(format, args...)
	}
}

func (d debugLogging) Println(args ...interface{}) {
	if d {
		debugLog.Println(args...)
	}
}

func (d errorLogging) Printf(format string, args ...interface{}) {
	if d {
		errorLog.Printf(format, args...)
	}
}

func (d errorLogging) Println(args ...interface{}) {
	if d {
		errorLog.Println(args...)
	}
}

func (d requestLogging) Printf(format string, args ...interface{}) {
	if d {
		requestLog.Printf(format, args...)
	}
}

func (d responseLogging) Printf(format string, args ...interface{}) {
	if d {
		responseLog.Printf(format, args...)
	}
}
