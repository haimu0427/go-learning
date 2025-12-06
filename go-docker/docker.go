package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("help")
	}
}

// 1. run 函数：这是宿主机视角的“入口”
// 它的作用是配置 Namespace，然后在一个新的隔离环境中“自我调用” child 函数
func run() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	// /proc/self/exe 指向当前正在运行的程序本身（也就是 mydocker）
	// 我们实际上是调用：./mydocker child [你的命令]
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	// 连接标准输入输出，这样你才能在容器里打字看到结果
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 【核心重点】：配置 Namespace
	// CLONE_NEWUTS: 隔离主机名 (Unix Timesharing System)
	// CLONE_NEWPID: 隔离进程 ID
	// CLONE_NEWNS:  隔离挂载点 (Mount Namespace)，防止容器内挂载影响宿主机
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running the /bin/sh command - %s\n", err)
		os.Exit(1)
	}
}

// 2. child 函数：这是容器视角的“入口”
// 当代码执行到这里时，已经处于被隔离的 Namespace 内部了
func child() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	// 系统调用：设置容器内的主机名
	syscall.Sethostname([]byte("container"))

	// 配置挂载点
	// 这一步是为了让 ps 命令能工作。ps 命令需要读取 /proc 目录。
	// 如果不重新挂载 /proc，容器内看到的还是宿主机的进程列表。
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "")

	// 执行用户请求的命令（比如 /bin/sh）
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running the child process - %s\n", err)
		os.Exit(1)
	}

	// 退出时取消挂载（虽然进程结束 Namespace 销毁时会自动清理，但这是好习惯）
	syscall.Unmount("proc", 0)
}
