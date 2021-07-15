# dog

watchdog for RSS/CPU

## 编译

- `go install -trimpath -ldflags='-extldflags=-static -s -w' ./...`
- `GOOS=linux GOARCH=amd64 go install -trimpath -ldflags='-extldflags=-static -s -w' ./...`
- `ldd /usr/local/bin/busy`

## 部署

1. `dog --init` 在当前目录创建 `ctl` 脚本和 示例 `dog.yml`
1. `./ctl start` 启动； `./ctl stop` 停止；`./ctl restart` 重新启动；`./ctl tail` 查看日志；

## demo

### 限制CPU

启动测试目标：

```sh
$ busy -p50             
2021/07/14 16:29:44 busy starting, pid 45198
2021/07/14 16:29:44  run 50% of 12/12 CPU cores forever
```

放狗，咬死超过50%：

```sh
$ GOLOG_STDOUT=true dog -max-pcpu 500 -filter busy -cond 2/30s -log ENV,CWD
2021-07-15 10:00:53.059 [INFO ] 625 --- [1    ] [-]  : log file created:~/logs/dog/dog.log
2021-07-15 10:00:53.060 [INFO ] 625 --- [1    ] [-]  : dog with config: &{Topn:0 Pid:0 Ppid:0 Self:false KillSignals:[INT KILL] Interval:10s MaxMem:0 MaxPmem:50 MaxPcpu:300 CmdFilter:[] LogItems:[ENV CWD] RateConfig:2/30s limiter:0xc000082f80 Jitter:1s} created
2021-07-15 10:00:56.515 [INFO ] 625 --- [1    ] [-]  : Dog barking for 3, config:2/30s, item User: bingoobjca Pid: 98283 Ppid: 66509 %cpu: 563 %mem: 0 VSZ: 5.1GB, RSS: 3.8MB Tty: s002 Stat: R+ Start: 2021-07-15 02:00:08 Time: 3:59.93 Command: busy -p50
2021-07-15 10:01:07.103 [INFO ] 625 --- [1    ] [-]  : Dog barking for 3, config:2/30s, item User: bingoobjca Pid: 98283 Ppid: 66509 %cpu: 599.8 %mem: 0 VSZ: 5.1GB, RSS: 3.8MB Tty: s002 Stat: R+ Start: 2021-07-15 02:00:08 Time: 4:57.69 Command: busy -p50
2021-07-15 10:01:17.721 [INFO ] 625 --- [1    ] [-]  : Dog barking for 3, config:2/30s, item User: bingoobjca Pid: 98283 Ppid: 66509 %cpu: 495.5 %mem: 0 VSZ: 5.1GB, RSS: 3.8MB Tty: s002 Stat: R+ Start: 2021-07-15 02:00:08 Time: 5:53.81 Command: busy -p50
2021-07-15 10:01:27.121 [INFO ] 625 --- [1    ] [-]  : Dog biting for 3, item User: bingoobjca Pid: 98283 Ppid: 66509 %cpu: 623.5 %mem: 0 VSZ: 5.1GB, RSS: 3.8MB Tty: s002 Stat: R+ Start: 2021-07-15 02:00:08 Time: 6:48.59 Command: busy -p50
2021-07-15 10:01:27.127 [INFO ] 625 --- [1    ] [-]  : LogItem: ENV, Value: 98283 s002  S+     7:11.22 busy -p50 PATH=/usr/local/go/bin:/Users/bingoobjca/go/bin:/usr/local/Cellar/mysql-client/8.0.23/bin:/Users/bingoobjca/go/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:/Applications/VMware Fusion.app/Contents/Public:/Library/TeX/texbin:/usr/local/go/bin:/usr/local/aria2/bin:/Users/bingoobjca/.cargo/bin:/Users/bingoobjca/.fzf/bin TERM=xterm-256color COMMAND_MODE=unix2003 __INTELLIJ_COMMAND_HISTFILE__=/Users/bingoobjca/Library/Caches/JetBrains/GoLand2021.1/terminal/history/dog-history1 LOGNAME=bingoobjca XPC_SERVICE_NAME=0 __CFBundleIdentifier=com.jetbrains.goland SHELL=/bin/zsh GOPATH=/Users/bingoobjca/go USER=bingoobjca GOROOT=/usr/local/go TMPDIR=/var/folders/c8/ft7qp47d6lj5579gmyflxbr80000gn/T/ TERMINAL_EMULATOR=JetBrains-JediTerm LOGIN_SHELL=1 GO111MODULE=on SSH_AUTH_SOCK=/private/tmp/com.apple.launchd.O4JlykSOLq/Listeners XPC_FLAGS=0x0 TERM_SESSION_ID=3beeb8a3-2d1f-479c-9165-2514faed7d26 __CF_USER_TEXT_ENCODING=0x1F5:0x19:0x34 LC_CTYPE=zh_CN.UTF-8 HOME=/Users/bingoobjca SHLVL=1 PWD=/Users/bingoobjca/github/dog OLDPWD=/Users/bingoobjca/github/dog ZSH=/Users/bingoobjca/.oh-my-zsh PAGER=less LESS=-R LSCOLORS=Gxfxcxdxbxegedabagacad http_proxy=http://127.0.0.1:9999 HTTP_PROXY=http://127.0.0.1:9999 https_proxy=http://127.0.0.1:9999 HTTPS_PROXY=http://127.0.0.1:9999 all_proxy=http://127.0.0.1:10000 ALL_PROXY=http://127.0.0.1:10000 HSTR_CONFIG=hicolor _=/Users/bingoobjca/go/bin/busy
2021-07-15 10:01:27.144 [INFO ] 625 --- [1    ] [-]  : LogItem: CWD, Value: /Users/bingoobjca/github/dog
2021-07-15 10:01:27.144 [INFO ] 625 --- [1    ] [-]  : Kill interrupt to 98283 succeeded
2021-07-15 10:01:27.144 [INFO ] 625 --- [1    ] [-]  : Kill killed to 98283 succeeded
```

测试目标打印：

```sh
2021/07/14 16:32:06 received signal interrupt, exiting
```

## help

```sh
$ dog -h                                      
Usage of dog:
  -filter value 命令包含，以!开头为不包含，可以多个值
  -cond string 发送条件，默认触发1次就发信号，eg.3/30s，在30s内发生3次，则触发 
  -kill string 发送信号，多个逗号分隔，eg. INT,TERM,KILL,QUIT,USR1,USR2 (默认 INT)
  -log  string 记录日志信息，多个逗号分隔，eg. ENV,CWD
  -max-mem value 允许最大内存 (默认 0B，不检查内存)
  -max-pcpu int 允许内存最大百分比, eg. 1-1200 (默认 600), 0 不查 CPU
  -max-pmem int 允许CPU最大百分比, eg. 1-100 (默认 50)
  -pid int 指定pid
  -ppid int 指定ppid
  -self 是否监控自身
  -span duration 检查时间间隔 (默认 10s)
  -topn int 只取前N个检查
  -v Print version info and exit
$ busy -h               
Usage of busy:
  -c int 使用核数，默认 12
  -d duration 跑多久，默认一直跑
  -m string 总内存,增量, eg. 1) 10M 直接达到10M 2) 10M,1K/10s 总用量10M,每10秒增加1K
  -p int 每核CPU百分比 (默认 100)
```

