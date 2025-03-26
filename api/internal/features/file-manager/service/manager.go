package service

import (
	"fmt"
	"os"
	"time"
)

type FileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	Mode    string    `json:"mode"`
	ModTime time.Time `json:"mod_time"`
	IsDir   bool      `json:"is_dir"`
}

type SSHClient interface {
	NewSftp() (SFTPClient, error)
}

type SFTPClient interface {
	Close() error
	ReadDir(path string) ([]os.FileInfo, error)
	Mkdir(path string) error
	Remove(path string) error
	Stat(path string) (os.FileInfo, error)
}

type SFTPFileInfo interface {
	Name() string
	Size() int64
	Mode() SFTPFileMode
	ModTime() time.Time
	IsDir() bool
}

type SFTPFileMode interface {
	String() string
}

// withSFTPClient safely executes an operation with an SFTP client
func (f *FileManagerService) withSFTPClient(operation func(SFTPClient) error) error {
	if f == nil {
		return fmt.Errorf("file manager service is nil")
	}

	if f.sshpkg == nil {
		return fmt.Errorf("ssh client is nil")
	}

	client, err := f.sshpkg.NewSftp()
	if err != nil {
		return err
	}
	defer client.Close()

	return operation(client)
}

// ListFiles returns a list of files in the given path
func (f *FileManagerService) ListFiles(path string) ([]FileInfo, error) {
	var fileInfos []FileInfo

	err := f.withSFTPClient(func(client SFTPClient) error {
		sftpFileInfos, err := client.ReadDir(path)
		if err != nil {
			return fmt.Errorf("failed to read directory %s: %w", path, err)
		}

		fileInfos = make([]FileInfo, 0, len(sftpFileInfos))
		for _, info := range sftpFileInfos {
			fileInfos = append(fileInfos, FileInfo{
				Name:    info.Name(),
				Size:    info.Size(),
				Mode:    info.Mode().String(),
				ModTime: info.ModTime(),
				IsDir:   info.IsDir(),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileInfos, nil
}

// CreateDirectory creates a new directory at the given path and returns its contents
func (f *FileManagerService) CreateDirectory(path string) ([]FileInfo, error) {
	err := f.withSFTPClient(func(client SFTPClient) error {
		if err := client.Mkdir(path); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return f.ListFiles(path)
}

func (f *FileManagerService) DeleteFile(path string) error {
	err := f.withSFTPClient(func(client SFTPClient) error {
		if err := client.Remove(path); err != nil {
			return fmt.Errorf("failed to delete file %s: %w", path, err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}