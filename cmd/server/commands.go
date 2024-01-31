package main

import "flag"

const CmdVersion = "version"

// Проверяем запрошено ли отображение данных о версии сервиса.
func hasVersionCommand() bool {
	return flag.Arg(0) == CmdVersion
}
