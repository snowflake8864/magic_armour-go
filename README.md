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

      ./app_cli.linux64 migrate -A (init|dictionary|menu|append)
      ./app_cli.exe queue  -A (init|consume) -T (common|process|threaten|alarm|syslog|crontab)
      ./app_cli.exe crontab -U (minute|hour|day) -M(every|once) -A (zero|one|two|three|four|five|six|ten|fifteen|twenty|thirty)


### 守护进程，注意<kbd>supervisord.conf</kbd>路径

      1. supervisorctl -c /data/supervisor/etc/supervisord.conf reload
      2. supervisorctl -c /data/supervisor/etc/supervisord.conf restart all
      3. systemctl restart supervisord.service


### 管理平台隐藏页面

      /#/menuConfiguration
      /#/dataDictionary


### 定时任务

      #每隔[1分钟], 已按顺序使用到[two]
      */1 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I one -A zero >> /var/log/armour_crontab.log 2>&1
      */1 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I one -A one  >> /var/log/armour_crontab.log 2>&1
      */1 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I one -A two >> /var/log/armour_crontab.log 2>&1
      */1 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I one -A three  >> /var/log/armour_crontab.log 2>&1
      */1 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I one -A four >> /var/log/armour_crontab.log 2>&1
      
      #每隔[2分钟], 已按顺序使用到[one]
      */2 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I two -A zero >> /var/log/armour_crontab.log 2>&1
      */2 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I two -A one >> /var/log/armour_crontab.log 2>&1
      */2 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I two -A two >> /var/log/armour_crontab.log 2>&1
      */2 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I two -A three >> /var/log/armour_crontab.log 2>&1
      
      #每隔[3分钟], 已按顺序使用到[two]
      */3 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I three -A zero  >> /var/log/armour_crontab.log 2>&1
      */3 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I three -A one  >> /var/log/armour_crontab.log 2>&1
      */3 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I three -A two >> /var/log/armour_crontab.log 2>&1
      */3 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I three -A three >> /var/log/armour_crontab.log 2>&1
      */3 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I three -A four >> /var/log/armour_crontab.log 2>&1
      
      #每隔[5分钟], 已按顺序使用到[zero]
      */5 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I five -A zero >> /var/log/armour_crontab.log 2>&1
      */5 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I five -A one >> /var/log/armour_crontab.log 2>&1
      */5 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I five -A two >> /var/log/armour_crontab.log 2>&1
      
      #每隔[10分钟], 已按顺序使用到[zero]
      */10 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I ten -A zero >> /var/log/armour_crontab.log 2>&1
      */10 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I ten -A one >> /var/log/armour_crontab.log 2>&1
      
      #每隔[15分钟], 已按顺序使用到[zero]
      */15 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I fifteen -A zero >> /var/log/armour_crontab.log 2>&1
      */15 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I fifteen -A one >> /var/log/armour_crontab.log 2>&1
      
      #每隔[30分钟], 已按顺序使用到[]
      */30 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I thirty -A zero >> /var/log/armour_crontab.log 2>&1
      */30 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I thirty -A one >> /var/log/armour_crontab.log 2>&1
      
      #每隔[1小时], 已按顺序使用到[]
      0 */1 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U hour -I one -A zero >> /var/log/armour_crontab.log 2>&1
      1 */1 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U hour -I one -A one >> /var/log/armour_crontab.log 2>&1
      
      #每天一次, 已按顺序使用到[eight]
      10 1 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A zero >> /var/log/armour_crontab.log 2>&1
      11 1 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A one >> /var/log/armour_crontab.log 2>&1
      12 1 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A two >> /var/log/armour_crontab.log 2>&1
      13 1 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A three >> /var/log/armour_crontab.log 2>&1
      14 1 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A four >> /var/log/armour_crontab.log 2>&1
      15 1 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A five >> /var/log/armour_crontab.log 2>&1
      16 1 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A six >> /var/log/armour_crontab.log 2>&1
      17 1 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A seven >> /var/log/armour_crontab.log 2>&1
      20 2 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A eight >> /var/log/armour_crontab.log 2>&1
      30 3 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A nine >> /var/log/armour_crontab.log 2>&1
      40 4 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A ten >> /var/log/armour_crontab.log 2>&1
      
      #每天两次, 已按顺序使用到[zero]
      30 3,13 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M every -A zero >> /var/log/armour_crontab.log 2>&1
      31 3,13 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M every -A one >> /var/log/armour_crontab.log 2>&1
    
