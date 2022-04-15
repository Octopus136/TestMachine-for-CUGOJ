package Tester

//GUN编译文件

const (
	c99   = "-std=c99"
	c11   = "-std=c11"
	cpp11 = "-std=c++11"
	cpp14 = "-std=c++14"
	cpp17 = "-std=c++17"
	cpp20 = "-std=c++2a"
)
const versions = "'c99','c11','cpp11','cpp14','cpp17','cpp20'"

type GUNTester struct {
	version     string
	path        string
	out         string
	timeLimit   int64
	memoryLimit int64
}

func NewGUNTester(version, path string, timeLimit, memoryLimit int64) GUNTester {
	var tmp GUNTester
	tmp.version = version

	if version == "c99" || version == "c11" {
		tmp.path = path + ".c"
	} else {
		tmp.path = path + ".cpp"
	}
	tmp.timeLimit = timeLimit
	tmp.memoryLimit = memoryLimit
	tmp.out = path

	return tmp
}

func (tester GUNTester) Compile() TestInfo {
	var cmd string
	var std string
	switch tester.version {
	case "c99":
		cmd = "gcc"
		std = c99
	case "c11":
		cmd = "gcc"
		std = c11
	case "cpp11":
		cmd = "g++"
		std = cpp11
	case "cpp14":
		cmd = "g++"
		std = cpp14
	case "cpp17":
		cmd = "g++"
		std = cpp17
	case "cpp20":
		cmd = "g++"
		std = cpp20
	default:
		return TestInfo{
			Statu:   "003",
			Info:    "编译器版本选择错误，输入为 " + tester.version + " ,但是期望的值只包括" + versions,
			RunTime: -1,
			Memory:  -1,
		}
	}
	cmdArgs := cmd + " -O2 " + std + " " + tester.path + " -o " + tester.out
	return CompileBase(cmdArgs, tester.timeLimit, tester.memoryLimit)
}

func (tester GUNTester) Run(in, out string) TestInfo {
	return RunBase(tester.out, in, out, tester.timeLimit, tester.memoryLimit)
}

func (tester GUNTester) Spj(in, out, spj string) TestInfo {
	return TestInfo{}
}
