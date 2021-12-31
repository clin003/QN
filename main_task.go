package main

import (
	"os"
	"os/exec"
)

func restart() error {
	if *version {
		return nil
	}
	// 设置传递给子进程的参数
	// args := []string{}
	var args []string
	if len(*cfg) > 0 {
		args = append(args, "-c")
		args = append(args, *cfg)
	}
	if len(*workdir) > 0 {
		args = append(args, "-w")
		args = append(args, *workdir)
	}
	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdout = os.Stdout // 标准输出
	cmd.Stderr = os.Stderr // 错误输出
	// 新建并执行子进程
	// cmd.Start()
	// cmd.Wait()
	return cmd.Start()
	// return cmd.Run()
}
