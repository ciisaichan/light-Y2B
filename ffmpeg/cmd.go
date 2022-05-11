package ffmpeg

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func StartCmd(ctx context.Context, ffmpegPath string, params []string, hideLog bool) error {
	c := exec.CommandContext(ctx, ffmpegPath, params...)
	SetCmdEnvironment(c)
	stdout, err := c.StderrPipe()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				buff := make([]byte, 1024)
				_, err := reader.Read(buff)
				if err != nil || err == io.EOF {
					return
				}
				if !hideLog {
					fmt.Print(string(buff))
				}
				/*
					readString, err := reader.ReadString('\r')
					if err != nil || err == io.EOF {
						return
					}
					readString = strings.Replace(readString, "\n", "", -1)

					if !debug && strings.HasPrefix(readString, "frame=") {
						logger.L.Infof("[FFmpeg] %s", readString)
					} else if debug {
						logger.L.Infof("[FFmpeg] %s", readString)
					}
				*/
			}
		}
	}(&wg)
	err = c.Start()
	wg.Wait()
	return err
}

func SetCmdEnvironment(cmd *exec.Cmd) error {
	env := os.Environ()
	cmdEnv := []string{}

	for _, e := range env {
		i := strings.Index(e, "=")
		if i > 0 && (e[:i] == "ENV_NAME") {
		} else {
			cmdEnv = append(cmdEnv, e)
		}
	}
	cmd.Env = cmdEnv

	return nil
}
