package jk_utils

import (
	"fmt"
	"runtime"
	"syscall"
	"time"
)

func PrintMemUsage() {

	for  {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		// For info on each, see: https://golang.org/pkg/runtime/#MemStats
		fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))//分配的堆对象的字节
		//fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))//分配给堆对象的累积字节
		fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))//从OS获得的内存的总字节数
		//fmt.Printf("\tLookups = %v", m.Lookups)//执行的指针查找的数量
		//fmt.Printf("\tMallocs = %v", m.Mallocs)//分配的堆对象的累积计数
		//fmt.Printf("\tFrees = %v", m.Frees)//已释放的堆对象的累积计数
		fmt.Printf("\tLastGC = %v", m.LastGC)//最后一次垃圾回收完成的时间
		fmt.Printf("\tNextGC MiB = %v", bToMb(m.NextGC))//下一个GC周期的目标堆大小
		fmt.Printf("\tNumGC = %v", m.NumGC)//已完成的GC周期数


		var mem syscall.Rusage
		syscall.Getrusage(syscall.RUSAGE_SELF, &mem)
		fmt.Printf("\tmem.Maxrss %v MiB \n",bToMb(uint64(mem.Maxrss)))

		time.Sleep(time.Second)
	}
}




func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}