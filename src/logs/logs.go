package logs

type Logger struct {
	logs []string
}

func (l *Logger) AddNewEntry(log string) {
	l.logs = append(l.logs, log)
}

func (l *Logger) GetLastLog() string {
	return l.logs[len(l.logs)-1]
}
