package Tester

import (
	"TestMachine/src/LimitExec"
	"bufio"
	"bytes"
	"encoding/json"

	//"fmt"
	"io"
	"os"
	//"strconv"
)

type TestBuffer struct {
	buf []byte
	siz int
	len int
	off int
	out bool
}

func NewTestBuffer(length int) TestBuffer {
	return TestBuffer{
		buf: make([]byte, length),
		len: 0,
		siz: length,
		off: 0,
		out: false,
	}
}

func (t *TestBuffer) Peek() (byte, error) {
	if t.off == t.len {
		return 0, LimitExec.ExecError{Info: "EOF"}
	}
	return t.buf[t.off], nil
}

func (t *TestBuffer) Move() {
	if t.off < t.len {
		t.off++
	}
}

func (t *TestBuffer) GetByte() (byte, error) {
	if t.off == t.len {
		return 0, LimitExec.ExecError{Info: "EOF"}
	}
	t.off++
	return t.buf[t.off-1], nil
}

func (t *TestBuffer) Read(p []byte) (n int, err error) {
	length := len(p)
	for i := 0; i < length; i++ {
		if t.off == t.len {
			if i == 0 {
				return 0, LimitExec.ExecError{Info: "EOF"}
			} else {
				return i, nil
			}
		}
		p[i] = t.buf[t.off]
		t.off++
	}
	return length, nil
}

func ReadLine(p io.Reader) ([]byte, error) {
	in1 := bufio.NewReader(p)
	return in1.ReadBytes('\n')
}

func CheckLine(p1, p2 []byte) bool {
	len1 := len(p1) - 1
	len2 := len(p2) - 1
	for len1 >= 0 && (p1[len1] == 0 || p1[len1] == '\n' || p1[len1] == ' ') {
		len1--
	}
	for len2 >= 0 && (p2[len2] == 0 || p2[len2] == '\n' || p2[len2] == ' ') {
		len2--
	}
	if len1 != len2 {
		return false
	}
	for i := 0; i <= len1; i++ {
		if p1[i] != p2[i] {
			return false
		}
	}
	return true
}

func (t *TestBuffer) Check(p io.Reader) (bool, error) {

	for {
		out := false
		p1, err := ReadLine(t)
		if err != nil {
			out = true
		}
		p2, err := ReadLine(p)
		if err != nil {
			out = true
		}
		if !CheckLine(p1, p2) {
			return false, nil
		}
		if out {
			break
		}
	}
	return true, nil
}

func (t *TestBuffer) Write(p []byte) (n int, err error) {
	plen := len(p)

	if t.len+plen > t.siz {
		t.out = true
		return 0, LimitExec.ExecError{Info: "OLE"}
	}
	tmp := t.buf[t.off:]
	copy(tmp, p)
	t.len += plen

	return plen, nil
}

type TestInfo struct {
	Statu   string
	Info    string
	RunTime int64
	Memory  int64
}

func (info TestInfo) ToBytes() []byte {
	tmp, err := json.Marshal(info)
	if err != nil {
		return []byte("{\"Statu\":\"005\",\"Info\":\"" + err.Error() + "\",\"RunTime\":-1,\"Memory\":-1}")
	}
	return tmp
}

func (info TestInfo) ToString() string {
	tmp := info.ToBytes()
	str := string(tmp)
	return str
}

func (info TestInfo) ToStdString() string {
	tmp := info.ToBytes()
	str := bytes.Buffer{}
	_ = json.Indent(&str, tmp, "", "	")
	return str.String()
}

type Tester interface {
	Compile() TestInfo
	Run(in, out string) TestInfo
	Spj(in, out, spj string) TestInfo
}

//
//
//
//
//
func CompileBase(cmd string, timeLimit, memoryLimit int64) TestInfo {
	stderr := bytes.Buffer{}
	timeUse, memoryUse, err := LimitExec.LimitExecE(timeLimit, memoryLimit, cmd, &stderr)
	var res TestInfo
	res.Statu = "006"
	if err != nil {
		if err.Error() == "MLE" {
			res.Info = "编译使用内存超限"
		} else if err.Error() == "TLE" {
			res.Info = "编译超时"
		} else if err.Error() == "RE" {
			res.Info = "编译出现错误" + stderr.String()
		} else {
			res.Info = err.Error() + stderr.String()
		}
		res.Statu = "007"
	} else {
		res.Info = "编译成功"
	}
	res.RunTime = timeUse
	res.Memory = memoryUse
	return res
}

func ReadFile(f io.Reader) *bytes.Buffer {
	resbuf := bytes.Buffer{}
	buf := make([]byte, 1024*1024)
	for {
		len, err := f.Read(buf)
		if len == 0 {
			break
		} else if err != nil {
			return nil
		} else {
			resbuf.Write(buf[0:len])
		}
	}
	return &resbuf
}

func RunBase(cmd, in, out string, timeLimit, memoryLimit int64) TestInfo {
	infile, err := os.OpenFile(in, os.O_RDONLY, 0444)
	if err != nil {
		return TestInfo{
			Statu: "009",
			Info:  "测试输入用例不存在：" + in,
		}
	}
	outfile, err := os.OpenFile(out, os.O_RDONLY, 0444)
	if err != nil {
		return TestInfo{
			Statu: "009",
			Info:  "测试输出用例不存在：" + out,
		}
	}

	stdoutbuf := ReadFile(outfile)
	if stdoutbuf == nil {
		return TestInfo{
			Statu: "016",
			Info:  "文件读取出错：" + out,
		}
	}
	stdout := NewTestBuffer(stdoutbuf.Len() + 1024*16)

	stderr := bytes.Buffer{}

	timeUse, memoryUse, err := LimitExec.LimitExec(timeLimit, memoryLimit, cmd, infile, &stdout, &stderr)

	res := TestInfo{}
	if err != nil {
		if err.Error() == "MLE" {
			res.Info = "内存超限"
			res.Statu = "013"
		} else if err.Error() == "TLE" {
			res.Info = "运行超时"
			res.Statu = "011"
		} else if err.Error() == "OLE" {
			res.Info = "输出超限"
			res.Statu = "015"
		} else {
			res.Statu = "012"
			res.Info = "运行错误"
		}
	} else {
		if stdout.out {
			res.Info = "输出超限"
			res.Statu = "015"
		} else {
			if ok, err := stdout.Check(stdoutbuf); err != nil {
				res.Info = "评测机内部错误"
				res.Statu = "017"
			} else if ok {
				res.Info = "运行通过"
				res.Statu = "010"
			} else {
				res.Info = "程序运行结果错误"
				res.Statu = "014"
			}
		}
	}
	res.RunTime = timeUse
	res.Memory = memoryUse
	return res
}
