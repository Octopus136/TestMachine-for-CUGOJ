package main

import (
	"TestMachine/src/Tester"
	"fmt"
	"os"
	"strconv"
)

const types = "'gnu','java'"

/*
args[1] 编译器类型(gnu,java)
args[2] 语言版本(
	gnu:c99,c11,cpp11,cpp14,cpp17,cpp20
	java:jdk
)
args[3] 执行方式(compile,run,spjrun)
args[4] 源文件、可执行文件路径，不包含后缀，例如/code/main，不需要/code/main.c
args[5] 执行时间限制 ms
args[6] 执行空间限制 MB
args[7] 测试数据路径，不包含后缀，默认后缀分别为.in和.out，如果是spj，允许.out文件不存在
args[8] spj路径
compile方式要求有7个参数
run方式要求有8个参数
spj run方式要求9个参数
*/
func main() {
	var args = os.Args
	if len(args) < 4 {
		show(Tester.TestInfo{
			Statu: "004",
			Info:  "调用参数过少，参数数量为 " + strconv.Itoa(len(args)) + " ,但是至少需要7个参数执行:对于compile模式，需要提供7个参数;对于run模式需要提供8个参数;对于spjrun模式需要提供9个参数",
		})
		return
	} else {
		if args[3] == "compile" {
			if len(args) < 7 {
				show(Tester.TestInfo{
					Statu: "004",
					Info:  "调用参数过少，参数数量为 " + strconv.Itoa(len(args)) + " ,但是compile模式至少需要7个参数执行",
				})
				return
			}
		} else if args[3] == "run" {
			if len(args) < 8 {
				show(Tester.TestInfo{
					Statu: "004",
					Info:  "调用参数过少，参数数量为 " + strconv.Itoa(len(args)) + " ,但是run模式至少需要8个参数执行",
				})
				return
			}
		} else if args[3] == "spjrun" {
			if len(args) < 9 {
				show(Tester.TestInfo{
					Statu: "004",
					Info:  "调用参数过少，参数数量为 " + strconv.Itoa(len(args)) + " ,但是run模式至少需要9个参数执行",
				})
				return
			}
		} else {
			show(Tester.TestInfo{
				Statu: "024",
				Info:  "参数3错误，应为:compile或run或spjrun，但是输入了：" + args[3],
			})
			return
		}
	}
	var tester Tester.Tester
	time, err := strconv.ParseInt(args[5], 10, 64)
	if err != nil {
		show(Tester.TestInfo{
			Statu: "001",
			Info:  "执行时间限制错误，输入为 " + args[5] + " ,但是期望的值为一个64位整数，单位ms。错误信息：" + err.Error(),
		})
		return
	}
	mem, err := strconv.ParseInt(args[6], 10, 64)
	if err != nil {
		show(Tester.TestInfo{
			Statu: "002",
			Info:  "执行时间限制错误，输入为 " + args[6] + " ,但是期望的值为一个64位整数，单位MB。错误信息：" + err.Error(),
		})
		return
	}
	mem *= 1024
	switch args[1] {
	case "gnu":
		tester = Tester.NewGUNTester(args[2], args[4], time, mem)
	case "java":
		tester = Tester.NewJAVATester(args[2], args[4], time, mem)
	default:
		show(Tester.TestInfo{
			Statu: "008",
			Info:  "语言类型输入错误，输入为：" + args[1] + " ，期望的值包括：" + types,
		})
		return
	}
	if args[3] == "compile" {
		show(tester.Compile())
	} else if args[3] == "run" {
		show(tester.Run(args[7]+".in", args[7]+".out"))
	}

}

func show(info Tester.TestInfo) {
	fmt.Println(info.ToStdString())
}
