package errors_util

import (
	"log"
	"os"

	"github.com/go-errors/errors"
)

func InitLogFile(logFileName string) {
	fLog, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	CrashOnError(err)
	defer saveLog(fLog)
}

func CrashOnError(err error) {
	if err != nil {
		log.Printf("ERROR CRASH: [%v]", err)
		panic(err)
	}
}

func LogError(err error) {
	if err != nil {
		log.Printf("ERROR CRASH: [%v]", err)
		println(err)
	}
}

func CrashWithError(err *errors.Error) {
	if err != nil {
		log.Printf("ERROR CRASH: [%v]", err)
		panic(err)
	}
}

func PrintOnError(err *errors.Error) {
	if err != nil {
		log.Printf("ERROR: [%v]", err)
		println(err)
	}
}

func RecoverOnError(err *errors.Error) {
	if err != nil {
		log.Printf("ERROR RECOVER: [%v]", err)
		recover()
	}
}

func saveLog(fLog *os.File) {
	log.SetOutput(fLog)
	fLog.Close()
}
