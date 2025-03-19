package ssh

import (
	"log"

	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

type SSH struct {
	PrivateKey          string `json:"private_key"`
	PublicKey           string `json:"public_key"`
	Host                string `json:"host"`
	User                string `json:"user"`
	Port                uint   `json:"port"`
	Password            string `json:"password"`
	PrivateKeyProtected string `json:"private_key_protected"`
}

func NewSSH() *SSH {
	return &SSH{}
}

func (s *SSH) ConnectWithPassword() (*goph.Client, error) {
	auth := goph.Password(s.Password)

	client, err := goph.NewConn(&goph.Config{
		User:     s.User,
		Addr:     s.Host,
		Port:     uint(s.Port),
		Auth:     auth,
		Callback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Fatalf("SSH connection failed: %v", err)
	}

	defer client.Close()
	return client, nil
}

func (s *SSH) ConnectWithPrivateKey() (*goph.Client, error) {
	auth, err := goph.Key(s.PrivateKey, "")

	if err != nil {
		log.Fatalf("SSH connection failed: %v", err)
	}

	client, err := goph.NewConn(&goph.Config{
		User:     s.User,
		Addr:     s.Host,
		Port:     uint(s.Port),
		Auth:     auth,
		Callback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Fatalf("SSH connection failed: %v", err)
	}

	defer client.Close()
	return client, nil
}

func (s *SSH) ConnectWithPrivateKeyProtected() (*goph.Client, error) {
	auth, err := goph.Key(s.PrivateKeyProtected, "")

	if err != nil {
		log.Fatalf("SSH connection failed: %v", err)
	}

	client, err := goph.NewConn(&goph.Config{
		User:     s.User,
		Addr:     s.Host,
		Port:     uint(s.Port),
		Auth:     auth,
		Callback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Fatalf("SSH connection failed: %v", err)
	}

	defer client.Close()
	return client, nil
}
