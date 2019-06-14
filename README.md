# dog
watchdog for RSS/CPU or Heartbeating

## demo

```go

dog := &Dog{
    MaxMem:        100 * 1024 * 1024, // 最大内存 100M
    MaxCpuPercent: 90,                // 最大CPU占比90%
    MaxMemPercent: 90,                // 最大内存占比90%
    BiteLive:      false,             // 咬了不活，直接死掉
}

dog.FreeDog()  // 放狗看门（不阻塞)

```
