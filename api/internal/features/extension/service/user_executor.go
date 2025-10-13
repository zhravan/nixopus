package service

import (
	"fmt"
	"strings"

	"github.com/raghavyuva/nixopus-api/internal/features/ssh"
)

func (s *ExtensionService) executeUserStep(sshClient *ssh.SSH, props map[string]interface{}, replacer func(string) string, timeout int) (string, error) {
	username, _ := props["username"].(string)
	username = replacer(username)
	action, _ := props["action"].(string)
	shell, _ := props["shell"].(string)
	home, _ := props["home"].(string)
	groups, _ := props["groups"].(string)

	if username == "" {
		return "", fmt.Errorf("username is required for user operations")
	}

	tools := s.userTools(sshClient)
	var cmd string
	switch action {
	case "ensure":
		var err error
		cmd, err = s.buildEnsureUserCmd(tools, username, shell, home, groups)
		if err != nil {
			return "", fmt.Errorf("failed to build ensure user command: %w", err)
		}
	case "delete":
		var err error
		cmd, err = s.buildDeleteUserCmd(tools, username)
		if err != nil {
			return "", fmt.Errorf("failed to build delete user command: %w", err)
		}
	case "modify":
		var err error
		cmd, err = s.buildModifyUserCmd(username, shell, home, groups)
		if err != nil {
			return "", fmt.Errorf("failed to build modify user command: %w", err)
		}
	case "add_groups":
		var err error
		cmd, err = s.buildAddGroupsCmd(username, groups)
		if err != nil {
			return "", fmt.Errorf("failed to build add groups command: %w", err)
		}
	case "remove_groups":
		var err error
		cmd, err = s.buildRemoveGroupsCmd(tools, username, groups)
		if err != nil {
			return "", fmt.Errorf("failed to build remove groups command: %w", err)
		}
	case "check":
		cmd = s.buildCheckUserCmd(username)
	default:
		return "", fmt.Errorf("unsupported user action: %s", action)
	}
	if timeout > 0 {
		cmd = fmt.Sprintf("timeout %ds %s", timeout, cmd)
	}

	output, err := sshClient.RunCommand(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to execute user command '%s': %w (output: %s)", cmd, err, output)
	}
	return output, nil
}

func (s *ExtensionService) buildEnsureUserCmd(tools UserToolset, username string, shell string, home string, groups string) (string, error) {
	var cmd string
	if strings.Contains(tools.add, "useradd") {
		cmd = fmt.Sprintf("id -u %s >/dev/null 2>&1 || sudo useradd %s", username, username)
	} else if strings.Contains(tools.add, "adduser") {
		cmd = fmt.Sprintf("id -u %s >/dev/null 2>&1 || sudo adduser -D %s", username, username)
	} else {
		return "", fmt.Errorf("user management tool not found")
	}
	if shell != "" {
		if strings.Contains(tools.mod, "usermod") {
			cmd += fmt.Sprintf(" && sudo usermod -s %s %s", shell, username)
		} else if strings.Contains(tools.mod, "chsh") {
			cmd += fmt.Sprintf(" && sudo chsh -s %s %s", shell, username)
		}
	}
	if home != "" {
		cmd += fmt.Sprintf(" && sudo usermod -d %s %s", home, username)
	}
	if groups != "" {
		cmd += fmt.Sprintf(" && sudo usermod -aG %s %s", groups, username)
	}
	return cmd, nil
}

func (s *ExtensionService) buildDeleteUserCmd(tools UserToolset, username string) (string, error) {
	if strings.Contains(tools.del, "userdel") {
		return fmt.Sprintf("id -u %s >/dev/null 2>&1 && sudo userdel -r %s || true", username, username), nil
	}
	if strings.Contains(tools.del, "deluser") {
		return fmt.Sprintf("id -u %s >/dev/null 2>&1 && sudo deluser --remove-home %s || true", username, username), nil
	}
	return "", fmt.Errorf("user management tool not found")
}

func (s *ExtensionService) buildModifyUserCmd(username string, shell string, home string, groups string) (string, error) {
	parts := []string{}
	if shell != "" {
		parts = append(parts, fmt.Sprintf("-s %s", shell))
	}
	if home != "" {
		parts = append(parts, fmt.Sprintf("-d %s", home))
	}
	if groups != "" {
		parts = append(parts, fmt.Sprintf("-aG %s", groups))
	}
	if len(parts) == 0 {
		return "", fmt.Errorf("no changes")
	}
	return fmt.Sprintf("sudo usermod %s %s", strings.Join(parts, " "), username), nil
}

func (s *ExtensionService) buildAddGroupsCmd(username string, groups string) (string, error) {
	if groups == "" {
		return "", fmt.Errorf("no groups provided")
	}
	cmds := []string{}
	for _, g := range strings.Split(groups, ",") {
		g = strings.TrimSpace(g)
		if g == "" {
			continue
		}
		cmds = append(cmds, fmt.Sprintf("sudo usermod -aG %s %s", g, username))
	}
	if len(cmds) == 0 {
		return "", fmt.Errorf("no groups provided")
	}
	return strings.Join(cmds, " && "), nil
}

func (s *ExtensionService) buildRemoveGroupsCmd(tools UserToolset, username string, groups string) (string, error) {
	if groups == "" {
		return "", fmt.Errorf("no groups provided")
	}
	cmds := []string{}
	for _, g := range strings.Split(groups, ",") {
		g = strings.TrimSpace(g)
		if g == "" {
			continue
		}
		if strings.Contains(tools.groupDel, "gpasswd") {
			cmds = append(cmds, fmt.Sprintf("sudo gpasswd -d %s %s || true", username, g))
		} else if strings.Contains(tools.groupDel, "deluser") {
			cmds = append(cmds, fmt.Sprintf("sudo deluser %s %s || true", username, g))
		}
	}
	if len(cmds) == 0 {
		return "", fmt.Errorf("no groups provided")
	}
	return strings.Join(cmds, " && "), nil
}

func (s *ExtensionService) buildCheckUserCmd(username string) string {
	return fmt.Sprintf("id -u %s >/dev/null 2>&1 && echo exists || echo missing", username)
}
