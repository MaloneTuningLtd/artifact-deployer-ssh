package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

type SSHSession struct {
	Connection *ssh.Client
	Session    *ssh.Session
}

func Connect() *SSHSession {
	dsn := fmt.Sprintf("%s:%s", config.SSHHost, config.SSHPort)
	conn, err := ssh.Dial("tcp", dsn, config.sshClientConfig())
	if err != nil {
		log.Fatal(err)
	}

	session, err := conn.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	return &SSHSession{conn, session}
}

func (session *SSHSession) Close() {
	session.Session.Close()
	session.Connection.Close()
}

func (session *SSHSession) RunOnce(t *Task) {
	defer session.Close()

	stdout, oerr := session.Session.StdoutPipe()
	if oerr != nil {
		log.Fatal(oerr)
	}

	stderr, eerr := session.Session.StderrPipe()
	if eerr != nil {
		log.Fatal(eerr)
	}

	err := session.Session.Run(t.GetCmd())
	if err != nil {
		log.Fatal(err)
	}

	cout, cerr := make(chan string), make(chan string)

	go func() {
		buf := bytes.NewBuffer(nil)
		io.Copy(buf, stdout)
		cout <- buf.String()
	}()

	go func() {
		buf := bytes.NewBuffer(nil)
		io.Copy(buf, stderr)
		cerr <- buf.String()
	}()

	if out := <-cout; out != "" {
		t.reply = out
	}

	if err := <-cerr; err != "" {
		t.reply = err
	}
}

// https://stackoverflow.com/a/53117797/5332177
func (session *SSHSession) RunTasks(tasks Tasks) {
	cmdLog := log.New(os.Stdout, "[cmd] ", log.Lshortfile)

	// get stdin pipe to shove commands in
	stdin, err := session.Session.StdinPipe()
	defer stdin.Close()

	if err != nil {
		log.Fatal(err)
	}

	// start remote shell
	if err := session.Session.Shell(); err != nil {
		log.Fatal(err)
	}

	// send the commands
	for _, task := range tasks {
		cmd := task.GetCmd()
		_, err := fmt.Fprintf(stdin, "%s\n", cmd)

		cmdLog.Println(cmd)

		if err != nil {
			log.Fatal(err)
		}
	}

	// wait for the session to senpuku
	// if err := session.Session.Wait(); err != nil {
	// 	log.Fatal(err)
	// }
}
