package utils

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io"
	"net"
	"os"
	"path"
	"strconv"
	"time"
)

/*
@Author : VictorTu
@Software: GoLand
*/

func connect(user, password, host string, port int64) (*ssh.Client, *sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connect to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return sshClient, nil, err
	}

	return sshClient, sftpClient, nil
}

type sftpUtil struct {
}

var SftpUtil sftpUtil

type SftpInfo struct {
	Username string
	Password string
	Host     string
	Port     string
	Dir      string
}

type SrcDistPath struct {
	Src  string
	Dist string
}

func (this *sftpUtil) ConnectGet(username, password, host, portS, remoteFileUri string, localFileUri string) (err error) {
	var (
		sftpClient *sftp.Client
		sshClient  *ssh.Client
	)
	// 这里换成实际的 SSH 连接的 用户名，密码，主机名或IP，SSH端口
	port, _ := strconv.ParseInt(portS, 10, 64)
	sshClient, sftpClient, err = connect(username, password, host, port)
	if err != nil {
		return err
	}

	defer func() {
		if sftpClient != nil {
			sftpClient.Close()
		}
		if sshClient != nil {
			sshClient.Close()
		}
	}()

	srcFile, err := sftpClient.Open(remoteFileUri)
	if err != nil {
		return err
	}
	defer func() {
		if srcFile != nil {
			srcFile.Close()
		}
	}()

	dstFile, err := os.Create(localFileUri)
	if err != nil {
		return err
	}
	defer func() {
		if dstFile != nil {
			dstFile.Close()
		}
	}()
	if _, err := srcFile.WriteTo(dstFile); err != nil {
		return err
	}
	return nil
}

func (this *sftpUtil) PushFile(username, password, host, port, localFileUri, remoteFileUri string) error {

	var auths []ssh.AuthMethod
	if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))

	}

	if password != "" {
		auths = append(auths, ssh.Password(password))
	}

	config := ssh.ClientConfig{
		User:            username,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := host + ":" + port
	conn, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return err
	}
	defer conn.Close()

	c, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer c.Close()

	dir := path.Dir(remoteFileUri)
	if err = c.MkdirAll(dir); err != nil {
		return err
	}
	f, err := os.Open(localFileUri)
	if err != nil {
		return err
	}
	defer f.Close()
	w, err := c.Create(remoteFileUri)
	if err != nil {
		return err
	}
	defer w.Close()
	buffer := make([]byte, 1024*5)
	for {
		n, err := f.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		w.Write(buffer[0:n])
	}

	return nil
}

var defaultBufferSize = 1024 * 512
