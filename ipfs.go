package seed

import (
	"bufio"
	"context"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
)

func RunIPFS(ctx context.Context, command string, path string, options ...string) (e error) {
	cmd := exec.CommandContext(ctx, command, options...)
	e = os.Setenv("IPFS_PATH", path)
	if e != nil {
		return e
	}
	cmd.Env = os.Environ()

	//显示运行的命令
	log.Info("[RunIPFS]:", cmd.Args)

	stdout, e := cmd.StdoutPipe()
	if e != nil {
		return e
	}

	stderr, e := cmd.StderrPipe()
	if e != nil {
		return e
	}

	e = cmd.Start()
	if e != nil {
		return e
	}

	reader := bufio.NewReader(io.MultiReader(stdout, stderr))

	//实时循环读取输出流中的一行内容
	for {
		line, e := reader.ReadString('\n')
		if e != nil || io.EOF == e {
			break
		}
		log.Info(line)
	}

	e = cmd.Wait()
	if e != nil {
		return e
	}
	return nil
}
