# dog

watchdog for RSS/CPU

## 特性

1. 配置主机最小可用内存大小，触发时，驱逐内存排名第一的进程（白名单除外）
1. 配置主机最大 CPU 百分比，触发时，驱逐CPU排名第一的进程（白名单除外）
1. 配置单个进程最大占用内存，单个进程最大内存占比，单个进程最大CPU占比，超过时，驱逐

## 编译

- 本机编译：`make`
- 交叉编译Linux：`make linux`

## 部署

1. `dogwatch --init` 在当前目录创建 `ctl` 启停 shell 脚本和 示例配置文件 `dogwatch.yml`
1. `./ctl start` 后台启动； `./ctl stop` 停止；`./ctl restart` 重新启动；`./ctl tail` 查看日志；

## demo

### 限制CPU

启动测试目标：

```sh
$ dogbusy -p50             
2021/07/14 16:29:44 busy starting, pid 45198
2021/07/14 16:29:44  run 50% of 12/12 CPU cores forever
```

放狗，咬死超过50%：

```sh
$ GOLOG_STDOUT=true dogwatch -max-pcpu 500 -filter busy -cond 2/30s -log ENV,CWD
2021-07-15 10:00:53.059 [INFO ] 625 --- [1    ] [-]  : log file created:~/logs/dog/dog.log
2021-07-15 10:00:53.060 [INFO ] 625 --- [1    ] [-]  : dog with config: &{Topn:0 Pid:0 Ppid:0 Self:false KillSignals:[INT KILL] Interval:10s MaxMem:0 MaxPmem:50 MaxPcpu:300 CmdFilter:[] LogItems:[ENV CWD] RateConfig:2/30s limiter:0xc000082f80 Jitter:1s} created
2021-07-15 10:00:56.515 [INFO ] 625 --- [1    ] [-]  : Dog barking for 3, config:2/30s, item User: bingoo Pid: 98283 Ppid: 66509 %cpu: 563 %mem: 0 VSZ: 5.1GB, RSS: 3.8MB Tty: s002 Stat: R+ Start: 2021-07-15 02:00:08 Time: 3:59.93 Command: busy -p50
2021-07-15 10:01:07.103 [INFO ] 625 --- [1    ] [-]  : Dog barking for 3, config:2/30s, item User: bingoo Pid: 98283 Ppid: 66509 %cpu: 599.8 %mem: 0 VSZ: 5.1GB, RSS: 3.8MB Tty: s002 Stat: R+ Start: 2021-07-15 02:00:08 Time: 4:57.69 Command: busy -p50
2021-07-15 10:01:17.721 [INFO ] 625 --- [1    ] [-]  : Dog barking for 3, config:2/30s, item User: bingoo Pid: 98283 Ppid: 66509 %cpu: 495.5 %mem: 0 VSZ: 5.1GB, RSS: 3.8MB Tty: s002 Stat: R+ Start: 2021-07-15 02:00:08 Time: 5:53.81 Command: busy -p50
2021-07-15 10:01:27.121 [INFO ] 625 --- [1    ] [-]  : Dog biting for 3, item User: bingoo Pid: 98283 Ppid: 66509 %cpu: 623.5 %mem: 0 VSZ: 5.1GB, RSS: 3.8MB Tty: s002 Stat: R+ Start: 2021-07-15 02:00:08 Time: 6:48.59 Command: busy -p50
2021-07-15 10:01:27.127 [INFO ] 625 --- [1    ] [-]  : LogItem: ENV, Value: 98283 s002  S+     7:11.22 busy -p50 PATH=/usr/local/go/bin:/Users/bingoo/go/bin:/usr/local/Cellar/mysql-client/8.0.23/bin:/Users/bingoo/go/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:/Applications/VMware Fusion.app/Contents/Public:/Library/TeX/texbin:/usr/local/go/bin:/usr/local/aria2/bin:/Users/bingoo/.cargo/bin:/Users/bingoo/.fzf/bin TERM=xterm-256color COMMAND_MODE=unix2003 _INTELLIJ_COMMAND_HISTFILE_=/Users/bingoo/Library/Caches/JetBrains/GoLand2021.1/terminal/history/dog-history1 LOGNAME=bingoo XPC_SERVICE_NAME=0 _CFBundleIdentifier=com.jetbrains.goland SHELL=/bin/zsh GOPATH=/Users/bingoo/go USER=bingoo GOROOT=/usr/local/go TMPDIR=/var/folders/c8/ft7qp47d6lj5579gmyflxbr80000gn/T/ TERMINAL_EMULATOR=JetBrains-JediTerm LOGIN_SHELL=1 GO111MODULE=on SSH_AUTH_SOCK=/private/tmp/com.apple.launchd.O4JlykSOLq/Listeners XPC_FLAGS=0x0 TERM_SESSION_ID=3beeb8a3-2d1f-479c-9165-2514faed7d26 _CF_USER_TEXT_ENCODING=0x1F5:0x19:0x34 LC_CTYPE=zh_CN.UTF-8 HOME=/Users/bingoo SHLVL=1 PWD=/Users/bingoo/github/dog OLDPWD=/Users/bingoo/github/dog ZSH=/Users/bingoo/.oh-my-zsh PAGER=less LESS=-R LSCOLORS=Gxfxcxdxbxegedabagacad http_proxy=http://127.0.0.1:9999 HTTP_PROXY=http://127.0.0.1:9999 https_proxy=http://127.0.0.1:9999 HTTPS_PROXY=http://127.0.0.1:9999 all_proxy=http://127.0.0.1:10000 ALL_PROXY=http://127.0.0.1:10000 HSTR_CONFIG=hicolor _=/Users/bingoo/go/bin/busy
2021-07-15 10:01:27.144 [INFO ] 625 --- [1    ] [-]  : LogItem: CWD, Value: /Users/bingoo/github/dog
2021-07-15 10:01:27.144 [INFO ] 625 --- [1    ] [-]  : Kill interrupt to 98283 succeeded
2021-07-15 10:01:27.144 [INFO ] 625 --- [1    ] [-]  : Kill killed to 98283 succeeded
```

测试目标打印：

```sh
2021/07/14 16:32:06 received signal interrupt, exiting
```

## 运行时长

启动测试目标：

```sh
$ SHORT_TASK=true dogbusy -p20
2021/07/15 10:27:56 busy starting, pid 63610
2021/07/15 10:27:56  run 20% of 12/12 CPU cores forever.
2021/07/15 10:29:09 received signal interrupt, exiting
```

放狗，咬死超过10s：

```sh
$ GOLOG_STDOUT=true dogwatch -max-time 10s -max-time-env SHORT_TASK -cond 2/30s -log ENV,CWD
2021-07-15 10:28:11.375 [INFO ] 64264 --- [1    ] [-]  : log file created:~/logs/dog/dog.log
2021-07-15 10:28:11.376 [INFO ] 64264 --- [1    ] [-]  : dog with config: &{Topn:0 Pid:0 Ppid:0 Self:false KillSignals:[INT KILL] Interval:10s MaxMem:0 MaxPmem:50 MaxPcpu:300 CmdFilter:[] LogItems:[ENV CWD] RateConfig:2/30s limiter:0xc0001100c0 Jitter:1s MaxTime:10s MaxTimeEnv:SHORT_TASK} created
2021-07-15 10:28:19.048 [INFO ] 64264 --- [1    ] [-]  : Dog barking for 运行时长超了, config:2/30s, item User: bingoo Pid: 63610 Ppid: 66509 %cpu: 203.6 %mem: 0 VSZ: 5.1GB, RSS: 3.7MB Tty: s002 Stat: S+ Start: 2021-07-15 02:27:e: 0:33.48 Command: busy -p20
2021-07-15 10:28:29.673 [INFO ] 64264 --- [1    ] [-]  : Dog barking for 运行时长超了, config:2/30s, item User: bingoo Pid: 63610 Ppid: 66509 %cpu: 230.4 %mem: 0 VSZ: 5.1GB, RSS: 3.7MB Tty: s002 Stat: S+ Start: 2021-07-15 02:27:e: 0:54.87 Command: busy -p20
2021-07-15 10:28:38.654 [INFO ] 64264 --- [1    ] [-]  : Dog barking for 运行时长超了, config:2/30s, item User: bingoo Pid: 63610 Ppid: 66509 %cpu: 208.1 %mem: 0 VSZ: 5.1GB, RSS: 3.7MB Tty: s002 Stat: S+ Start: 2021-07-15 02:27:e: 1:16.40 Command: busy -p20
2021-07-15 10:28:49.925 [INFO ] 64264 --- [1    ] [-]  : Dog barking for 运行时长超了, config:2/30s, item User: bingoo Pid: 63610 Ppid: 66509 %cpu: 206.8 %mem: 0 VSZ: 5.1GB, RSS: 3.7MB Tty: s002 Stat: R+ Start: 2021-07-15 02:27:e: 1:39.28 Command: busy -p20
2021-07-15 10:29:00.074 [INFO ] 64264 --- [1    ] [-]  : Dog barking for 运行时长超了, config:2/30s, item User: bingoo Pid: 63610 Ppid: 66509 %cpu: 232.4 %mem: 0 VSZ: 5.1GB, RSS: 3.8MB Tty: s002 Stat: R+ Start: 2021-07-15 02:27:e: 2:00.14 Command: busy -p20
2021-07-15 10:29:09.495 [INFO ] 64264 --- [1    ] [-]  : Dog biting for 运行时长超了, item User: bingoo Pid: 63610 Ppid: 66509 %cpu: 239.6 %mem: 0 VSZ: 5.1GB, RSS: 3.8MB Tty: s002 Stat: S+ Start: 2021-07-15 02:27:56 Time: 2:21.0and: busy -p20
2021-07-15 10:29:09.509 [INFO ] 64264 --- [1    ] [-]  : LogItem: ENV, Value: 63610 s002  S+     2:38.94 busy -p20 PATH=/usr/local/go/bin:/Users/bingoo/go/bin:/usr/local/Cellar/mysql-client/8.0.23/bin:/Users/bingoo/go/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:/Applications/VMware Fusion.app/Contents/Public:/Library/TeX/texbin:/usr/local/go/bin:/usr/local/aria2/bin:/Users/bingoo/.cargo/bin:/Users/bingoo/.fzf/bin TERM=xterm-256color COMMAND_MODE=unix2003 _INTELLIJ_COMMAND_HISTFILE_=/Users/bingoo/Library/Caches/JetBrains/GoLand2021.1/terminal/history/dog-history1 LOGNAME=bingoo XPC_SERVICE_NAME=0 _CFBundleIdentifier=com.jetbrains.goland SHELL=/bin/zsh GOPATH=/Users/bingoo/go USER=bingoo GOROOT=/usr/local/go TMPDIR=/var/folders/c8/ft7qp47d6lj5579gmyflxbr80000gn/T/ TERMINAL_EMULATOR=JetBrains-JediTerm LOGIN_SHELL=1 GO111MODULE=on SSH_AUTH_SOCK=/private/tmp/com.apple.launchd.O4JlykSOLq/Listeners XPC_FLAGS=0x0 TERM_SESSION_ID=3beeb8a3-2d1f-479c-9165-2514faed7d26 _CF_USER_TEXT_ENCODING=0x1F5:0x19:0x34 LC_CTYPE=zh_CN.UTF-8 HOME=/Users/bingoo SHLVL=1 PWD=/Users/bingoo/github/dog OLDPWD=/Users/bingoo/github/dog ZSH=/Users/bingoo/.oh-my-zsh PAGER=less LESS=-R LSCOLORS=Gxfxcxdxbxegedabagacad http_proxy=http://127.0.0.1:9999 HTTP_PROXY=http://127.0.0.1:9999 https_proxy=http://127.0.0.1:9999 HTTPS_PROXY=http://127.0.0.1:9999 all_proxy=http://127.0.0.1:10000 ALL_PROXY=http://127.0.0.1:10000 HSTR_CONFIG=hicolor SHORT_TASK=true _=/Users/bingoo/go/bin/busy
2021-07-15 10:29:09.529 [INFO ] 64264 --- [1    ] [-]  : LogItem: CWD, Value: /Users/bingoo/github/dog
2021-07-15 10:29:09.529 [INFO ] 64264 --- [1    ] [-]  : Kill interrupt to 63610 succeeded
2021-07-15 10:29:09.529 [INFO ] 64264 --- [1    ] [-]  : Kill killed to 63610 succeeded
```

## help

```sh
$ dogwatch -h
Usage of dogwatch:
  -filter value 命令包含，以!开头为不包含，可以多个值
  -cond string 发送条件，默认触发1次就发信号，eg.3/30s，在30s内发生3次，则触发 
  -kill string 发送信号，多个逗号分隔，eg. INT,TERM,KILL,QUIT,USR1,USR2 (默认 INT)
  -log  string 记录日志信息，多个逗号分隔，eg. ENV,CWD
  -max-time value 允许最大启动时长 (默认 0，不检查启动时长)
  -max-time-env value 允许最大启动时长包含的环境变量
  -max-mem value 允许最大内存 (默认 0B，不检查内存)
  -max-pcpu int 允许内存最大百分比, eg. 1-1200 (默认 600), 0 不查 CPU
  -max-pmem int 允许CPU最大百分比, eg. 1-100 (默认 50)
  -min-available-memory 允许最小总可用内存 (默认 0B，不检查此项)
  -max-host-cpu 允许最大机器CPU百分比（0-100） (默认 0，不检查此项)
  -whites value 总最小内存/最大机器CPU百分比触发时，驱逐进程命令行包含白名单，可以多个值
  -pid int 指定pid
  -ppid int 指定ppid
  -self 是否监控自身
  -span duration 检查时间间隔 (默认 10s)
  -jitter duration 最大抖动 (默认 1s)
  -topn int 只取前N个检查
  -v Print version info and exit

$ dogbusy -h               
Usage of dogbusy:
  -c int      使用核数，默认 12
  -p int      每核 CPU 百分比 (默认 100), 0 时不开启 CPU 耗用
  -l          是否在 CPU 耗用时锁定 OS 线程
  -m string   总内存耗用，默认不开启, e.g. 1) 10M 直达10M 2) 10M,1K/10s 总10M,每10秒加1K
  -d duration 跑多久，默认一直跑，e.g. 10s 20m 30h
  -v          看下版本号
```

## 检查狗咬日志

```sh
$ bssh -H q2,q11,q15,q14,q7,q8 "watchdog/ctl log"
Select Server :q2,q11,q15,q14,q7,q8
Run Command   :watchdog/ctl log
q2  ::  2021-07-22 12:18:14.716 [INFO ] 4396 --- [1    ] [-]  : Dog barking for CPU占比超了, config:2/30s, item User: root Pid: 24440 Ppid: 24405 %cpu: 273 %mem: 1.6 VSZ: 4GB, RSS: 543.2MB Tty: ? Stat: Sl Start: 2021-07-22 12:18:07 Time: 00:00:16 Command: | \_ java -server -Xmx768m -Xms768m -Xmn384m -XX:PermSize=128m -Xss256k -XX:+DisableExplicitGC -XX:+UseConcMarkSweepGC -XX:+CMSParallelRemarkEnabled -XX:+UseCMSCompactAtFullCollection -XX:LargePageSizeInBytes=128m -XX:+UseFastAccessorMethods -XX:+UseCMSInitiatingOccupancyOnly -XX:CMSInitiatingOccupancyFraction=70 -Dlogpath.base=/root/logs/ids-mini -cp /app/ids-mini/config:/app/ids-mini/lib/* -jar /app/ids-mini/lib/IDS-Mini-server-1.0.0-SNAPSHOT.jar ids-mini
q2  ::  2021-07-22 12:18:24.179 [INFO ] 4396 --- [1    ] [-]  : Dog biting for CPU占比超了, item User: root Pid: 24440 Ppid: 24405 %cpu: 255 %mem: 2.1 VSZ: 4.5GB, RSS: 719.6MB Tty: ? Stat: Sl Start: 2021-07-22 12:18:07 Time: 00:00:40 Command: | \_ java -server -Xmx768m -Xms768m -Xmn384m -XX:PermSize=128m -Xss256k -XX:+DisableExplicitGC -XX:+UseConcMarkSweepGC -XX:+CMSParallelRemarkEnabled -XX:+UseCMSCompactAtFullCollection -XX:LargePageSizeInBytes=128m -XX:+UseFastAccessorMethods -XX:+UseCMSInitiatingOccupancyOnly -XX:CMSInitiatingOccupancyFraction=70 -Dlogpath.base=/root/logs/ids-mini -cp /app/ids-mini/config:/app/ids-mini/lib/* -jar /app/ids-mini/lib/IDS-Mini-server-1.0.0-SNAPSHOT.jar ids-mini
q2  ::  2021-07-22 12:18:24.209 [INFO ] 4396 --- [1    ] [-]  : LogItem: ENV, Value: 24440 ?        Sl     0:40 java -server -Xmx768m -Xms768m -Xmn384m -XX:PermSize=128m -Xss256k -XX:+DisableExplicitGC -XX:+UseConcMarkSweepGC -XX:+CMSParallelRemarkEnabled -XX:+UseCMSCompactAtFullCollection -XX:LargePageSizeInBytes=128m -XX:+UseFastAccessorMethods -XX:+UseCMSInitiatingOccupancyOnly -XX:CMSInitiatingOccupancyFraction=70 -Dlogpath.base=/root/logs/ids-mini -cp /app/ids-mini/config:/app/ids-mini/lib/* -jar /app/ids-mini/lib/IDS-Mini-server-1.0.0-SNAPSHOT.jar ids-mini IDSBP_SERVICE_PORT_10620_TCP_PROTO=tcp IDS_SS_PORT_10640_TCP_PORT=10640 IDSBP_SZ_SS_PORT_HTTP=10622 IDS_MSS_PORT=10650 HOSTNAME=ids-mini-8675f586dc-mqzk2 IDS_QSS_PORT=10660 IDS_QS_PORT_10660_TCP=tcp://10.4.2.1:10660 IDSBP_DEMO_SS_PORT=10621 IDSPV2_SS_PORT=10610 IDS_SERVICE_PORT_10630_TCP_ADDR=127.0.4.1 AS_PORT_9611_TCP_PORT=9611 IDSBP_SZ_SS_PORT=10622 IDS_SERVICE_PORT_10630_TCP_PROTO=tcp IDS_CONSUMER_SS_HOST=127.0.2.1 KUBERNETES_PORT=tcp://127.0.0.1:443 KUBERNETES_PORT_443_TCP_PORT=443 IDSBP_DEMO_SERVICE_PORT_10621_TCP=tcp://127.0.8.3:10621 enableDataPersist= IDSBP_SERVICE_PORT_10620_TCP_PORT=10620 IDS_QS_PORT_10660_TCP_PROTO=tcp IDSPV2_DEMO_SERVICE_PORT_10601_TCP=tcp://127.0.9.2:10601 umpBaseDir=/var/log/footstone/metrics KUBERNETES_SERVICE_PORT=443 AS_PORT_9611_TCP=tcp://127.0.227.224:9611 OLDPWD=/app/ids-mini runOnBackground=false IDSPV2_SERVICE_PORT=tcp://10.4.1.0:10610 IDS_MINI_SERVICE_PORT_10650_TCP_PROTO=tcp IDSPV2_DEMO_SERVICE_PORT_10601_TCP_PORT=10601 KUBERNETES_SERVICE_HOST=127.0.0.1 IDS_SERVICE_PORT=tcp://127.0.4.1:10630 IDSPV2_SERVICE_PORT_10610_TCP_ADDR=10.4.1.0 IDS_CONSUMER_SERVICE_PORT_10401_TCP_PORT=10401 IDSBP_SS_PORT=10620 IDSBP_SERVICE_PORT=tcp://10.4.8.21:10620 IDS_QS_PORT_10660_TCP_PORT=10660 IDS_SERVICE_PORT_10400_TCP=tcp://1.4.2.2:10400 IDS_SERVICE_PORT=tcp://1.4.2.2:10400 IDSBP_DEMO_SERVICE_PORT_10621_TCP_PORT=10621 AUTHCODE_SERVICE_PORT=tcp://127.0.123.155:10007 IDSPV2_DEMO_SS_PORT_HTTP=10601 IS_PORT_10600_TCP_PORT=10600 IDS_SS_SERVICE_HOST=127.0.2.2 AUTHCODE_SERVICE_PORT_10007_TCP_PORT=10007 IDS_CONSUMER_SERVICE_PORT_10401_TCP=tcp://127.0.2.1:10401 AUTHCODE_SERVICE_PORT_10007_TCP_PROTO=tcp IDSPV2_SERVICE_PORT_10610_TCP_PORT=10610 IS_PORT_10600_TCP_ADDR=127.0.1.2 IDS_SS_HOST=127.0.4.1 IDS_SS_PORT_10640_TCP=tcp://127.0.2.2:10640 IDSPV2_DEMO_SERVICE_PORT_10601_TCP_PROTO=tcp IS_PORT=tcp://127.0.1.2:10600 IDS_CONSUMER_SERVICE_PORT_10401_TCP_PROTO=tcp IDS_SS_PORT_HTTP=10400 IDS_MINI_SERVICE_PORT_10650_TCP=tcp://10.4.18.21:10650 SERVER_ACCESSABLE_PORT=10000 IDS_SS_PORT=10400 IDS_MSS_PORT_HTTP=10650 ASS_HOST=127.0.227.224 IS_SERVICE_PORT=10600 IDSBP_SS_PORT_HTTP=10620 ASS_HOST=127.0.123.155 SERVER_ENVIROMENT= IDSBP_SZ_SERVICE_PORT_10622_TCP_ADDR=127.0.67.26 PATH=/opt/jdk/default/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin ASS_PORT_HTTP=9611 IDS_CONSUMER_SERVICE_PORT=tcp://127.0.2.1:10401 IDS_MINI_SERVICE_PORT=tcp://10.4.18.21:10650 MY_POD_NAME=ids-mini-8675f586dc-mqzk2 IDS_SS_PORT_10640_TCP_ADDR=127.0.2.2 IDSBP_DEMO_SERVICE_PORT_10621_TCP_PROTO=tcp IS_SERVICE_PORT_HTTP=10600 IDSPV2_SERVICE_PORT_10610_TCP=tcp://10.4.1.0:10610 PWD=/app/ids-mini IDSBP_SZ_SS_HOST=127.0.67.26 AS_PORT=tcp://127.0.227.224:9611 JAVA_HOME=/opt/jdk/default IDS_QS_PORT=tcp://10.4.2.1:10660 IDSBP_SERVICE_PORT_10620_TCP_ADDR=10.4.8.21 IDSBP_DEMO_SERVICE_PORT=tcp://127.0.8.3:10621 IDS_MINI_SERVICE_PORT_10650_TCP_PORT=10650 IDS_SERVICE_PORT_10400_TCP_PROTO=tcp IDSBP_SZ_SERVICE_PORT_10622_TCP=tcp://127.0.67.26:10622 AS_PORT_9611_TCP_ADDR=127.0.227.224 IDS_CONSUMER_SS_PORT_HTTP=10401 IDS_SERVICE_PORT_10630_TCP_PORT=10630 LANG=zh_CN.UTF-8 AS_PORT_9611_TCP_PROTO=tcp IDSPV2_DEMO_SERVICE_PORT=tcp://127.0.9.2:10601 applicationDir=/app/ids-mini IDSBP_SZ_SERVICE_PORT_10622_TCP_PROTO=tcp serviceName=ids-mini IDS_SS_PORT=10630 IDSBP_SZ_SERVICE_PORT=tcp://127.0.67.26:10622 IDS_SERVICE_PORT_10630_TCP=tcp://127.0.4.1:10630 IDSPV2_DEMO_SS_PORT=10601 jdkDir=jdk1.8.0_191 IDSBP_SERVICE_PORT_10620_TCP=tcp://10.4.8.21:10620 ASS_PORT=10007 IDSPV2_DEMO_SS_HOST=127.0.9.2 HOME=/root SHLVL=2 IDS_SERVICE_PORT_10400_TCP_ADDR=1.4.2.2 IDSPV2_SS_HOST=10.4.1.0 IDSBP_SZ_SERVICE_PORT_10622_TCP_PORT=10622 KUBERNETES_PORT_443_TCP_PROTO=tcp IDS_MINI_SERVICE_PORT_10650_TCP_ADDR=10.4.18.21 IDSBP_DEMO_SERVICE_PORT_10621_TCP_ADDR=127.0.8.3 KUBERNETES_SERVICE_PORT_HTTPS=443 IDS_SS_PORT=tcp://127.0.2.2:10640 IDS_SS_PORT_10640_TCP_PROTO=tcp SERVER_ACCESSABLE_IP= IDS_CONSUMER_SERVICE_PORT_10401_TCP_ADDR=127.0.2.1 IDS_SS_HOST=1.4.2.2 IDSPV2_SS_PORT_HTTP=10610 IDSBP_SS_HOST=10.4.8.21 IDS_SS_PORT_HTTP=10630 AUTHCODE_SERVICE_PORT_10007_TCP_ADDR=127.0.123.155 IDS_MSS_HOST=10.4.18.21 IDSPV2_DEMO_SERVICE_PORT_10601_TCP_ADDR=127.0.9.2 CLASSPATH=/opt/jdk/default/lib/*.jar:/opt/jdk/default/jre/lib/*.jar env=TEST ASS_PORT=9611 IDS_SS_SERVICE_PORT_HTTP=10640 IS_PORT_10600_TCP=tcp://127.0.1.2:10600 IDSBP_DEMO_SS_HOST=127.0.8.3 IDS_SS_SERVICE_PORT=10640 IDS_QSS_HOST=10.4.2.1 ASS_PORT_HTTP=10007 applicationLogDir=/root/logs/IDS-Mini enableUMP=true IDS_CONSUMER_SS_PORT=10401 AUTHCODE_SERVICE_PORT_10007_TCP=tcp://127.0.123.155:10007 KUBERNETES_PORT_443_TCP_ADDR=127.0.0.1 IDS_SERVICE_PORT_10400_TCP_PORT=10400 IDS_QS_PORT_10660_TCP_ADDR=10.4.2.1 IS_PORT_10600_TCP_PROTO=tcp KUBERNETES_PORT_443_TCP=tcp://127.0.0.1:443 IS_SERVICE_HOST=127.0.1.2 IDSBP_DEMO_SS_PORT_HTTP=10621 IDS_QSS_PORT_HTTP=10660 IDSPV2_SERVICE_PORT_10610_TCP_PROTO=tcp enableLogPersist=false _=/opt/jdk/default/bin/java
q2  ::  2021-07-22 12:18:24.333 [INFO ] 4396 --- [1    ] [-]  : LogItem: CWD, Value: /app/ids-mini
q2  ::  2021-07-22 12:18:24.333 [INFO ] 4396 --- [1    ] [-]  : Kill interrupt to 24440 succeeded
q2  ::  2021-07-22 12:18:24.334 [INFO ] 4396 --- [1    ] [-]  : Kill killed to 24440 succeeded
q2  ::  2021-07-22 12:18:53.573 [INFO ] 4396 --- [1    ] [-]  : Dog barking for cpu percent: 93.000000 > config max: 90.000000
q2  ::  2021-07-22 12:18:54.650 [INFO ] 4396 --- [1    ] [-]  : Dog barking for CPU占比超了, config:2/30s, item User: root Pid: 26500 Ppid: 26450 %cpu: 281 %mem: 1.9 VSZ: 4.3GB, RSS: 646.8MB Tty: ? Stat: Sl Start: 2021-07-22 12:18:41 Time: 00:00:33 Command: | \_ java -server -Xmx768m -Xms768m -Xmn384m -XX:PermSize=128m -Xss256k -XX:+DisableExplicitGC -XX:+UseConcMarkSweepGC -XX:+CMSParallelRemarkEnabled -XX:+UseCMSCompactAtFullCollection -XX:LargePageSizeInBytes=128m -XX:+UseFastAccessorMethods -XX:+UseCMSInitiatingOccupancyOnly -XX:CMSInitiatingOccupancyFraction=70 -Dlogpath.base=/root/logs/ids-mini -cp /app/ids-mini/config:/app/ids-mini/lib/* -jar /app/ids-mini/lib/IDS-Mini-server-1.0.0-SNAPSHOT.jar ids-mini
q11 ::  2021-07-22 13:08:54.675 [INFO ] 5235 --- [1    ] [-]  : Dog biting for cpu percent: 100.000000 > config max: 90.000000
q11 ::  2021-07-22 13:08:54.750 [INFO ] 5235 --- [1    ] [-]  : Dog barking for CPU占用第一, config:2/30s, item User: root Pid: 6460 Ppid: 6439 %cpu: 15 %mem: 2.2 VSZ: 6.3GB, RSS: 768.4MB Tty: ? Stat: Sl Start: 2021-07-19 19:35:47 Time: 09:52:21 Command: java -server -Xmx384m -Xms384m -Xmn384m -XX:PermSize=128m -Xss256k -XX:+DisableExplicitGC -XX:+UseConcMarkSweepGC -XX:+CMSParallelRemarkEnabled -XX:+UseCMSCompactAtFullCollection -XX:LargePageSizeInBytes=128m -XX:+UseFastAccessorMethods -XX:+UseCMSInitiatingOccupancyOnly -XX:CMSInitiatingOccupancyFraction=70 -Dlogpath.base=/root/logs/idsbp-server -cp /app/idsbp/config:/app/idsbp/lib/* -jar /app/idsbp/lib/idsbp-server-1.0.0-SNAPSHOT.jar idsbp-server
q11 ::  2021-07-22 13:21:14.669 [INFO ] 5235 --- [1    ] [-]  : Dog barking for cpu percent: 100.000000 > config max: 90.000000
q11 ::  2021-07-22 13:36:34.666 [INFO ] 5235 --- [1    ] [-]  : Dog biting for cpu percent: 100.000000 > config max: 90.000000
q11 ::  2021-07-22 13:36:34.723 [INFO ] 5235 --- [1    ] [-]  : Dog biting for CPU占用第一, item User: root Pid: 6460 Ppid: 6439 %cpu: 15 %mem: 2.2 VSZ: 6.3GB, RSS: 768.4MB Tty: ? Stat: Sl Start: 2021-07-19 19:35:47 Time: 09:57:51 Command: java -server -Xmx384m -Xms384m -Xmn384m -XX:PermSize=128m -Xss256k -XX:+DisableExplicitGC -XX:+UseConcMarkSweepGC -XX:+CMSParallelRemarkEnabled -XX:+UseCMSCompactAtFullCollection -XX:LargePageSizeInBytes=128m -XX:+UseFastAccessorMethods -XX:+UseCMSInitiatingOccupancyOnly -XX:CMSInitiatingOccupancyFraction=70 -Dlogpath.base=/root/logs/idsbp-server -cp /app/idsbp/config:/app/idsbp/lib/* -jar /app/idsbp/lib/idsbp-server-1.0.0-SNAPSHOT.jar idsbp-server
q11 ::  2021-07-22 13:36:34.748 [INFO ] 5235 --- [1    ] [-]  : LogItem: ENV, Value:  6460 ?        Sl   597:51 java -server -Xmx384m -Xms384m -Xmn384m -XX:PermSize=128m -Xss256k -XX:+DisableExplicitGC -XX:+UseConcMarkSweepGC -XX:+CMSParallelRemarkEnabled -XX:+UseCMSCompactAtFullCollection -XX:LargePageSizeInBytes=128m -XX:+UseFastAccessorMethods -XX:+UseCMSInitiatingOccupancyOnly -XX:CMSInitiatingOccupancyFraction=70 -Dlogpath.base=/root/logs/idsbp-server -cp /app/idsbp/config:/app/idsbp/lib/* -jar /app/idsbp/lib/idsbp-server-1.0.0-SNAPSHOT.jar idsbp-server IDSBP_SERVICE_PORT_10620_TCP_PROTO=tcp IDS_SS_PORT_10640_TCP_PORT=10640 IDSBP_SZ_SS_PORT_HTTP=10622 IDS_MSS_PORT=10650 HOSTNAME=idsbp-7fb5654c9b-lrmsb IDSBP_DEMO_SS_PORT=10621 IDS_QSS_PORT=10660 IDS_QS_PORT_10660_TCP=tcp://10.4.2.1:10660 AS_PORT_9611_TCP_PORT=9611 IDSPV2_SS_PORT=10610 IDS_SERVICE_PORT_10630_TCP_ADDR=127.0.4.1 IDSBP_SZ_SS_PORT=10622 IDS_SERVICE_PORT_10630_TCP_PROTO=tcp IDSBP_DEMO_SERVICE_PORT_10621_TCP=tcp://127.0.8.3:10621 IDS_CONSUMER_SS_HOST=127.0.2.1 KUBERNETES_PORT_443_TCP_PORT=443 KUBERNETES_PORT=tcp://127.0.0.1:443 enableDataPersist= IDSBP_SERVICE_PORT_10620_TCP_PORT=10620 IDS_QS_PORT_10660_TCP_PROTO=tcp IDSPV2_DEMO_SERVICE_PORT_10601_TCP=tcp://127.0.9.2:10601 umpBaseDir=/var/log/footstone/metrics KUBERNETES_SERVICE_PORT=443 AS_PORT_9611_TCP=tcp://127.0.227.224:9611 OLDPWD=/app/idsbp runOnBackground=false IDSPV2_SERVICE_PORT=tcp://10.4.1.0:10610 IDS_MINI_SERVICE_PORT_10650_TCP_PROTO=tcp IDS_SERVICE_PORT=tcp://127.0.4.1:10630 KUBERNETES_SERVICE_HOST=127.0.0.1 IDSPV2_SERVICE_PORT_10610_TCP_ADDR=10.4.1.0 IDSPV2_DEMO_SERVICE_PORT_10601_TCP_PORT=10601 IDS_CONSUMER_SERVICE_PORT_10401_TCP_PORT=10401 IDSBP_SERVICE_PORT=tcp://10.4.8.21:10620 IDSBP_SS_PORT=10620 IDS_QS_PORT_10660_TCP_PORT=10660 IDS_SERVICE_PORT=tcp://1.4.2.2:10400 IDS_SERVICE_PORT_10400_TCP=tcp://1.4.2.2:10400 IDSBP_DEMO_SERVICE_PORT_10621_TCP_PORT=10621 AUTHCODE_SERVICE_PORT=tcp://127.0.123.155:10007 IDSPV2_DEMO_SS_PORT_HTTP=10601 IS_PORT_10600_TCP_PORT=10600 IDS_SS_SERVICE_HOST=127.0.2.2 AUTHCODE_SERVICE_PORT_10007_TCP_PORT=10007 IDS_CONSUMER_SERVICE_PORT_10401_TCP=tcp://127.0.2.1:10401 AUTHCODE_SERVICE_PORT_10007_TCP_PROTO=tcp IDSPV2_SERVICE_PORT_10610_TCP_PORT=10610 IDSPV2_DEMO_SERVICE_PORT_10601_TCP_PROTO=tcp IS_PORT_10600_TCP_ADDR=127.0.1.2 IDS_SS_PORT_10640_TCP=tcp://127.0.2.2:10640 IS_PORT=tcp://127.0.1.2:10600 IDS_SS_HOST=127.0.4.1 IDS_SS_PORT_HTTP=10400 IDS_CONSUMER_SERVICE_PORT_10401_TCP_PROTO=tcp IDS_MINI_SERVICE_PORT_10650_TCP=tcp://10.4.18.21:10650 SERVER_ACCESSABLE_PORT=10000 IDS_SS_PORT=10400 IDS_MSS_PORT_HTTP=10650 ASS_HOST=127.0.227.224 ASS_HOST=127.0.123.155 IS_SERVICE_PORT=10600 IDSBP_SS_PORT_HTTP=10620 SERVER_ENVIROMENT= IDSBP_SZ_SERVICE_PORT_10622_TCP_ADDR=127.0.67.26 PATH=/opt/jdk/default/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin IDS_CONSUMER_SERVICE_PORT=tcp://127.0.2.1:10401 ASS_PORT_HTTP=9611 IDS_MINI_SERVICE_PORT=tcp://10.4.18.21:10650 MY_POD_NAME=idsbp-7fb5654c9b-lrmsb IDS_SS_PORT_10640_TCP_ADDR=127.0.2.2 IDSBP_DEMO_SERVICE_PORT_10621_TCP_PROTO=tcp IS_SERVICE_PORT_HTTP=10600 IDSPV2_SERVICE_PORT_10610_TCP=tcp://10.4.1.0:10610 PWD=/app/idsbp IDSBP_SZ_SS_HOST=127.0.67.26 AS_PORT=tcp://127.0.227.224:9611 JAVA_HOME=/opt/jdk/default IDSBP_SZ_SERVICE_PORT_10622_TCP=tcp://127.0.67.26:10622 IDS_QS_PORT=tcp://10.4.2.1:10660 IDS_SERVICE_PORT_10400_TCP_PROTO=tcp IDSBP_DEMO_SERVICE_PORT=tcp://127.0.8.3:10621 IDSBP_SERVICE_PORT_10620_TCP_ADDR=10.4.8.21 IDS_MINI_SERVICE_PORT_10650_TCP_PORT=10650 IDS_SERVICE_PORT_10630_TCP_PORT=10630 AS_PORT_9611_TCP_ADDR=127.0.227.224 IDS_CONSUMER_SS_PORT_HTTP=10401 LANG=zh_CN.UTF-8 AS_PORT_9611_TCP_PROTO=tcp IDSPV2_DEMO_SERVICE_PORT=tcp://127.0.9.2:10601 applicationDir=/app/idsbp IDSBP_SZ_SERVICE_PORT_10622_TCP_PROTO=tcp serviceName=idsbp IDSBP_SZ_SERVICE_PORT=tcp://127.0.67.26:10622 IDS_SS_PORT=10630 IDSPV2_DEMO_SS_PORT=10601 IDS_SERVICE_PORT_10630_TCP=tcp://127.0.4.1:10630 jdkDir=jdk1.8.0_191 IDSBP_SERVICE_PORT_10620_TCP=tcp://10.4.8.21:10620 ASS_PORT=10007 IDSPV2_DEMO_SS_HOST=127.0.9.2 HOME=/root SHLVL=2 IDS_SERVICE_PORT_10400_TCP_ADDR=1.4.2.2 KUBERNETES_PORT_443_TCP_PROTO=tcp IDS_MINI_SERVICE_PORT_10650_TCP_ADDR=10.4.18.21 IDSBP_SZ_SERVICE_PORT_10622_TCP_PORT=10622 IDSPV2_SS_HOST=10.4.1.0 IDSBP_DEMO_SERVICE_PORT_10621_TCP_ADDR=127.0.8.3 KUBERNETES_SERVICE_PORT_HTTPS=443 IDS_SS_PORT=tcp://127.0.2.2:10640 IDS_SS_PORT_10640_TCP_PROTO=tcp SERVER_ACCESSABLE_IP= IDS_CONSUMER_SERVICE_PORT_10401_TCP_ADDR=127.0.2.1 IDS_SS_HOST=1.4.2.2 IDSPV2_SS_PORT_HTTP=10610 IDS_SS_PORT_HTTP=10630 IDSBP_SS_HOST=10.4.8.21 AUTHCODE_SERVICE_PORT_10007_TCP_ADDR=127.0.123.155 IDS_MSS_HOST=10.4.18.21 IDSPV2_DEMO_SERVICE_PORT_10601_TCP_ADDR=127.0.9.2 CLASSPATH=/opt/jdk/default/lib/*.jar:/opt/jdk/default/jre/lib/*.jar env=TEST ASS_PORT=9611 IDS_SS_SERVICE_PORT_HTTP=10640 IDS_SS_SERVICE_PORT=10640 ASS_PORT_HTTP=10007 IS_PORT_10600_TCP=tcp://127.0.1.2:10600 IDS_QSS_HOST=10.4.2.1 IDSBP_DEMO_SS_HOST=127.0.8.3 enableUMP=true applicationLogDir=/root/logs/idsbp IDS_CONSUMER_SS_PORT=10401 KUBERNETES_PORT_443_TCP_ADDR=127.0.0.1 IDS_SERVICE_PORT_10400_TCP_PORT=10400 AUTHCODE_SERVICE_PORT_10007_TCP=tcp://127.0.123.155:10007 IS_PORT_10600_TCP_PROTO=tcp IDS_QS_PORT_10660_TCP_ADDR=10.4.2.1 KUBERNETES_PORT_443_TCP=tcp://127.0.0.1:443 IS_SERVICE_HOST=127.0.1.2 IDSBP_DEMO_SS_PORT_HTTP=10621 IDS_QSS_PORT_HTTP=10660 IDSPV2_SERVICE_PORT_10610_TCP_PROTO=tcp enableLogPersist=false _=/opt/jdk/default/bin/java
q11 ::  2021-07-22 13:36:34.951 [INFO ] 5235 --- [1    ] [-]  : LogItem: CWD, Value: /app/idsbp
q11 ::  2021-07-22 13:36:34.951 [INFO ] 5235 --- [1    ] [-]  : Kill interrupt to 6460 succeeded
q11 ::  2021-07-22 13:36:34.952 [INFO ] 5235 --- [1    ] [-]  : Kill killed to 6460 succeeded
q11 ::  2021-07-22 13:43:54.667 [INFO ] 5235 --- [1    ] [-]  : Dog barking for cpu percent: 90.599998 > config max: 90.000000
q15 ::  2021-07-22 11:29:48.921 [INFO ] 1428 --- [1    ] [-]  : log file created:/var/log/watchdog/dog.log
q15 ::  2021-07-22 11:29:48.921 [INFO ] 1428 --- [1    ] [-]  : dog with config: &{Topn:0 Pid:0 Ppid:0 Self:false KillSignals:[INT KILL] Interval:10s MaxMem:4294967296 MaxPmem:50 MaxPcpu:300 CmdFilter:[] MinAvailableMemory:2147483648 MaxHostPcpu:90 Whites:[kubelet kube- etcd] LogItems:[ENV CWD] RateConfig:2/30s limiter:0xc00000c660 Jitter:1s MaxTime:0s MaxTimeEnv:SHORT_TASK_FOR_DOG} created
q14 ::  2021-07-22 11:33:00.228 [INFO ] 847 --- [1    ] [-]  : Dog barking for cpu percent: 90.300003 > config max: 90.000000
q14 ::  2021-07-22 11:45:20.231 [INFO ] 847 --- [1    ] [-]  : Dog biting for cpu percent: 90.400002 > config max: 90.000000
q14 ::  2021-07-22 11:45:20.300 [INFO ] 847 --- [1    ] [-]  : Dog barking for CPU占用第一, config:2/30s, item User: elastic+ Pid: 13989 Ppid: 1 %cpu: 143 %mem: 16.4 VSZ: 70.5GB, RSS: 5.5GB Tty: ? Stat: SLsl Start: 2021-07-20 18:53:37 Time: 2-10:38:30 Command: /bin/java -Xms4g -Xmx4g -XX:+UseConcMarkSweepGC -XX:CMSInitiatingOccupancyFraction=75 -XX:+UseCMSInitiatingOccupancyOnly -XX:+AlwaysPreTouch -server -Djava.awt.headless=true -Dfile.encoding=UTF-8 -Djna.nosys=true -Djdk.io.permissionsUseCanonicalPath=true -Dio.netty.noUnsafe=true -Dio.netty.noKeySetOptimization=true -Dio.netty.recycler.maxCapacityPerThread=0 -Dlog4j.shutdownHookEnabled=false -Dlog4j2.disable.jmx=true -Dlog4j.skipJansi=true -XX:+HeapDumpOnOutOfMemoryError -Des.path.home=/usr/share/elasticsearch -Des.path.conf=/etc/elasticsearch/node1 -Des.distribution.flavor=default -Des.distribution.type=rpm -cp /usr/share/elasticsearch/lib/* org.elasticsearch.bootstrap.Elasticsearch -p /var/run/elasticsearch/192.168.1.1-node1/elasticsearch.pid --quiet
q14 ::  2021-07-22 12:10:20.261 [INFO ] 847 --- [1    ] [-]  : Dog barking for cpu percent: 93.800003 > config max: 90.000000
q14 ::  2021-07-22 12:32:40.230 [INFO ] 847 --- [1    ] [-]  : Dog biting for cpu percent: 93.599998 > config max: 90.000000
q14 ::  2021-07-22 12:32:40.308 [INFO ] 847 --- [1    ] [-]  : Dog biting for CPU占用第一, item User: elastic+ Pid: 13989 Ppid: 1 %cpu: 143 %mem: 16.7 VSZ: 70.5GB, RSS: 5.6GB Tty: ? Stat: SLsl Start: 2021-07-20 18:53:37 Time: 2-11:56:36 Command: /bin/java -Xms4g -Xmx4g -XX:+UseConcMarkSweepGC -XX:CMSInitiatingOccupancyFraction=75 -XX:+UseCMSInitiatingOccupancyOnly -XX:+AlwaysPreTouch -server -Djava.awt.headless=true -Dfile.encoding=UTF-8 -Djna.nosys=true -Djdk.io.permissionsUseCanonicalPath=true -Dio.netty.noUnsafe=true -Dio.netty.noKeySetOptimization=true -Dio.netty.recycler.maxCapacityPerThread=0 -Dlog4j.shutdownHookEnabled=false -Dlog4j2.disable.jmx=true -Dlog4j.skipJansi=true -XX:+HeapDumpOnOutOfMemoryError -Des.path.home=/usr/share/elasticsearch -Des.path.conf=/etc/elasticsearch/node1 -Des.distribution.flavor=default -Des.distribution.type=rpm -cp /usr/share/elasticsearch/lib/* org.elasticsearch.bootstrap.Elasticsearch -p /var/run/elasticsearch/192.168.1.1-node1/elasticsearch.pid --quiet
q14 ::  2021-07-22 12:32:40.349 [INFO ] 847 --- [1    ] [-]  : LogItem: ENV, Value: 13989 ?        SLsl 3596:36 /bin/java -Xms4g -Xmx4g -XX:+UseConcMarkSweepGC -XX:CMSInitiatingOccupancyFraction=75 -XX:+UseCMSInitiatingOccupancyOnly -XX:+AlwaysPreTouch -server -Djava.awt.headless=true -Dfile.encoding=UTF-8 -Djna.nosys=true -Djdk.io.permissionsUseCanonicalPath=true -Dio.netty.noUnsafe=true -Dio.netty.noKeySetOptimization=true -Dio.netty.recycler.maxCapacityPerThread=0 -Dlog4j.shutdownHookEnabled=false -Dlog4j2.disable.jmx=true -Dlog4j.skipJansi=true -XX:+HeapDumpOnOutOfMemoryError -Des.path.home=/usr/share/elasticsearch -Des.path.conf=/etc/elasticsearch/node1 -Des.distribution.flavor=default -Des.distribution.type=rpm -cp /usr/share/elasticsearch/lib/* org.elasticsearch.bootstrap.Elasticsearch -p /var/run/elasticsearch/192.168.1.1-node1/elasticsearch.pid --quiet HOSTNAME=tencent-beta19 SHELL=/sbin/nologin ES_GROUP=elasticsearch OLDPWD=/usr/share/elasticsearch DATA_DIR=/opt/elasticsearch/data/192.168.1.1-node1 ES_STARTUP_SLEEP_TIME=5 USER=elasticsearch ES_HOME=/usr/share/elasticsearch ES_USER=elasticsearch MAX_LOCKED_MEMORY=unlimited PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin CONF_DIR=/etc/elasticsearch/node1 PWD=/usr/share/elasticsearch LOG_DIR=/opt/elasticsearch/logs/192.168.1.1-node1 LANG=en_US.utf8 ES_JVM_OPTIONS=/etc/elasticsearch/node1/jvm.options SHLVL=0 HOME=/nonexistent LOGNAME=elasticsearch MAX_OPEN_FILES=65536 MAX_THREADS=8192 MAX_MAP_COUNT=262144 ES_PATH_CONF=/etc/elasticsearch/node1 PID_DIR=/var/run/elasticsearch/192.168.1.1-node1
q14 ::  2021-07-22 12:32:52.835 [INFO ] 847 --- [1    ] [-]  : LogItem: CWD, Value: /usr/share/elasticsearch
q14 ::  2021-07-22 12:32:52.835 [INFO ] 847 --- [1    ] [-]  : Kill interrupt to 13989 succeeded
q14 ::  2021-07-22 12:32:52.836 [INFO ] 847 --- [1    ] [-]  : Kill killed to 13989 succeeded
q7  ::  2021-07-22 12:11:51.467 [INFO ] 13927 --- [1    ] [-]  : log file created:/var/log/watchdog/dog.log
q7  ::  2021-07-22 12:11:51.467 [INFO ] 13927 --- [1    ] [-]  : dog with config: &{Topn:0 Pid:0 Ppid:0 Self:false KillSignals:[INT KILL] Interval:10s MaxMem:4294967296 MaxPmem:50 MaxPcpu:300 CmdFilter:[] MinAvailableMemory:2147483648 MaxHostPcpu:90 Whites:[kubelet kube- etcd nexus] LogItems:[ENV CWD] RateConfig:2/30s limiter:0xc00000c0a8 Jitter:1s MaxTime:0s MaxTimeEnv:SHORT_TASK_FOR_DOG} created
q8  ::  2021-07-22 11:29:52.248 [INFO ] 9036 --- [1    ] [-]  : log file created:/var/log/watchdog/dog.log
q8  ::  2021-07-22 11:29:52.248 [INFO ] 9036 --- [1    ] [-]  : dog with config: &{Topn:0 Pid:0 Ppid:0 Self:false KillSignals:[INT KILL] Interval:10s MaxMem:4294967296 MaxPmem:50 MaxPcpu:300 CmdFilter:[] MinAvailableMemory:2147483648 MaxHostPcpu:90 Whites:[kubelet kube- etcd] LogItems:[ENV CWD] RateConfig:2/30s limiter:0xc000112090 Jitter:1s MaxTime:0s MaxTimeEnv:SHORT_TASK_FOR_DOG} created
```
