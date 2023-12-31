package server

import (
	"math"
	"sync"

	"github.com/Ruslan-Androsenko/system-monitoring/api/proto"
)

// Заполнение массива данных полученными значениями из различных каналов.
func fillDataSlice(dataItem *proto.MonitoringResponse, metricsChs MetricsChannels) *proto.MonitoringResponse {
	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)

	if metricsConf.LoadAverage {
		go func() {
			defer mu.Unlock()
			defer wg.Done()

			wg.Add(1)
			mu.Lock()
			dataItem.LoadAverage = <-metricsChs.loadAverageCh
		}()
	}

	if metricsConf.CPULoad {
		go func() {
			defer mu.Unlock()
			defer wg.Done()

			wg.Add(1)
			mu.Lock()
			dataItem.CpuLoad = <-metricsChs.cpuLoadCh
		}()
	}

	if metricsConf.DiskLoad {
		go func() {
			defer mu.Unlock()
			defer wg.Done()

			wg.Add(1)
			mu.Lock()
			dataItem.DiskLoad = <-metricsChs.diskLoadCh
		}()
	}

	if metricsConf.DiskInfo {
		go func() {
			defer mu.Unlock()
			defer wg.Done()

			wg.Add(1)
			mu.Lock()
			dataItem.DiskInfo = <-metricsChs.diskInfoCh
		}()
	}

	wg.Wait()

	return dataItem
}

// Сформировать массив данных из необходимого диапазона, для дальнейших расчетов усредненных значений.
func makeDataSlice(data []*proto.MonitoringResponse, currentIndex, avgSeconds int) []*proto.MonitoringResponse {
	var (
		dataSlice  []*proto.MonitoringResponse
		startIndex int
	)

	// Первый проход во время заполнения массива
	if currentIndex >= avgSeconds {
		if currentIndex != avgSeconds {
			startIndex = currentIndex - avgSeconds
		}

		dataSlice = data[startIndex:currentIndex]
	} else {
		// Когда массив уже заполнен и идет перезапись по второму кругу на новой минуте
		startIndex = countSeconds - avgSeconds + currentIndex
		dataSlice = data[startIndex:countSeconds]

		// Проверяем сколько элементов не хватает для полного снепшота, и дополняем их с начала массива
		incomplete := avgSeconds - len(dataSlice)

		if incomplete > 0 {
			dataSlice = append(dataSlice, data[0:currentIndex]...)
		}
	}

	return dataSlice
}

// Расчет усредненных значений перед отправкой данных клиенту.
func calculateAverageOfSlice(data []*proto.MonitoringResponse) *proto.MonitoringResponse {
	var (
		length      = float64(len(data))
		loadAverage float64
		cpuLoad     = &proto.CpuLoad{}
		diskLoad    = &proto.DiskLoad{}
		diskInfo    = make(map[string]*proto.DiskInfo)
	)

	for _, item := range data {
		if metricsConf.LoadAverage {
			loadAverage += item.LoadAverage
		}

		cpuLoad = sumCPULoad(cpuLoad, item.CpuLoad)
		diskLoad = sumDiskLoad(diskLoad, item.DiskLoad)
		diskInfo = sumDiskInfo(diskInfo, item.DiskInfo)
	}

	return &proto.MonitoringResponse{
		LoadAverage: roundNumber(loadAverage / length),
		CpuLoad:     makeAverageCPULoad(cpuLoad, length),
		DiskLoad:    makeAverageDiskLoad(diskLoad, length),
		DiskInfo:    makeAverageDiskInfo(diskInfo, length),
	}
}

// Округляем число до 2-х знаков после запятой.
func roundNumber(number float64) float64 {
	return math.Round(number*100.0) / 100.0
}

// Суммирование данных по загрузке процессора.
func sumCPULoad(cpuLoad, itemCPULoad *proto.CpuLoad) *proto.CpuLoad {
	if !metricsConf.CPULoad {
		return nil
	}

	cpuLoad.UserMode += itemCPULoad.UserMode
	cpuLoad.SystemMode += itemCPULoad.SystemMode
	cpuLoad.Idle += itemCPULoad.Idle

	return cpuLoad
}

// Расчет усредненных данных по загрузке процессора.
func makeAverageCPULoad(cpuLoad *proto.CpuLoad, length float64) *proto.CpuLoad {
	if !metricsConf.CPULoad {
		return nil
	}

	return &proto.CpuLoad{
		UserMode:   roundNumber(cpuLoad.UserMode / length),
		SystemMode: roundNumber(cpuLoad.SystemMode / length),
		Idle:       roundNumber(cpuLoad.Idle / length),
	}
}

// Суммирование данных по загрузке диска.
func sumDiskLoad(diskLoad, itemDiskLoad *proto.DiskLoad) *proto.DiskLoad {
	if !metricsConf.DiskLoad {
		return nil
	}

	diskLoad.TransferPerSecond += itemDiskLoad.TransferPerSecond
	diskLoad.ReadPerSecond += itemDiskLoad.ReadPerSecond
	diskLoad.WritePerSecond += itemDiskLoad.WritePerSecond

	return diskLoad
}

// Расчет усредненных данных по загрузке диска.
func makeAverageDiskLoad(diskLoad *proto.DiskLoad, length float64) *proto.DiskLoad {
	if !metricsConf.DiskLoad {
		return nil
	}

	return &proto.DiskLoad{
		TransferPerSecond: roundNumber(diskLoad.TransferPerSecond / length),
		ReadPerSecond:     roundNumber(diskLoad.ReadPerSecond / length),
		WritePerSecond:    roundNumber(diskLoad.WritePerSecond / length),
	}
}

// Суммирование данных использования диска.
func sumDiskInfo(diskInfo, itemDiskInfo map[string]*proto.DiskInfo) map[string]*proto.DiskInfo {
	if !metricsConf.DiskInfo {
		return nil
	}

	for fileSystem, item := range itemDiskInfo {
		if _, ok := diskInfo[fileSystem]; ok {
			diskInfo[fileSystem].UsageSize += item.UsageSize
			diskInfo[fileSystem].UsageInode += item.UsageInode
		} else {
			diskInfo[fileSystem] = &proto.DiskInfo{
				UsageSize:  item.UsageSize,
				UsageInode: item.UsageInode,
			}
		}
	}

	return diskInfo
}

// Расчет усредненных данных использования диска.
func makeAverageDiskInfo(diskInfo map[string]*proto.DiskInfo, length float64) map[string]*proto.DiskInfo {
	if !metricsConf.DiskInfo {
		return nil
	}

	for fileSystem, item := range diskInfo {
		diskInfo[fileSystem].UsageSize = roundNumber(item.UsageSize / length)
		diskInfo[fileSystem].UsageInode = roundNumber(item.UsageInode / length)
	}

	return diskInfo
}
