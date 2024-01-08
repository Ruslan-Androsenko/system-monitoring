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

	buffer := make([]byte, 1024)
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

	buffer := make([]byte, 1024)
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

			pattern := fmt.Sprintf(".*average:%s,.*", percentPattern)
			re := regexp.MustCompile(pattern)
			numbers := re.ReplaceAllString(string(buffer), "$1.$2")

			_, err = fmt.Sscanf(numbers, "%f", &oneMinute)
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

			var userMode, systemMode, idle float64

			pattern := fmt.Sprintf("%sus,%ssy,.*ni,%sid.*", percentPattern, percentPattern, percentPattern)
			re := regexp.MustCompile(pattern)
			numbers := re.ReplaceAllString(string(buffer), "$1.$2, $3.$4, $5.$6")

			_, err = fmt.Sscanf(numbers, "%f, %f, %f", &userMode, &systemMode, &idle)
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
			buffer, err := executeWithPipe(iostatCmd, grepCmd, iostatArgs, diskLoadArgs)
			if err != nil {
				errCh <- err
				return
			}

			var transferPerSecond, readPerSecond, writePerSecond float64

			pattern := fmt.Sprintf(".*sda%s%s%s", percentPattern, percentPattern, percentPattern)
			re := regexp.MustCompile(pattern)
			numbers := re.ReplaceAllString(string(buffer), "$1.$2, $3.$4, $5.$6")

			_, err = fmt.Sscanf(numbers, "%f, %f, %f", &transferPerSecond, &readPerSecond, &writePerSecond)
			if err != nil {
				errCh <- fmt.Errorf("failed to parse numbers, error: %w", err)
				return
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
