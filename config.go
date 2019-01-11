package main

import (
	"bytes"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	SSHHost                 string
	SSHPort                 string
	SSHUsername             string
	SSHPrivateKey           *ssh.Signer
	SSHPrivateKeyPassPhrase []byte

	Binary      string
	Tag         string
	ProjectPath string
	SharedDirs  []string
}

func getEnv(key string) (result string, found bool) {
	result, found = os.LookupEnv(key)
	if !found {
		result, found = os.LookupEnv("PLUGIN_" + key)
	}

	return
}

func envNotFound(key string) {
	log.Fatalf("%s should be defined \n", key)
}

func loadIntoConfig(c *Config) {
	if host, f := getEnv("SSH_HOST"); f {
		c.SSHHost = host
	} else {
		envNotFound("SSH_HOST")
	}

	if user, f := getEnv("SSH_USER"); f {
		c.SSHUsername = user
	} else {
		envNotFound("SSH_USER")
	}

	if passphrase, f := getEnv("SSH_PASSPHRASE"); f {
		c.SSHPrivateKeyPassPhrase = []byte(passphrase)
	}

	if port, p := getEnv("SSH_PORT"); p {
		c.SSHPort = port
	} else {
		c.SSHPort = "22"
	}

	if keyPrivate, k := getEnv("SSH_PRIVATEKEY"); k {
		var (
			signer ssh.Signer
			err    error
		)

		keyBytes := []byte(keyPrivate)
		keyBytes = bytes.Replace(keyBytes, []byte("\\n"), []byte("\n"), -1)

		if len(c.SSHPrivateKeyPassPhrase) > 0 {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(keyBytes, c.SSHPrivateKeyPassPhrase)
		} else {
			signer, err = ssh.ParsePrivateKey(keyBytes)
		}

		if err != nil {
			log.Fatal(err)
		}

		c.SSHPrivateKey = &signer
	}

	if b, k := getEnv("BINARY"); !k {
		err := errors.New("Binary should be defined")
		log.Fatal(err)
	} else {
		binPath, err := filepath.Abs(b)
		if err != nil {
			err = errors.Wrap(err, "Invalid Path")
			log.Fatal(err)
		}

		if _, err := os.Stat(binPath); err != nil {
			// TODO: use file os error wrap checker func
			err = errors.Wrap(err, "Binary File Not Found")
			log.Fatal(err)
		}

		c.Binary = binPath
	}

	if t, k := getEnv("TAG"); !k {
		err := errors.New("Tag should be defined")
		log.Fatal(err)
	} else {
		c.Tag = t
	}

	if projectPath, k := getEnv("PROJECT_PATH"); !k {
		err := errors.New("Project Path should be defined")
		log.Fatal(err)
	} else {
		c.ProjectPath = projectPath
	}

	if directories, k := getEnv("SHARED_DIRS"); k {
		delimit := regexp.MustCompile("[\\,\\s+]")
		c.SharedDirs = delimit.Split(directories, -1)
	}
}

func (c Config) sshClientConfig() *ssh.ClientConfig {
	cfg := &ssh.ClientConfig{
		User: c.SSHUsername,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(*c.SSHPrivateKey)},
		// TODO: most likely unsafe
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	cfg.SetDefaults()

	return cfg
}
