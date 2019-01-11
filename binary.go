package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
)

func (s *SSHSession) SFTPClient() *sftp.Client {
	conn := s.Connection

	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func getBinary() *os.File {
	f, err := os.Open(config.Binary)
	if err != nil {
		log.Fatal(err)
	}

	return f
}

func (s *SSHSession) PutBinary(binary *os.File, releasePath string) error {
	sftpClient := s.SFTPClient()

	name := filepath.Base(binary.Name())
	copyPath := filepath.Join(releasePath, name)

	file, err := sftpClient.Create(copyPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if _, err := io.Copy(file, binary); err != nil {
		return err
	}

	if err := file.Chmod(os.FileMode(0755)); err != nil {
		log.Println("unable to assign binary execute")
	}

	return nil
}
