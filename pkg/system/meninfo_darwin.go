package system

import (
	"fmt"
)

/*
#include <stdlib.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/sysctl.h>

struct xsw_usage* getswapusage(){
	struct xsw_usage* vmusage =  NULL;
	size_t size = sizeof(struct xsw_usage);
	int iret = 0;

	vmusage = (struct xsw_usage*)malloc(size);
	iret = sysctlbyname("vm.swapusage", vmusage, &size, NULL, 0);

	if(iret != 0){
		free(vmusage);
		return NULL;
	}
	return vmusage;
}
void freeswapusage(struct xsw_usage* usage){
	if (usage == NULL) return;
	free(usage);
}
size_t getfreememory(){
	size_t pagesize = 0;
	size_t freepages = 0;
	size_t speculativepages = 0;
	size_t size = sizeof(pagesize);
	int iret = -1;

	iret = sysctlbyname("vm.pagesize",&pagesize,&size,NULL,0);
	if(iret != 0){
		return -1;
	}

	iret = sysctlbyname("vm.page_free_count",&freepages,&size,NULL,0);
	if(iret != 0){
		return -1;
	}

	iret = sysctlbyname("vm.page_speculative_count",&speculativepages,&size,NULL,0);
	if(iret != 0){
		return -1;
	}

	return (freepages + speculativepages) * pagesize;
}
*/
import "C"

// Get the system memory info using sysconf same as prtconf
func getTotalMem() int64 {
	pagesize := C.sysconf(C._SC_PAGESIZE)
	npages := C.sysconf(C._SC_PHYS_PAGES)
	return int64(pagesize * npages)
}

func getFreeMem() int64 {
	return int64(C.getfreememory())
}

// ReadMemInfo retrieves memory statistics of the host system and returns a
//  MemInfo type.
func ReadMemInfo() (*MemInfo, error) {
	MemTotal := getTotalMem()
	MemFree := getFreeMem()
	SwapTotal, SwapFree, err := getSysSwap()

	if MemTotal < 0 || MemFree < 0 || SwapTotal < 0 || SwapFree < 0 {
		return nil, fmt.Errorf("Error getting system memory info %v\n", err)
	}

	meminfo := &MemInfo{}
	// Total memory is total physical memory less than memory locked by kernel
	meminfo.MemTotal = MemTotal //- int64(ppKernel)
	meminfo.MemFree = MemFree
	meminfo.SwapTotal = SwapTotal
	meminfo.SwapFree = SwapFree

	return meminfo, nil
}

func getSysSwap() (int64, int64, error) {
	var tSwap int64
	var fSwap int64
	//	var pagesize int64

	usage := C.getswapusage()
	if usage == nil {
		return 0, 0, fmt.Errorf(`get swap file usage fail`)
	}
	//	pagesize = int64(usage.xsu_pagesize)
	tSwap = int64(usage.xsu_total)
	fSwap = int64(usage.xsu_avail)

	C.freeswapusage(usage)
	return tSwap, fSwap, nil
}
