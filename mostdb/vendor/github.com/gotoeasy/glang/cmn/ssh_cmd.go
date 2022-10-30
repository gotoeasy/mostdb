package cmn

import (
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// 远程ssh执行命令
func SshCmd(host string, port string, user string, password string, cmd ...string) (string, error) {

	config := &ssh.ClientConfig{
		Timeout:         time.Second * 3,
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
	}

	sshClient, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		log.Fatal("远程登录失败", err)
		return "", err
	}
	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatal("远程会话失败", err)
		return "", err
	}
	defer session.Close()

	combo, err := session.CombinedOutput(strings.Join(cmd, "\n"))
	if err != nil {
		log.Fatal("执行命令失败", err)
	}
	return string(combo), err
}
