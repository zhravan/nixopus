package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/raghavyuva/nixopus-api/internal/features/ssh"
)

func executeFileStep(sshClient *ssh.SSH, props map[string]interface{}, replacer func(string) string) (string, error) {
	action, _ := props["action"].(string)
	src, _ := props["src"].(string)
	dest, _ := props["dest"].(string)
	src = replacer(src)
	dest = replacer(dest)

	// For mkdir action, dest is required but src is not
	if action == "mkdir" && dest == "" {
		return "", fmt.Errorf("dest is required for mkdir action")
	}

	switch action {
	case "move":
		client, err := sshClient.Connect()
		if err != nil {
			return "", fmt.Errorf("failed to connect SSH for move operation: %w", err)
		}
		s, err := client.NewSftp()
		if err != nil {
			return "", fmt.Errorf("failed to create SFTP session for move operation: %w", err)
		}
		defer s.Close()
		if err := s.Rename(src, dest); err != nil {
			return "", fmt.Errorf("failed to move %s to %s: %w", src, dest, err)
		}
		return fmt.Sprintf("moved %s to %s", src, dest), nil
	case "copy":
		// Fallback to remote cp for recursive copy
		output, err := sshClient.RunCommand(fmt.Sprintf("cp -r %s %s", src, dest))
		if err != nil {
			return "", fmt.Errorf("failed to copy %s to %s: %w (output: %s)", src, dest, err, output)
		}
		return fmt.Sprintf("copied %s to %s", src, dest), nil
	case "upload":
		client, err := sshClient.Connect()
		if err != nil {
			return "", fmt.Errorf("failed to connect SSH for upload operation: %w", err)
		}
		s, err := client.NewSftp()
		if err != nil {
			return "", fmt.Errorf("failed to create SFTP session for upload operation: %w", err)
		}
		defer s.Close()
		f, err := os.Open(src)
		if err != nil {
			return "", fmt.Errorf("failed to open source file %s: %w", src, err)
		}
		defer f.Close()
		if err := s.MkdirAll(filepath.Dir(dest)); err != nil {
			return "", fmt.Errorf("failed to create destination directory %s: %w", filepath.Dir(dest), err)
		}
		out, err := s.Create(dest)
		if err != nil {
			return "", fmt.Errorf("failed to create destination file %s: %w", dest, err)
		}
		defer out.Close()
		if _, err := io.Copy(out, f); err != nil {
			return "", fmt.Errorf("failed to copy file content from %s to %s: %w", src, dest, err)
		}
		return fmt.Sprintf("uploaded to %s", dest), nil
	case "delete":
		client, err := sshClient.Connect()
		if err != nil {
			return "", fmt.Errorf("failed to connect SSH for delete operation: %w", err)
		}
		s, err := client.NewSftp()
		if err != nil {
			return "", fmt.Errorf("failed to create SFTP session for delete operation: %w", err)
		}
		defer s.Close()
		if err := s.Remove(dest); err != nil {
			return "", fmt.Errorf("failed to delete %s: %w", dest, err)
		}
		return fmt.Sprintf("deleted %s", dest), nil
	case "mkdir":
		client, err := sshClient.Connect()
		if err != nil {
			return "", fmt.Errorf("failed to connect SSH for mkdir operation: %w", err)
		}
		s, err := client.NewSftp()
		if err != nil {
			return "", fmt.Errorf("failed to create SFTP session for mkdir operation: %w", err)
		}
		defer s.Close()
		if err := s.MkdirAll(dest); err != nil {
			return "", fmt.Errorf("failed to create directory %s: %w", dest, err)
		}
		return fmt.Sprintf("mkdir %s", dest), nil
	default:
		return "", fmt.Errorf("unsupported file action: %s", action)
	}
}
