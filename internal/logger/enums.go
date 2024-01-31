package logger

// Числовые значения уровней лога.
const (
	fatalLevel logLevel = iota + 1
	errorLevel
	warningLevel
	infoLevel
	debugLevel
)

// Строковые значения уровней лога.
const (
	fatalTitle   = "FATAL"
	errorTitle   = "ERROR"
	warningTitle = "WARN"
	infoTitle    = "INFO"
	debugTitle   = "DEBUG"
)

// Словарь для получения числового значения уровня лога из строкового.
var titleToLevel = map[string]logLevel{
	fatalTitle:   fatalLevel,
	errorTitle:   errorLevel,
	warningTitle: warningLevel,
	infoTitle:    infoLevel,
	debugTitle:   debugLevel,
}

// Словарь для получения строкового значения уровня лога из числового.
var levelToTitle = map[logLevel]string{
	fatalLevel:   fatalTitle,
	errorLevel:   errorTitle,
	warningLevel: warningTitle,
	infoLevel:    infoTitle,
	debugLevel:   debugTitle,
}

// Получение строкового значения уровня лога.
func (level logLevel) String() string {
	return levelToTitle[level]
}
