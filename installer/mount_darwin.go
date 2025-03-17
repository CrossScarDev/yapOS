//go:build darwin

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type DeviceInfo struct {
	MountPoint string
	DeviceNode string
}

func FindMount(label string) (DeviceInfo, error) {
	deviceNode := ""

	out, err := exec.Command("diskutil", "info", "-plist", label).Output()
	if err == nil {
		for _, line := range strings.Split(string(out), "\n") {
			if strings.Contains(line, "<key>DeviceNode</key>") {
				parts := strings.Split(strings.TrimSpace(line), "<string>")
				if len(parts) > 1 {
					deviceNode = strings.TrimSuffix(parts[1], "</string>")
					break
				}
			}
		}
	}

	mounts, _ := exec.Command("mount").Output()
	for _, line := range strings.Split(string(mounts), "\n") {
		if strings.Contains(line, deviceNode) {
			parts := strings.Fields(line)
			if len(parts) > 2 {
				return DeviceInfo{
					MountPoint: parts[2],
					DeviceNode: deviceNode,
				}, nil
			}
		}
	}

	return DeviceInfo{}, fmt.Errorf("mount point not found")
}

func UnmountAndEject(info DeviceInfo) error {
	script := fmt.Sprintf(`
    tell application "System Events"
        activate
        do shell script "diskutil unmount '%s' && diskutil eject '%s'" with administrator privileges
    end tell
    `, info.MountPoint, info.DeviceNode)

	cmd := exec.Command("osascript", "-e", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
