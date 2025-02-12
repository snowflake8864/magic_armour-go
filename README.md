### 编译
#### 目标平台为 windows

      go build -o app_agent.exe    -ldflags "-w -s"  -trimpath ./cmd/agent/main.go
      go build -o app_api.exe      -ldflags "-w -s"  -trimpath ./cmd/api/main.go
      go build -o app_cli.exe      -ldflags "-w -s"  -trimpath ./cmd/cli/main.go
      go build -o app_org.exe      -ldflags "-w -s"  -trimpath ./cmd/organization/main.go
      go build -o app_platform.exe -ldflags "-w -s"  -trimpath ./cmd/platform/main.go


#### 目标平台为 linux(通过windows系统直接编译出linux可执行程序)

      // cpu架构，amd64、arm64
      go env -w GOARCH=amd64
      // 程序运行的最终系统，linux、windows、darwin(苹果macos系统)
      go env -w GOOS=linux
      // window编译设置Cgo模块关闭，因为windows上做cgo开发太麻烦，如果引用了Cgo库库，那么请在linux环境开发、编译
      go env -w CGO_ENABLED=0
      // 编译出最终可执行文件，进入根目录（也就是 go.mod 所在的目录）
      go build -o app_agent.linux64     -ldflags "-w -s"  -trimpath  ./cmd/agent/main.go
      go build -o app_api.linux64       -ldflags "-w -s"  -trimpath  ./cmd/api/main.go
      go build -o app_cli.linux64       -ldflags "-w -s"  -trimpath  ./cmd/cli/main.go
      go build -o app_org.linux64       -ldflags "-w -s"  -trimpath  ./cmd/organization/main.go
      go build -o app_platform.linux64  -ldflags "-w -s"  -trimpath  ./cmd/platform/main.go


### 命令行入口

      ./app_cli.linux64 migrate -A (init|dictionary|menu|append)
      ./app_cli.linux64 queue  -A (init|consume) -T (common|process|dll|threaten|alarm|syslog|crontab)
      ./app_cli.linux64 crontab -U (minute|hour|day) -M(every|once) -A (zero|one|two|three|four|five|six|ten|fifteen|twenty|thirty)


### 守护进程，注意<kbd>supervisord.conf</kbd>路径

      supervisorctl -c /data/supervisor/etc/supervisord.conf reload
      supervisorctl -c /data/supervisor/etc/supervisord.conf restart all
      systemctl status supervisord.service


### 管理平台隐藏页面

- /#/menuConfiguration
- /#/dataDictionary

### 注意事项

- 服务器时区应设置为东8区: 日期转成时间戳时默认按照东8区计算
- 每个路由都有限流设置
- Doris 数据库

      DROP TABLE IF EXISTS table_name FORCE;


### 定时任务，参考 <kbd>docs/Crontab.txt</kbd>