// +build windows

package disk

import (
	"bytes"
	"unsafe"

	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/internal/common"
	"golang.org/x/sys/windows"
)

var (
	procGetDiskFreeSpaceExW     = common.Modkernel32.NewProc("GetDiskFreeSpaceExW")
	procGetLogicalDriveStringsW = common.Modkernel32.NewProc("GetLogicalDriveStringsW")
	procGetDriveType            = common.Modkernel32.NewProc("GetDriveTypeW")
	provGetVolumeInformation    = common.Modkernel32.NewProc("GetVolumeInformationW")
)

var (
	FileFileCompression = int64(16)     // 0x00000010
	FileReadOnlyVolume  = int64(524288) // 0x00080000
)

type Win32_PerfFormattedData struct {
	Name                    string
	AvgDiskBytesPerRead     uint64
	AvgDiskBytesPerWrite    uint64
	AvgDiskReadQueueLength  uint64
	AvgDiskWriteQueueLength uint64
	AvgDisksecPerRead       uint64
	AvgDisksecPerWrite      uint64
}

type Win32_PerfFormattedData_PerfDisk_PhysicalDisk struct {
	AvgDiskBytesPerRead     uint64
	AvgDiskBytesPerTransfer uint64
	AvgDiskBytesPerWrite    uint64
	AvgDiskQueueLength      uint64
	AvgDiskReadQueueLength  uint64
	AvgDiskSecPerRead       uint32
	AvgDiskSecPerTransfer   uint32
	AvgDiskSecPerWrite      uint32
	AvgDiskWriteQueueLength uint64
	CurrentDiskQueueLength  uint32
	DiskBytesPerSec         uint64
	DiskReadBytesPerSec     uint64
	DiskReadsPerSec         uint32
	DiskTransfersPerSec     uint32
	DiskWriteBytesPerSec    uint64
	DiskWritesPerSec        uint32
	Name                    string
	PercentDiskReadTime     uint64
	PercentDiskTime         uint64
	PercentDiskWriteTime    uint64
	PercentIdleTime         uint64
	SplitIOPerSec           uint32
}

type Win32_PerfFormattedData_PerfDisk_LogicalDisk struct {
	AvgDiskBytesPerRead     uint64
	AvgDiskBytesPerTransfer uint64
	AvgDiskBytesPerWrite    uint64
	AvgDiskQueueLength      uint64
	AvgDiskReadQueueLength  uint64
	AvgDiskSecPerRead       uint32
	AvgDiskSecPerTransfer   uint32
	AvgDiskSecPerWrite      uint32
	AvgDiskWriteQueueLength uint64
	CurrentDiskQueueLength  uint32
	DiskBytesPerSec         uint64
	DiskReadBytesPerSec     uint64
	DiskReadsPerSec         uint32
	DiskTransfersPerSec     uint32
	DiskWriteBytesPerSec    uint64
	DiskWritesPerSec        uint32
	Name                    string
	PercentDiskReadTime     uint64
	PercentDiskTime         uint64
	PercentDiskWriteTime    uint64
	PercentIdleTime         uint64
	SplitIOPerSec           uint32
	FreeMegabytes           uint32
	PercentFreeSpace        uint32
}

type Win32_LogicalDisk struct {
	Name      string
	Size      *uint64 // null if no disk in CD/DVD drive
	FreeSpace *uint64 // null if no disk in CD/DVD drive
	DriveType uint32
}

type Win32_DiskDrive struct {
	Name       string
	Model      string
	Index      uint32
	Partitions uint32
	Size       *uint64
}

const WaitMSec = 500

func Usage(path string) (*UsageStat, error) {
	ret := &UsageStat{}

	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)
	diskret, _, err := procGetDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)))
	if diskret == 0 {
		return nil, err
	}
	ret = &UsageStat{
		Path:        path,
		Total:       uint64(lpTotalNumberOfBytes),
		Free:        uint64(lpTotalNumberOfFreeBytes),
		Used:        uint64(lpTotalNumberOfBytes) - uint64(lpTotalNumberOfFreeBytes),
		UsedPercent: (float64(lpTotalNumberOfBytes) - float64(lpTotalNumberOfFreeBytes)) / float64(lpTotalNumberOfBytes) * 100,
		// InodesTotal: 0,
		// InodesFree: 0,
		// InodesUsed: 0,
		// InodesUsedPercent: 0,
	}
	return ret, nil
}

func Partitions(all bool) ([]PartitionStat, error) {
	var ret []PartitionStat
	lpBuffer := make([]byte, 254)
	diskret, _, err := procGetLogicalDriveStringsW.Call(
		uintptr(len(lpBuffer)),
		uintptr(unsafe.Pointer(&lpBuffer[0])))
	if diskret == 0 {
		return ret, err
	}
	for _, v := range lpBuffer {
		if v >= 65 && v <= 90 {
			path := string(v) + ":"
			if path == "A:" || path == "B:" { // skip floppy drives
				continue
			}
			typepath, _ := windows.UTF16PtrFromString(path)
			typeret, _, _ := procGetDriveType.Call(uintptr(unsafe.Pointer(typepath)))
			if typeret == 0 {
				return ret, windows.GetLastError()
			}
			// 2: DRIVE_REMOVABLE 3: DRIVE_FIXED 4: DRIVE_REMOTE 5: DRIVE_CDROM

			if typeret == 2 || typeret == 3 || typeret == 4 || typeret == 5 {
				lpVolumeNameBuffer := make([]byte, 256)
				lpVolumeSerialNumber := int64(0)
				lpMaximumComponentLength := int64(0)
				lpFileSystemFlags := int64(0)
				lpFileSystemNameBuffer := make([]byte, 256)
				volpath, _ := windows.UTF16PtrFromString(string(v) + ":/")
				driveret, _, err := provGetVolumeInformation.Call(
					uintptr(unsafe.Pointer(volpath)),
					uintptr(unsafe.Pointer(&lpVolumeNameBuffer[0])),
					uintptr(len(lpVolumeNameBuffer)),
					uintptr(unsafe.Pointer(&lpVolumeSerialNumber)),
					uintptr(unsafe.Pointer(&lpMaximumComponentLength)),
					uintptr(unsafe.Pointer(&lpFileSystemFlags)),
					uintptr(unsafe.Pointer(&lpFileSystemNameBuffer[0])),
					uintptr(len(lpFileSystemNameBuffer)))
				if driveret == 0 {
					if typeret == 5 || typeret == 2 {
						continue //device is not ready will happen if there is no disk in the drive
					}
					return ret, err
				}
				opts := "rw"
				if lpFileSystemFlags&FileReadOnlyVolume != 0 {
					opts = "ro"
				}
				if lpFileSystemFlags&FileFileCompression != 0 {
					opts += ".compress"
				}

				d := PartitionStat{
					Mountpoint: path,
					Device:     path,
					Fstype:     string(bytes.Replace(lpFileSystemNameBuffer, []byte("\x00"), []byte(""), -1)),
					Opts:       opts,
				}
				ret = append(ret, d)
			}
		}
	}
	return ret, nil
}

func IOCounters(names ...string) (map[string]IOCountersStat, error) {
	ret := make(map[string]IOCountersStat, 0)
	var dst []Win32_PerfFormattedData

	err := wmi.Query("SELECT * FROM Win32_PerfFormattedData_PerfDisk_LogicalDisk ", &dst)
	if err != nil {
		return ret, err
	}
	for _, d := range dst {

		if len(d.Name) > 3 { // not get _Total or Harddrive
			continue
		}
		if len(names) > 0 && !common.StringsHas(names, d.Name) {
			continue
		}

		ret[d.Name] = IOCountersStat{
			Name:       d.Name,
			ReadCount:  uint64(d.AvgDiskReadQueueLength),
			WriteCount: d.AvgDiskWriteQueueLength,
			ReadBytes:  uint64(d.AvgDiskBytesPerRead),
			WriteBytes: uint64(d.AvgDiskBytesPerWrite),
			ReadTime:   d.AvgDisksecPerRead,
			WriteTime:  d.AvgDisksecPerWrite,
		}
	}
	return ret, nil
}

func PhysicalDisksStats() ([]Win32_PerfFormattedData_PerfDisk_PhysicalDisk, error) {
	var ret []Win32_PerfFormattedData_PerfDisk_PhysicalDisk
	q := wmi.CreateQuery(&ret, "")
	err := wmi.Query(q, &ret)
	return ret, err
}

func LogicalPartitionsStats() ([]Win32_PerfFormattedData_PerfDisk_LogicalDisk, error) {
	var ret []Win32_PerfFormattedData_PerfDisk_LogicalDisk
	q := wmi.CreateQuery(&ret, "")
	err := wmi.Query(q, &ret)
	return ret, err
}

func LogicalDiskSize() ([]Win32_LogicalDisk, error) {
	var ret []Win32_LogicalDisk
	q := wmi.CreateQuery(&ret, "")
	err := wmi.Query(q, &ret)
	return ret, err
}

func Info() ([]Win32_DiskDrive, error) {
	var ret []Win32_DiskDrive
	q := wmi.CreateQuery(&ret, "")
	err := wmi.Query(q, &ret)
	return ret, err
}
