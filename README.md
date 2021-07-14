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
$ dog -max-pcpu 500 -filter busy   
2021/07/14 16:32:06 Dog biting for 3, item User: bingoobjca Pid: 47526 Ppid: 72811 %cpu: 550.2 %mem: 0 VSZ: 5.1GB, RSS: 3.6MB Tty: s002 Stat: S+ Start: 2021-07-14 08:31:59 Time: 0:14.49 Command: busy -p50
2021/07/14 16:32:06 Kill interrupt to 47526 succeeded
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
  -kill string 发送信号，eg INT TERM KILL QUIT USR1 USR2 (default "INT")
  -max-mem value 允许最大内存 (默认 0B，不检查内存)
  -max-pcpu int 允许内存最大百分比, eg 0-1200 (默认 600), 0 不检查 CPU
  -max-pmem int 允许CPU最大百分比, eg 1-100 (默认 50)
  -pid int 指定pid
  -ppid int 指定ppid
  -self 是否监控自身
  -span duration 检查时间间隔 (默认 10s)
  -topn int 只取前N个检查
  
$ busy -h               
Usage of busy:
  -c int 使用核数，默认 12
  -d duration 跑多久，默认一直跑
  -m string 总内存,增量, eg. 1) 10M 直接达到10M 2) 10M,1K/10s 总用量10M,每10秒增加1K
  -p int 每核CPU百分比 (默认 100)
```

