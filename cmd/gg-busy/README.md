# busy

from https://github.com/vikyd/go-cpu-load

Generate CPU loads on Windows/Linux/Mac.

# Usage

example 01: run 30% of all CPU cores for 10 seconds.

`gg-busy -p 30 -d 10s`

example 02: run 30% of all CPU cores forever.

`gg-busy -p 30`

example 03: run 30% of 2 of CPU cores for 10 seconds.

`gg-busy -p 30 -c 2 -d 10s`

- `top CPU load` = `c` \* `p`
- may not specify cores run the load only, it just promises the `all CPU load`, and not promise each cores run the same
  load

# Parameters

```
Usage of gg-busy:
  -c int
        how many cores (default cpu cores)
  -d duration
        how long
  -p int
        percentage of each specify cores (default 100)
```

## How it runs

- Giving a range of time(e.g. 100ms)
- Want to run 30% of all CPU cores
    - 30ms: run (CPU 100%)
    - 70ms: sleep(CPU 0%)

## resources

### stress

1. 模拟一个 CPU 使用率 100% 的场景，跑 600 秒：`stress --cpu 1 --timeout 600`
2. 模拟 I/O 压力，即不停地执行 sync： `stress -i 1 --timeout 600`  -i, --io N 产生 N 个进程，每个进程反复调用 sync() 将内存上的内容写到硬盘上
3. 模拟的是 8 个进程： `stress -c 8 --timeout 600`， -c, --cpu N 产生 N 个进程，每个进程都反复不停的计算随机数的平方根
4. 消耗内存：产生两个子进程，每个进程分配 300M 内存：`stress --vm 2 --vm-bytes 300M --vm-keep`  -m, --vm N 产生 N 个进程，每个进程不断分配和释放内存；
