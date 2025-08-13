package logger

type LogLevel string

// Уровни логов
const (
	LevelDebug   LogLevel = "debug"
	LevelInfo    LogLevel = "info"
	LevelWarning LogLevel = "warn"
	LevelError   LogLevel = "error"
	LevelFatal   LogLevel = "fatal"
)

var mapPriority = map[LogLevel]int{
	LevelDebug:   1,
	LevelInfo:    2,
	LevelWarning: 3,
	LevelError:   4,
	LevelFatal:   5,
}

func (l LogLevel) GreaterThan(other LogLevel) bool {
	return mapPriority[l] > mapPriority[other]
}

func (l LogLevel) ToUpper() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO "
	case LevelWarning:
		return "WARN "
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarning:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		return ""
	}
}
