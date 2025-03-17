//go:build windows

package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type DeviceInfo struct {
	MountPoint string
	VolumeName string
}

func FindMount(label string) (DeviceInfo, error) {
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	getDrives := kernel32.NewProc("GetLogicalDriveStringsW")
	getVolume := kernel32.NewProc("GetVolumeInformationW")

	buf := make([]uint16, 256)
	ret, _, _ := getDrives.Call(
		uintptr(len(buf)),
		uintptr(unsafe.Pointer(&buf[0])),
	)

	if ret == 0 {
		return DeviceInfo{}, fmt.Errorf("failed to get drives")
	}

	drives := windows.UTF16ToString(buf)
	for _, drive := range strings.Split(drives, "\x00") {
		drive = strings.TrimSpace(drive)
		if len(drive) < 3 {
			continue
		}

		var volName [256]uint16
		var fsName [256]uint16
		ret, _, _ := getVolume.Call(
			uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(drive))),
			uintptr(unsafe.Pointer(&volName[0])),
			uintptr(len(volName)),
			0, 0, 0,
			uintptr(unsafe.Pointer(&fsName[0])),
			uintptr(len(fsName)),
		)

		if ret != 0 && windows.UTF16ToString(volName[:]) == label {
			return DeviceInfo{
				MountPoint: filepath.Clean(drive),
				VolumeName: windows.UTF16ToString(volName[:]),
			}, nil
		}
	}
	return DeviceInfo{}, fmt.Errorf("drive not found")
}

func UnmountAndEject(info DeviceInfo) error {
	script := fmt.Sprintf(
		`$drive = Get-Volume -FileSystemLabel "%s";`+
			`Remove-PartitionAccessPath -DiskNumber $drive.DiskNumber -PartitionNumber $drive.PartitionNumber -AccessPath "%s";`+
			`Dismount-Disk -Number $drive.DiskNumber -Confirm:$false;`+
			`Set-Disk -Number $drive.DiskNumber -IsOffline $true`,
		info.VolumeName, info.MountPoint)

	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf(`Start-Process -Verb RunAs PowerShell -ArgumentList "-Command", "%s"`, script))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}
