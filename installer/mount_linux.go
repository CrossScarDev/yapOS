//go:build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

type DeviceInfo struct {
	MountPoint string
	DevicePath string
}

func FindMount(label string) (DeviceInfo, error) {
	devicePath, err := filepath.EvalSymlinks(filepath.Join("/dev/disk/by-label", label))
	if err != nil {
		return DeviceInfo{}, fmt.Errorf("error resolving label: %v", err)
	}

	var st syscall.Stat_t
	if err := syscall.Stat(devicePath, &st); err != nil {
		return DeviceInfo{}, fmt.Errorf("error stating device: %v", err)
	}

	data, err := os.ReadFile("/proc/self/mountinfo")
	if err != nil {
		return DeviceInfo{}, fmt.Errorf("error reading mountinfo: %v", err)
	}

	major := int(st.Rdev / 256)
	minor := int(st.Rdev % 256)

	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 7 {
			continue
		}
		devParts := strings.Split(fields[2], ":")
		if len(devParts) != 2 {
			continue
		}
		currMajor, _ := strconv.Atoi(devParts[0])
		currMinor, _ := strconv.Atoi(devParts[1])
		if currMajor == major && currMinor == minor {
			return DeviceInfo{
				MountPoint: fields[4],
				DevicePath: devicePath,
			}, nil
		}
	}
	return DeviceInfo{}, fmt.Errorf("mount point not found")
}

func UnmountAndEject(info DeviceInfo) error {
	if err := syscall.Unmount(info.MountPoint, 0); err != nil {
		cmd := exec.Command("pkexec", "umount", info.MountPoint)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	ejectCmd := exec.Command("udisksctl", "power-off", "-b", info.DevicePath)
	ejectCmd.Stdout = os.Stdout
	ejectCmd.Stderr = os.Stderr
	return ejectCmd.Run()
}
