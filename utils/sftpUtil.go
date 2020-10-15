package utils

import (
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
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
		return nil, nil, errors.New(fmt.Sprintf("ssh.Dial tcp new sshClient err %s ", err.Error()))
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return sshClient, nil, errors.New(fmt.Sprintf("sshClient new sftpClient err %s ", err.Error()))
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

	var remoteFile *sftp.File
	if remoteFile, err = sftpClient.Open(remoteFileUri); err != nil {
		return errors.New(fmt.Sprintf("sftpClient.Open remoteFile err %s ", err.Error()))
	}
	defer func() {
		if remoteFile != nil {
			remoteFile.Close()
		}
	}()

	var localFile *os.File
	if localFile, err = os.Create(localFileUri); err != nil {
		return errors.New(fmt.Sprintf("os.Create localFile err %s ", err.Error()))
	}
	defer func() {
		if localFile != nil {
			localFile.Close()
		}
	}()
	if _, err := remoteFile.WriteTo(localFile); err != nil {
		return errors.New(fmt.Sprintf("remoteFile.WriteTo localFile err %s ", err.Error()))
	}
	return nil
}

func (this *sftpUtil) PushFile(username, password, host, port, localFileUri, remoteFileUri string) (err error) {

	var auths []ssh.AuthMethod

	auths = make([]ssh.AuthMethod, 0)
	auths = append(auths, ssh.Password(password))

	if password != "" {
		auths = append(auths, ssh.Password(password))
	}

	config := ssh.ClientConfig{
		User:            username,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := host + ":" + port
	var conn *ssh.Client
	if conn, err = ssh.Dial("tcp", addr, &config); err != nil {
		return errors.New(fmt.Sprintf("ssh.Dial err %s ", err.Error()))
	}
	defer conn.Close()

	var sftpClient *sftp.Client
	if sftpClient, err = sftp.NewClient(conn); err != nil {
		return errors.New(fmt.Sprintf("sftp.NewClient err %s ", err.Error()))
	}
	defer sftpClient.Close()

	dir := path.Dir(remoteFileUri)
	if err = sftpClient.MkdirAll(dir); err != nil {
		return errors.New(fmt.Sprintf("sftpClient MkdirAll err %s ", err.Error()))
	}
	var f *os.File
	if f, err = os.Open(localFileUri); err != nil {
		return errors.New(fmt.Sprintf("os.Open local file err %s ", err.Error()))
	}
	defer f.Close()

	var sftpFile *sftp.File
	if sftpFile, err = sftpClient.Create(remoteFileUri); err != nil {
		return errors.New(fmt.Sprintf("sftpClient create file err %s ", err.Error()))
	}
	defer sftpFile.Close()
	buffer := make([]byte, defaultBufferSize)
	for {
		n, err := f.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		sftpFile.Write(buffer[0:n])
	}

	return nil
}

var defaultBufferSize = 1024 * 32 // 32k buffer
