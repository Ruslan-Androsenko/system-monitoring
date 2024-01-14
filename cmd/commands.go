package main

import "flag"

const (
	CmdVersion    = "version"
	CmdGrpcClient = "grpc-client"
)

// Проверяем запрошено ли отображение данных о версии сервиса.
func hasVersionCommand() bool {
	return flag.Arg(0) == CmdVersion
}

// Проверяем запрошен ли запуск Grpc клиента.
func hasGrpcClientCommand() bool {
	return flag.Arg(0) == CmdGrpcClient
}
