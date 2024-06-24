package svc

import (
	"fmt"
	"os/exec"
)

// StartBind9 启动 BIND9
func StartBind9() error {
	cmd := exec.Command("/usr/sbin/service", "named", "start")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to start bind9: %v, output: %s", err, string(output))
	}
	return nil
}

// ReloadBind9 重载 BIND9
func ReloadBind9() error {
	cmd := exec.Command("/usr/sbin/service", "named", "reload")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reload bind9: %v, output: %s", err, string(output))
	}
	return nil
}
