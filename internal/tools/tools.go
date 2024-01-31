package tools

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"

	"github.com/Ruslan-Androsenko/system-monitoring/api/proto"
)

// Выполнить консольную команду без фильтрации данных.
func execute(name string, args []string) ([]byte, error) {
	cmd := exec.Command(name, args...)

	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to pipe on cmd, error: %w", err)
	}
	defer pipe.Close()

	if err = cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start cmd, error: %w", err)
	}

	buffer := make([]byte, bufferSize)
	bytes, err := pipe.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read pipe out, error: %w", err)
	}

	if err = cmd.Wait(); err != nil {
		return nil, fmt.Errorf("failed to wait cmd, error: %w", err)
	}

	return buffer[:bytes], nil
}

// Выполнить консольную команду с фильтрацией данных.
func executeWithPipe(mainName, secondName string, mainArgs, secondArgs []string) ([]byte, error) {
	cmdMain := exec.Command(mainName, mainArgs...)
	cmdSecond := exec.Command(secondName, secondArgs...)

	pipeMain, err := cmdMain.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to pipe on main cmd, error: %w", err)
	}
	defer pipeMain.Close()

	cmdSecond.Stdin = pipeMain
	pipeSecond, err := cmdSecond.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to pipe on second cmd, error: %w", err)
	}
	defer pipeSecond.Close()

	if err = cmdMain.Start(); err != nil {
		return nil, fmt.Errorf("failed to start main cmd, error: %w", err)
	}

	if err = cmdSecond.Start(); err != nil {
		return nil, fmt.Errorf("failed to start second cmd, error: %w", err)
	}

	buffer := make([]byte, bufferSize)
	bytes, err := pipeSecond.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read second pipe out, error: %w", err)
	}

	if err = cmdMain.Wait(); err != nil {
		return nil, fmt.Errorf("failed to wait main cmd, error: %w", err)
	}

	if err = cmdSecond.Wait(); err != nil {
		return nil, fmt.Errorf("failed to wait second cmd, error: %w", err)
	}

	return buffer[:bytes], nil
}

// Парсинг трёх числовых значений по указанному формату.
func parsingThreeNumbers(buffer []byte, format string) (float64, float64, float64, error) {
	var first, second, third float64

	pattern := fmt.Sprintf(format, percentWithComaPattern, percentWithComaPattern, percentWithComaPattern)
	re := regexp.MustCompile(pattern)

	// Проверяем в каком формате необходимо парсить значения
	if !re.Match(buffer) {
		pattern = fmt.Sprintf(format, percentWithPointPattern, percentWithPointPattern, percentWithPointPattern)
		re = regexp.MustCompile(pattern)
	}

	numbers := re.ReplaceAllString(string(buffer), "$1.$2, $3.$4, $5.$6")
	_, err := fmt.Sscanf(numbers, "%f, %f, %f", &first, &second, &third)

	return first, second, third, err
}

// GetLoadAverage Получить значение для расчета средней нагрузки системы.
func GetLoadAverage(ctx context.Context, resCh chan<- float64, errCh chan<- error) {
	defer close(resCh)

	for {
		select {
		case <-ctx.Done():
			return

		default:
			buffer, err := executeWithPipe(topCmd, grepCmd, topArgs, loadAverageArgs)
			if err != nil {
				errCh <- err
				return
			}

			var oneMinute float64

			pattern := fmt.Sprintf(loadAveragePatternFormat, percentWithComaPattern)
			re := regexp.MustCompile(pattern)

			if !re.Match(buffer) {
				pattern = fmt.Sprintf(loadAveragePatternFormat, percentWithPointPattern)
				re = regexp.MustCompile(pattern)
			}

			number := re.ReplaceAllString(string(buffer), "$1.$2")
			_, err = fmt.Sscanf(number, "%f", &oneMinute)
			if err != nil {
				errCh <- fmt.Errorf("failed to parse numbers for load average, error: %w", err)
				return
			}

			resCh <- oneMinute
		}
	}
}

// GetCPULoad Получить значения для расчета средней нагрузки процессора.
func GetCPULoad(ctx context.Context, resCh chan<- *proto.CpuLoad, errCh chan<- error) {
	defer close(resCh)

	for {
		select {
		case <-ctx.Done():
			return

		default:
			buffer, err := executeWithPipe(topCmd, grepCmd, topArgs, cpuLoadArgs)
			if err != nil {
				errCh <- err
				return
			}

			userMode, systemMode, idle, err := parsingThreeNumbers(buffer, cpuLoadPatternFormat)
			if err != nil {
				errCh <- fmt.Errorf("failed to parse numbers for cpu load, error: %w", err)
				return
			}

			resCh <- &proto.CpuLoad{
				UserMode:   userMode,
				SystemMode: systemMode,
				Idle:       idle,
			}
		}
	}
}

// GetDiskLoad Получить значения для расчета средней нагрузки диска.
func GetDiskLoad(ctx context.Context, resCh chan<- *proto.DiskLoad, errCh chan<- error) {
	defer close(resCh)

	for {
		select {
		case <-ctx.Done():
			return

		default:
			buffer, err := execute(iostatCmd, iostatArgs)
			if err != nil {
				errCh <- err
				return
			}

			re := regexp.MustCompile("\n")
			items := re.Split(string(buffer), -1)

			var transferPerSecond, readPerSecond, writePerSecond float64

			if len(items) > 2 {
				pattern := fmt.Sprintf(diskLoadPatternFormat, percentWithPointPattern,
					percentIntegerPattern, percentWithPointPattern)
				re = regexp.MustCompile(pattern)

				numbers := re.ReplaceAllString(items[2], "$1.$2, $3, $4.$5")
				_, err = fmt.Sscanf(numbers, "%f, %f, %f", &writePerSecond, &transferPerSecond, &readPerSecond)
				if err != nil {
					errCh <- fmt.Errorf("failed to parse numbers for disk load, error: %w", err)
					return
				}
			}

			resCh <- &proto.DiskLoad{
				TransferPerSecond: transferPerSecond,
				ReadPerSecond:     readPerSecond,
				WritePerSecond:    writePerSecond,
			}
		}
	}
}

// GetDiskInfo Получить информацию об файловых системах диска.
func GetDiskInfo(ctx context.Context, resCh chan<- map[string]*proto.DiskInfo, errCh chan<- error) {
	defer close(resCh)

	for {
		select {
		case <-ctx.Done():
			return

		default:
			diskUsedSize, err := getUsedPercents(diskSizeArgs)
			if err != nil {
				errCh <- err
				return
			}

			diskUsedInode, err := getUsedPercents(diskInodesArgs)
			if err != nil {
				errCh <- err
				return
			}

			res := make(map[string]*proto.DiskInfo)

			for fileSystem, usedSize := range diskUsedSize {
				res[fileSystem] = &proto.DiskInfo{
					UsageSize:  usedSize,
					UsageInode: diskUsedInode[fileSystem],
				}
			}

			resCh <- res
		}
	}
}

// Получить процентное значение использования файловой системы диска.
func getUsedPercents(diskArgs []string) (map[string]float64, error) {
	res := make(map[string]float64)
	buffer, err := execute(diskFreeCmd, diskArgs)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile("\n")
	items := re.Split(string(buffer), -1)

	var (
		fileSystem        string
		availableQuantity float64
	)

	re = regexp.MustCompile(diskFreePattern)

	for i := 1; i < len(items); i++ {
		item := strings.ReplaceAll(items[i], "-", "0%")
		line := re.ReplaceAllString(item, "$1 $2")

		_, err = fmt.Sscanf(line, "%s %f", &fileSystem, &availableQuantity)
		if errors.Is(err, io.EOF) {
			continue
		} else if err != nil {
			return nil, err
		}

		res[fileSystem] += availableQuantity
	}

	return res, nil
}

// GetNetworkStats Получить информацию по сетевым соединениям.
func GetNetworkStats(ctx context.Context, resCh chan<- *proto.NetworkStats, errCh chan<- error) {
	defer close(resCh)

	var sudoPassword string
	sudoPasswordValue := ctx.Value(SudoPassCtxKey)
	if value, ok := sudoPasswordValue.(string); ok {
		sudoPassword = value
	}

	for {
		select {
		case <-ctx.Done():
			return

		default:
			listenerSockets, err := getListenerSockets(sudoPassword)
			if err != nil {
				errCh <- err
				return
			}

			tcpConnections, err := getCounterConnections(ssTCPArgs)
			if err != nil {
				errCh <- err
				return
			}

			udpConnections, err := getCounterConnections(ssUDPArgs)
			if err != nil {
				errCh <- err
				return
			}

			resCh <- &proto.NetworkStats{
				ListenerSocket: listenerSockets,
				CounterConnections: &proto.CounterConnections{
					Tcp: tcpConnections,
					Udp: udpConnections,
				},
			}
		}
	}
}

// Получить список прослушиваемых сокетов.
func getListenerSockets(sudoPassword string) ([]*proto.ListenerSocket, error) {
	var (
		buffer []byte
		err    error
	)

	if len(sudoPassword) > 0 {
		buffer, err = executeWithPipe(echoCmd, sudoCmd, []string{sudoPassword}, netstatWithSudoArgs)
	} else {
		buffer, err = execute(netstatCmd, netstatArgs)
	}

	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile("\n")
	items := re.Split(string(buffer), -1)

	var (
		user, pid, port   uint32
		command, protocol string
		res               = make([]*proto.ListenerSocket, 0, len(items)-2)
	)

	re = regexp.MustCompile(netstatPattern)

	for i := 2; i < len(items); i++ {
		item := strings.ReplaceAll(items[i], " - ", " 0/hide ")
		line := re.ReplaceAllString(item, "$1 $2 $3 $4 $5")

		_, err = fmt.Sscanf(line, "%s %d %d %d %s", &protocol, &port, &user, &pid, &command)
		if errors.Is(err, io.EOF) {
			continue
		} else if err != nil {
			return nil, err
		}

		res = append(res, &proto.ListenerSocket{
			User:     user,
			Pid:      pid,
			Command:  command,
			Protocol: protocol,
			Port:     port,
		})
	}

	return res, nil
}

// Получить количество TCP|UDP соединений, находящихся в разных состояниях.
func getCounterConnections(args []string) (map[string]uint32, error) {
	res := make(map[string]uint32)
	buffer, err := execute(ssCmd, args)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile("\n")
	items := re.Split(string(buffer), -1)

	var state string
	re = regexp.MustCompile(ssPattern)

	for i := 1; i < len(items); i++ {
		line := re.ReplaceAllString(items[i], "$1")

		_, err = fmt.Sscanf(line, "%s", &state)
		if errors.Is(err, io.EOF) {
			continue
		} else if err != nil {
			return nil, err
		}

		res[state]++
	}

	return res, nil
}
