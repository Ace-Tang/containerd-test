# simple test for containerd run with runv

## contianerd config

```
$ cat /etc/containerd/config.toml 
root = "/var/lib/containerd"
state = "/run/containerd"
subreaper = true
oom_score = 0

[grpc]
  address = "/run/containerd/containerd.sock"
  uid = 0
  gid = 0

[debug]
  address = "/run/containerd/debug.sock"
  uid = 0
  gid = 0
  level = "debug"

[metrics]
  address = ""

[cgroup]
  path = ""

[plugins.linux]
 shim = "containerd-shim"
 no_shim = false
 runtime = "runv"
 shim_debug = true
```

## run test

1. go build -o run run.go

2. run test
```
$ sudo ./run docker.io/library/busybox:latest runv79 runv  echo 1
INFO[0000] start new task                               
INFO[0005] start wait task                              
INFO[0005] start task successful                        
1
INFO[0005] task exit with status 0         
```
