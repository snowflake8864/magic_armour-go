### 编译
- 目标平台为 windows

      go build -o app_cli.exe -ldflags "-w -s"  -trimpath ./cmd/cli/main.go

- 目标平台为 linux(通过 windows 系统直接编译出 linux 可执行文件)
  
      go build -o app_cli.exe -ldflags "-w -s"  -trimpath ./cmd/cli/main.go

### 命令行入口

      ./app_cli.exe migrate -A (init|dictionary|menu|append)
      ./app_cli.exe queue  queue -A (init|consume) -T (common|process|threaten|alarm|syslog)
      ./app_cli.exe queue  crontab -U (minute|hour|day) -M(every|once) -A (zero|one|two|three|four|five|six|ten|fifteen|twenty|thirty)

### 定时任务

      #每天1:00执行
      00 1 * * * app_cli crontab -U day -M once -A five 