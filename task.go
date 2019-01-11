package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

type Task struct {
	command string
	args    []string
	reply   string
}

type Tasks []Task

func HasMkDir(dir string) Task {
	cmd := fmt.Sprintf(`if [ ! -d "%[1]v" ]; then mkdir -p "%[1]v"; fi`, dir)

	return Task{
		command: cmd,
	}
}

func SymlinkSharedDirs(dirs []string, releasePath string) (tasks Tasks) {
	sharedPath := filepath.Join(config.ProjectPath, "shared")
	hasDir := TestOnce(fmt.Sprintf(`[[ -d "%s" ]]`, sharedPath))

	if hasDir {
		for _, d := range dirs {
			sharedDir := filepath.Join(sharedPath, d)
			releaseDir := filepath.Join(releasePath, d)

			tasks = append(tasks, Symlink(sharedDir, releaseDir))
		}
	}

	return
}

func Symlink(dir1, dir2 string) Task {
	var symlink, opt string

	if !supportSymlink {
		log.Fatal("symlink is not supported!")
	} else {
		symlink = "ln"
		opt = "--symbolic"
	}

	return Task{
		command: symlink,
		args:    []string{opt, dir1, dir2},
	}
}

func HasRmSymlink(symlinkDir string) (task Task) {
	if supportSymlink {
		task.command = fmt.Sprintf(`if [ -L "%[1]v" ]; then rm "%[1]v"; fi`, symlinkDir)
	}

	return
}

func (t Task) GetCmd() string {
	return fmt.Sprintf("%s %s", t.command, strings.Join(t.args, " "))
}

func Test(cmd string) Task {
	return Task{
		command: fmt.Sprintf("if %s; then echo 'true'; fi", cmd),
	}
}

func TestOnce(cmd string) bool {
	s := Connect()
	defer s.Close()

	t := Test(cmd)
	s.RunOnce(&t)

	return t.reply != "true"
}

func SupportsCmd(cmd, opt string) bool {
	command := fmt.Sprintf("[[ $(man %[1]s 2>&1 || %[1]s -h 2>&1 || %[1]s --help 2>&1) =~ '%[2]s' ]]", cmd, opt)

	return TestOnce(command)
}
