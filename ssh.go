package gocommon

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"time"
)

const (
	CERT_PASSWORD        = 1
	CERT_PUBLIC_KEY_FILE = 2
	DEFAULT_TIMEOUT      = 3 // second
)

type FileList struct {
	SourceFilePath      *string
	SourceFileReader    *bytes.Reader
	DestinationFilePath string
}

type SSH struct {
	Ip      string
	User    string
	Cert    string //password or key file path
	session *ssh.Session
	client  *ssh.Client
}

func (ssh_client *SSH) readPublicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, appendError(err, "ReadFile ")
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, appendError(err, "ParsePrivateKey ")
	}
	return ssh.PublicKeys(key), nil
}

func (ssh_client *SSH) Connect(mode int) error {

	var ssh_config *ssh.ClientConfig
	var auth []ssh.AuthMethod
	if mode == CERT_PASSWORD {
		auth = []ssh.AuthMethod{ssh.Password(ssh_client.Cert)}
	} else if mode == CERT_PUBLIC_KEY_FILE {
		key, err := ssh_client.readPublicKeyFile(ssh_client.Cert)
		if err != nil {
			return appendError(err, "readPublicKeyFile ")
		}
		auth = []ssh.AuthMethod{key}
	} else {
		return makeError("Mode not supported: %d", mode)
	}

	ssh_config = &ssh.ClientConfig{
		User: ssh_client.User,
		Auth: auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: time.Second * 5,
	}

	client, err := ssh.Dial("tcp", ssh_client.Ip, ssh_config)
	if err != nil {
		return appendError(err, "Dial ")
	}

	session, err := client.NewSession()
	if err != nil {
		_ = client.Close()
		return appendError(err, "NewSession ")
	}

	ssh_client.session = session
	ssh_client.client = client
	return nil
}

func RemoteAttachAndRun(user string, addr string, privateKey string, cmd string, stdin *os.File, stdout *os.File, stderr *os.File) error {
	client := &SSH{
		Ip:   addr,
		User: user,
		Cert: privateKey,
	}
	err := client.Connect(CERT_PUBLIC_KEY_FILE)
	if err != nil {
		return appendError(err, "Connect ")
	}
	err = client.RunAttachCmd(cmd, stdin, stdout, stderr)
	client.Close()
	return appendError(err, "RunAttachCmd ")
}

func (ssh_client *SSH) RunAttachCmd(cmd string, stdin *os.File, stdout *os.File, stderr *os.File) error {
	if stdin == nil {
		stdin = os.Stdin
	}
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}
	ssh_client.session.Stdin = stdin
	ssh_client.session.Stdout = stdout
	ssh_client.session.Stderr = stderr
	fileDescriptor := int(stdin.Fd())
	var err error
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if terminal.IsTerminal(fileDescriptor) {
		originalState, err := terminal.MakeRaw(fileDescriptor)
		if err != nil {
			return appendError(err, "terminal.MakeRaw ")
		}
		defer func() { _ = terminal.Restore(fileDescriptor, originalState) }()

		termWidth, termHeight, err := terminal.GetSize(fileDescriptor)
		if err != nil {
			return appendError(err, "terminal.GetSize ")
		}

		err = ssh_client.session.RequestPty("xterm-256color", termHeight, termWidth, modes)
		if err != nil {
			return appendError(err, "RequestPty ")
		}
	} else {
		err = ssh_client.session.RequestPty("vt100", 24, 80, modes)
		if err != nil {
			return appendError(err, "RequestPty ")
		}
	}
	err = ssh_client.session.Run(cmd)
	if err != nil {
		return appendError(err, "session.Run ")
	}
	return nil
}

func (ssh_client *SSH) RunCmd(cmd string) ([]byte, error) {
	out, err := ssh_client.session.CombinedOutput(cmd)
	if err != nil {
		return out, appendError(err, "session.CombinedOutput ")
	}
	return out, nil
}

func (ssh_client *SSH) Close() {
	_ = ssh_client.session.Close()
	_ = ssh_client.client.Close()
}

func RemoteRun(user string, addr string, privateKey string, cmd string) ([]byte, error) {
	client := &SSH{
		Ip:   addr,
		User: user,
		Cert: privateKey,
	}
	err := client.Connect(CERT_PUBLIC_KEY_FILE)
	if err != nil {
		return nil, appendError(err, "Connect ")
	}
	ret, err := client.RunCmd(cmd)
	client.Close()
	return ret, appendError(err, "RunCmd ")
}

func remoteSession(user string, addr string, privateKey string) (*SSH, error) {
	client := &SSH{
		Ip:   addr,
		User: user,
		Cert: privateKey,
	}
	err := client.Connect(CERT_PUBLIC_KEY_FILE)
	if err != nil {
		return nil, appendError(err, "Connect ")
	}
	return client, nil
}

func Scp(user string, addr string, privateKey string, files []FileList) error {

	for _, file := range files {
		sess, err := remoteSession(user, addr, privateKey)
		if err != nil {
			return appendError(err, "remoteSession ")
		}
		session := sess.session

		var size int64
		var f *os.File
		if file.SourceFilePath != nil {
			f, err = os.Open(*file.SourceFilePath)
			if err != nil {
				return appendError(err, "os.Open ")
			}
			defer func() { _ = f.Close() }()
			fi, err := f.Stat()
			if err != nil {
				return appendError(err, "f.Stat ")
			}
			size = fi.Size()
		} else {
			if file.SourceFileReader == nil {
				return makeError("Either SourceFilePath or SourceFileReader must be set for %s", file.DestinationFilePath)
			}
			size = file.SourceFileReader.Size()
		}
		go func() {
			w, err := session.StdinPipe()
			if err != nil {
				return
			}
			defer func() { _ = w.Close() }()
			_, _ = fmt.Fprintln(w, "C"+"0755", size, path.Base(file.DestinationFilePath))
			if file.SourceFilePath != nil {
				_, _ = io.Copy(w, f)
			} else {
				_, _ = io.Copy(w, file.SourceFileReader)
			}
			_, _ = fmt.Fprintln(w, "\x00")
		}()

		err = session.Run("/usr/bin/scp -qt " + path.Dir(file.DestinationFilePath))

		if err != nil && err.Error() != "Process exited with status 1" {
			return appendError(err, "session.Run(scp) ")
		}
		_ = session.Close()
		sess.Close()
	}

	return nil
}
