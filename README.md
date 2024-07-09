### 编译
- 目标平台为 windows
     
      go build -o app_agent.exe -ldflags "-w -s"  -trimpath ./cmd/agent/main.go
      go build -o app_api.exe -ldflags "-w -s"  -trimpath ./cmd/api/main.go
      go build -o app_cli.exe -ldflags "-w -s"  -trimpath ./cmd/cli/main.go
      go build -o app_org.exe -ldflags "-w -s"  -trimpath ./cmd/organization/main.go
      go build -o app_platform.exe -ldflags "-w -s"  -trimpath ./cmd/platform/main.go

  
- 目标平台为 linux(通过 windows 系统直接编译出 linux 可执行文件)
  
      // cpu架构
      go env -w GOARCH=amd64
      // 程序运行的最终系统，linux、windows、darwin(苹果macos系统)
      go env -w GOOS=linux
      // window编译设置Cgo模块关闭，因为windows上做cgo开发太麻烦，如果引用了Cgo库库，那么请在linux环境开发、编译  
      go env -w CGO_ENABLED=0
      // 编译出最终可执行文件，进入根目录（GinSkeleton所在目录，也就是 go.mod 所在的目录）
      go build -o app_agent.linux64  -ldflags "-w -s"  -trimpath  ./cmd/agent/main.go
      go build -o app_api.linux64  -ldflags "-w -s"  -trimpath  ./cmd/api/main.go
      go build -o app_cli.linux64  -ldflags "-w -s"  -trimpath  ./cmd/cli/main.go
      go build -o app_org.linux64  -ldflags "-w -s"  -trimpath  ./cmd/organization/main.go
      go build -o app_platform.linux64  -ldflags "-w -s"  -trimpath  ./cmd/platform/main.go


### 命令行入口

      ./app_cli.exe migrate -A (init|dictionary|menu|append)
      ./app_cli.exe queue  -A (init|consume) -T (common|process|threaten|alarm|syslog|crontab)
      ./app_cli.exe crontab -U (minute|hour|day) -M(every|once) -A (zero|one|two|three|four|five|six|ten|fifteen|twenty|thirty)

### 定时任务

      #每隔[1分钟], 已按顺序使用到[two]
      */1 * * * * app_cli.linux64   crontab -U minute -I one -A zero
      */1 * * * * app_cli.linux64   crontab -U minute -I one  -A one
      */1 * * * * app_cli.linux64   crontab -U minute -I one -A two
      */1 * * * * app_cli.linux64   crontab -U minute -I one -A four
      */1 * * * * app_cli.linux64   crontab -U minute -I one -A three

      #每隔[2分钟], 已按顺序使用到[one]
      */2 * * * * app_cli.linux64   crontab -U minute -I two -A zero
      */2 * * * * app_cli.linux64   crontab -U minute -I two -A one
      */2 * * * * app_cli.linux64   crontab -U minute -I two -A two
      */2 * * * * app_cli.linux64   crontab -U minute -I two -A three
      
      #每隔[3分钟], 已按顺序使用到[two]
      */3 * * * * app_cli.linux64   crontab -U minute -I three -A zero
      */3 * * * * app_cli.linux64   crontab -U minute -I three -A one
      */3 * * * * app_cli.linux64   crontab -U minute -I three -A two
      */3 * * * * app_cli.linux64   crontab -U minute -I three -A three
      */3 * * * * app_cli.linux64   crontab -U minute -I three -A four

      #每隔[5分钟], 已按顺序使用到[zero]
      */5 * * * * app_cli.linux64   crontab -U minute -I five -A zero
      */5 * * * * app_cli.linux64   crontab -U minute -I five -A one
      */5 * * * * app_cli.linux64   crontab -U minute -I five -A two

      #每隔[10分钟], 已按顺序使用到[zero]
      */10 * * * * app_cli.linux64   crontab -U minute -I ten -A zero
      */10 * * * * app_cli.linux64   crontab -U minute -I ten -A one

      #每隔[15分钟], 已按顺序使用到[zero]
      */15 * * * * app_cli.linux64   crontab -U minute -I fifteen -A zero
      */15 * * * * app_cli.linux64   crontab -U minute -I fifteen -A one

       #每隔[30分钟], 已按顺序使用到[]
      */30 * * * * app_cli.linux64   crontab -U minute -I thirty -A zero
      */30 * * * * app_cli.linux64   crontab -U minute -I thirty -A one

      #每隔[1小时], 已按顺序使用到[]
      0 */1 * * * app_cli.linux64   crontab -U hour -I one -A zero
      1 */1 * * * app_cli.linux64   crontab -U hour -I one -A one

      #每天一次, 已按顺序使用到[eight]
      10 1 * * * app_cli.linux64   crontab -U day -M once -A zero
      11 1 * * * app_cli.linux64   crontab -U day -M once -A one
      12 1 * * * app_cli.linux64   crontab -U day -M once -A two
      13 1 * * * app_cli.linux64   crontab -U day -M once -A three
      14 1 * * * app_cli.linux64   crontab -U day -M once -A four
      15 1 * * * app_cli.linux64   crontab -U day -M once -A five
      16 1 * * * app_cli.linux64   crontab -U day -M once -A six
      17 1 * * * app_cli.linux64   crontab -U day -M once -A seven
      20 2 * * * app_cli.linux64   crontab -U day -M once -A eight
      30 3 * * * app_cli.linux64   crontab -U day -M once -A nine
      40 4 * * * app_cli.linux64   crontab -U day -M once -A ten

      #每天两次, 已按顺序使用到[zero]
      30 3,13 * * * app_cli.linux64   crontab -U day -M every -A zero
      31 3,13 * * * app_cli.linux64   crontab -U day -M every -A one
