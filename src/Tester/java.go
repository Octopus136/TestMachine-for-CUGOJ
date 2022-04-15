package Tester

//GUN编译文件

const (
	jdk   = "javac"
)
const versions = "'java'"

type JAVATester struct {
	version     string
	path        string
	out         string
	timeLimit   int64
	memoryLimit int64
}

func NewJAVATester(version, path string, timeLimit, memoryLimit int64) JAVATester {
	var tmp JAVATester
	tmp.version = version

	tmp.path = path + ".java"
	tmp.timeLimit = timeLimit
	tmp.memoryLimit = memoryLimit
	tmp.out = path

	return tmp
}

func (tester JAVATester) Compile() TestInfo {
	var cmd string
	if(tester.version == version){
		cmd = javac
	}
	else{
		return TestInfo{
			Statu:   "003",
			Info:    "编译器版本选择错误，输入为 " + tester.version + " ,但是期望的值只包括" + versions,
			RunTime: -1,
			Memory:  -1,
		}
	}
	cmdArgs := cmd + " -d " + tester.out + " " + tester.path
	return CompileBase(cmdArgs, tester.timeLimit, tester.memoryLimit)
}

func (tester JAVATester) Run(in, out string) TestInfo {
	runCmd := "java " + tester.out + ".class"
	return RunBase(runCmd, in, out, tester.timeLimit, tester.memoryLimit)
}

func (tester JAVATester) Spj(in, out, spj string) TestInfo {
	return TestInfo{}
}
