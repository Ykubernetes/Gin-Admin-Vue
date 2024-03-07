package system

import (
	"fmt"
	"gitee.com/go-server/global"
	"gitee.com/go-server/service/response"
	"gitee.com/go-server/utils"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"strconv"
	"time"
)

var (
	cpuUsedPercent  float64
	memUsedPercent  float64
	diskUsed        uint64
	diskUsedPercent float64
)

func SystemState(c *gin.Context) {
	// cpu
	cpuPercents, _ := cpu.Percent(time.Second, true)
	for _, percent := range cpuPercents {
		cpuUsedPercent += percent
	}
	// 将cpuUsedPercent除以CPU核心的数量（即cpuPercents切片的长度），从而得到平均的CPU使用率。
	cpuUsedPercent /= float64(len(cpuPercents))

	// mem
	vms, _ := mem.VirtualMemory()
	memUsedPercent = vms.UsedPercent
	// disk
	partitions, _ := disk.Partitions(true)
	for _, partition := range partitions {
		us, _ := disk.Usage(partition.Mountpoint)
		diskUsed += us.Used
	}
	allUsage, _ := disk.Usage("/")
	diskUsedPercent = float64(diskUsed) / float64(allUsage.Total) * 100
	resData := make(map[string]string)
	resData["cpu_used_percent"] = strconv.FormatFloat(cpuUsedPercent, 'f', 0, 64)
	resData["mem_used_percent"] = strconv.FormatFloat(memUsedPercent, 'f', 0, 64)
	resData["disk_used_percent"] = strconv.FormatFloat(diskUsedPercent, 'f', 0, 64)
	response.ResSuccess(c, resData)
}

func GetSystemInfo(c *gin.Context) {
	server, err := GetServerInfo()
	if err != nil {
		global.Log.Errorf("获取信息信息失败!", err)
		response.ResFail(c, "获取信息信息失败!")
		return
	}
	resData := make(map[string]*utils.Server)
	resData["server"] = server
	response.ResSuccess(c, resData)

}

// 获取utils中system-info系统信息的数据
func GetServerInfo() (server *utils.Server, err error) {
	var s utils.Server
	s.Os = utils.InitOS()
	if s.Cpu, err = utils.InitCPU(); err != nil {
		fmt.Printf("func utils.InitCPU() Failed :", err.Error())
		return &s, err
	}
	if s.Ram, err = utils.InitRAM(); err != nil {
		fmt.Printf("func utils.InitRAM() Failed :", err.Error())
		return &s, err
	}
	if s.Disk, err = utils.InitDisk(); err != nil {
		fmt.Printf("func utils.InitDisk() Failed :", err.Error())
		return &s, err
	}

	return &s, nil
}
