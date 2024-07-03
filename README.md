### 编译
- 目标平台为 windows

      go build -o demo_cli.exe -ldflags "-w -s"  -trimpath ./cmd/cli/main.go

- 目标平台为 linux(通过 windows 系统直接编译出 linux 可执行文件)
  
      go build -o demo_cli.exe -ldflags "-w -s"  -trimpath ./cmd/cli/main.go

### 初始化
      ./demo_cli.exe migrate -A (init|dictionary|menu|append)
      ./demo_cli.exe queue  queue -A (init|consume) -T (common|process|threaten|alarm|syslog)
      ./demo_cli.exe queue  crontab -U (minute|hour|day) -M(every|once) -A (zero|one|two|three|four|five|six|ten|fifteen|twenty|thirty)

