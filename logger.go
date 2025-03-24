package capybara

import "log"

func InitLogger() *CapybaraLogger {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	return &CapybaraLogger{serviceName: "server"}
}

type CapybaraLogger struct {
	serviceName string
}

func (l *CapybaraLogger) Info(msg string) {
	log.Printf("[%s] INFO - %s", l.serviceName, msg)
}
