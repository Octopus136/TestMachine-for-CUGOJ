package LimitExec

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"syscall"
)

type ExecError struct {
	Info string
}

func (e ExecError) Error() string {
	return e.Info
}

//以限制时间，限制内存的方式执行一个命令
//timeLimit：时间限制 ms
//memoryLimit：空间限制 KB
//args：命令
//stdin：标准输入
//stdout：标准输出
//stderr：标准错误
//返回 时间消耗（ms）内存消耗（KB） 错误信息
func LimitExec(timeLimit, memoryLimit int64, args string, stdin io.Reader, stdout, stderr io.Writer) (time, memery int64, err error) {
	//构造资源限制串，加入用户需要命令
	var str = fmt.Sprintf("ulimit -t %d;ulimit -m %d;", (timeLimit+999)/1000, memoryLimit+32768) + args
	//构造cmd
	print(str)
	cmd := exec.Command("sh", "-c", str)
	//重定向输入输出
	if stdin != nil {
		cmd.Stdin = stdin
	}
	if stdout != nil {
		cmd.Stdout = stdout
	} else {
		cmd.Stdout = &bytes.Buffer{}
	}
	if stderr != nil {
		cmd.Stderr = stderr
	} else {
		cmd.Stderr = &bytes.Buffer{}
	}
	//运行cmd并获取其rusage
	err = cmd.Run()
	statu := (*cmd).ProcessState
	rusage := (*statu).SysUsage().(*syscall.Rusage)

	//计算使用时间和内存
	timeUse := (rusage.Utime.Sec+rusage.Stime.Sec)*1000 + (rusage.Utime.Usec+rusage.Stime.Usec)/1000
	memoryUse := rusage.Maxrss

	if timeUse > timeLimit {
		timeUse = timeLimit
	}
	if memoryUse > memoryLimit {
		memoryUse = memoryLimit
	}

	if timeUse >= timeLimit {
		return timeUse, memoryUse, ExecError{
			Info: "TLE",
		}
	}
	if memoryUse >= memoryLimit {
		return timeUse, memoryUse, ExecError{
			Info: "MLE",
		}
	}
	print("\n")
	print("\n")
	print("\n")
	if err != nil {
		print(err.Error())
	}
	print("\n")
	print("\n")
	print("\n")
	print("\n")
	if err != nil {
		if err.Error() == "exit status 137" {
			timeUse = timeLimit
			return timeUse, memoryUse, ExecError{
				Info: "TLE",
			}
		} else if err.Error() == "exit status 141" {
			return timeUse, memoryUse, ExecError{
				Info: "OLE",
			}
		} else {
			return timeUse, memoryUse, ExecError{
				Info: "RE",
			}
		}
	}
	return timeUse, memoryUse, err
}

//自定义输出的LimitExec  Out
func LimitExecO(timeLimit, memoryLimit int64, args string, stdout *bytes.Buffer) (time, memery int64, err error) {
	return LimitExec(timeLimit, memoryLimit, args, nil, stdout, nil)
}

//自定义错误输出的LimitExec  Err
func LimitExecE(timeLimit, memoryLimit int64, args string, stderr *bytes.Buffer) (time, memery int64, err error) {
	return LimitExec(timeLimit, memoryLimit, args, nil, nil, stderr)
}

//自定义标准输出和错误输出的LimitExec Out/Err
func LimitExecOE(timeLimit, memoryLimit int64, args string, stdout, stderr *bytes.Buffer) (time, memery int64, err error) {
	return LimitExec(timeLimit, memoryLimit, args, nil, stdout, stderr)
}

//不处理输入输出的LimitExec None
func LimitExecN(timeLimit, memoryLimit int64, args string) (time, memery int64, err error) {
	return LimitExec(timeLimit, memoryLimit, args, nil, nil, nil)
}
