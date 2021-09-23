# busy

from https://github.com/vikyd/go-cpu-load

Generate CPU load on Windows/Linux/Mac.

# Usage

example 01: run 30% of all CPU cores for 10 seconds.

`dogbusy -p 30 -d 10s`

example 02: run 30% of all CPU cores forever.

`dogbusy -p 30`

example 03: run 30% of 2 of CPU cores for 10 seconds.

`dogbusy -p 30 -c 2 -d 10s`

- `top CPU load` = `c` \* `p`
- may not specify cores run the load only, it just promises the `all CPU load`, and not promise each cores run the same
  load

# Parameters

```
Usage of dogbusy:
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
