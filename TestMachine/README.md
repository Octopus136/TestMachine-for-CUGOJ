# TestMachine
## 介绍
TestMachine是CUGOJ所依赖的评测核心，负责执行评测工作，包括：编译、普通评测、Special Judge、交互式评测等。
## 安装
TestMachine使用Golang实现，由于需要实现多环境评测，因此有较多的依赖需要进行处理。
这里提供两种安装方式：
*  **容器**
   
   由于TestMachine最终会以容器镜像的形式进行发布，因此用户可以直接下载镜像使用OCI容器运行时运行。
* ### **Ubuntu**
  
   * 提供评测机运行环境
 
      安装Golang：
      ~~~bash
      # 将文件下载到local文件夹下，注意中间的版本可以在官网自行选择，不要低于1.18
      cd /usr/local
      sudo wget https://golang.google.cn/dl/go1.18.linux-amd64.tar.gz
      # 解压文件到/usr/local，解压完成后可以执行以下命令检查文件内容
      # ls /usr/local/go
      sudo tar xfz go1.18.linux-amd64.tar.gz -C /usr/local
      # 删除下载的包
      sudo rm go1.18.linux-amd64.tar.gz
      # 编辑配置文件，配置环境变量
      vim /etc/profile
      # 按i开始编辑，在文件末尾中加入以下内容：
      # export GOROOT=/usr/local/go
      # export GOPATH=$HOME/gowork
      # export GOBIN=$GOPATH/bin
      # export PATH=$GOPATH:$GOBIN:$GOROOT/bin:$PATH
      # 添加完成后ctrl+c退出编辑模式，ctrl+: 输入指令 wq 回车保存退出
      
      # 使环境变量生效
      source /etc/profile

      # 持久化环境变量，避免下次重启失效
      cd ~
      vim .bashrc
      # 在末尾添加以下内容：
      # source /etc/profile

      # 检查Golang是否安装成功
      go version
      ~~~
  * 提供C/C++编译环境
  
    安装GUN：
    ~~~bash
    # 安装GNU编译环境
    sudo apt install build-essential

    # 测试是否安装成功
    gcc -v
    ~~~
* **其他Linux版本**
    
    TODO
## 使用方法
TestMachine主要负责进行最基础的编译、评测工作，目前已实现C/C++的编译以及常规测试

* 参数列表：

  args[1] 编译器类型(gun)

  args[2] 语言版本(
    gun:c99,c11,cpp11,cpp14,cpp17,cpp20
  )

  args[3] 执行方式(compile,run,spjrun)
  args[4] 源文件、可执行文件路径，不包含后缀，例如/code/main，不需要/code/main.c

  args[5] 执行时间限制，单位ms

  args[6] 执行空间限制，单位KB

  args[7] 测试数据路径，不包含后缀，默认后缀分别为.in和.out，如果是spj，允许.out文件不存在

  args[8] spj路径
  (compile方式要求有7个参数
  ;run方式要求有8个参数
  ;spj run方式要求9个参数)

* C/C++编译
    
    ~~~bash
    # ...请填写自己的安装位置
    cd .../TestMachine/bin
    # 使用GNU编译器，以C99为测试语言，以编译方式运行，编译文件位置为../test/main (对于C99、C11语言会寻找文件main.c，对于CPP会寻找main.cpp) 编译时限10000ms，内存限制262144KB
    ./cugtm gnu c99 compile ../test/main 10000 262144
    ~~~
* C/C++运行

    ~~~bash
    # ...请填写自己的安装位置
    cd .../TestMachine/bin
    # 使用GNU编译器，以C99为测试语言，以常规评测方式运行，评测文件位置为../test/main 评测时限10000ms，内存限制262144KB 测试文件名为../test/test1，系统将会读取../test/test1.in为输入文件，../test/test1.out为标准答案
    ./cugtm gnu c99 run ../test/main 10000 262144 ../test/test1
    ~~~

## 文件清单
```
TestMachine
├─ LICENSE
├─ README.md
├─ bin
│  └─ cugtm
├─ doc
│  └─ ErrorCodes.json
├─ go.mod
├─ src
│  ├─ LimitExec
│  │  └─ LimitExec.go
│  ├─ Tester
│  │  ├─ Tester.go
│  │  └─ gnu.go
│  └─ main.go
└─ test
   ├─ main.c
   ├─ test1.in
   └─ test1.out

bin/cugtm 为Linux可执行文件
test目录包含编译、评测用测试文件
src目录包含项目源代码
```
## 更新
  * 4.7 初次发布，实现了评测框架以及针对C/C++的基本评测和编译功能
