# go-cron - A user-land cron alternative for use in Docker containers

This app reads a crontab file and makes runs it just like `crond` would.
The reason this app was created is specifically geared towards the use in a Docker container:
- `go-cron` can run as non-root (`crond` wants to create a PID file and can only run as root)
- `go-cron` can run with per second accuracy if the crontab has 6 time fragments
- `go-cron` is not strict about the ownership and permissions of the crontab files.


## Building
1. Build a binary using `go build`
2. Now you have `go-cron` as an executable binary

## Usage
```
./go-cron -file crontab.txt
```
Consider using `dumb-init` to prevent zombie processes. I'm not sure if it is needed, as my understanding of the init system is limited, but using `dumb-init` does not hurt either.
```
/usr/bin/dumb-init -- ./go-cron -file crontab.txt
```
