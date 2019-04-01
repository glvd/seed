package seed

import (
	"bufio"
	"context"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"io"
	"os"
	"os/exec"
)

// RunIPFS ...
func RunIPFS(ctx context.Context, path string, command string, options ...string) (e error) {
	defer func() {
		if e != nil {
			log.Error(e)
		}
	}()
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
		select {
		case <-ctx.Done():
			e = xerrors.New("exit with done")
			return
		default:
			line, _, e := reader.ReadLine()
			if e != nil || io.EOF == e {
				goto END
			}
			log.Debug(string(line))
		}

	}
END:
	e = cmd.Wait()
	if e != nil {
		return e
	}
	return nil
}
