package main

import (
	"fmt"
	"path/filepath"
)

var (
	config         Config
	supportSymlink bool
)

func prerelease(s *SSHSession, releasePath string) {
	symlinkPath := filepath.Join(config.ProjectPath, "current")

	tasks := Tasks{
		HasMkDir(config.ProjectPath),
		HasMkDir(releasePath),
		HasRmSymlink(symlinkPath),
		Symlink(releasePath, symlinkPath),
	}

	symlinks := SymlinkSharedDirs(config.SharedDirs, releasePath)
	tasks = append(tasks, symlinks...)

	s.RunTasks(tasks)
}

func release(s *SSHSession, releasePath string) {
	binaryFile := getBinary()
	defer binaryFile.Close()

	s.PutBinary(binaryFile, releasePath)
}

func main() {
	loadIntoConfig(&config)

	fmt.Println("artifact-deployer-ssh")

	// test for symlink support
	supportSymlink = SupportsCmd("ln", "--symbolic")

	// join/make release path
	releasePath := filepath.Join(config.ProjectPath, "releases", config.Tag)

	// start release deploy
	s := Connect()
	defer s.Close()

	// 1. prerelease
	// - ensure binary file exists
	// - create the right directories
	// - create the symlink(s)
	prerelease(s, releasePath)

	// 2. release
	// - copy binary over
	// - give binary execute permissions
	release(s, releasePath)
}
