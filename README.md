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

      #每隔[1分钟]
      */1 * * * * app_cli crontab -U minute -M every -I one -A zero
      */1 * * * * app_cli crontab -U minute -M every -I one  -A one
      */1 * * * * app_cli crontab -U minute -M every -A -I one  two
      */1 * * * * app_cli crontab -U minute -M every -A -I one  three

      #每隔[2分钟]
      */2 * * * * app_cli crontab -U minute -M every -A -I two  zero
      */2 * * * * app_cli crontab -U minute -M every -A -I two  one
      */2 * * * * app_cli crontab -U minute -M every -A -I two  two
      */2 * * * * app_cli crontab -U minute -M every -A -I two  three
      
      #每隔[3分钟]
      */3 * * * * app_cli crontab -U minute -M every  -I three  -A zero
      */3 * * * * app_cli crontab -U minute -M every -I three  -A one
      */3 * * * * app_cli crontab -U minute -M every -I three  -A two

      #每隔[5分钟]
      */5 * * * * app_cli crontab -U minute -M every -I five  -A zero
      */5 * * * * app_cli crontab -U minute -M every -I five  -A one
      */5 * * * * app_cli crontab -U minute -M every -I five  -A two

      #每隔[10分钟]
      */10 * * * * app_cli crontab -U minute -M every -I ten -A zero
      */10 * * * * app_cli crontab -U minute -M every -I ten -A one

      #每隔[15分钟]
      */15 * * * * app_cli crontab -U minute -M every -I fifteen -A zero
      */15 * * * * app_cli crontab -U minute -M every -I fifteen -A one

       #每隔[30分钟]
      */30 * * * * app_cli crontab -U minute -M every -I thirty -A zero
      */30 * * * * app_cli crontab -U minute -M every -I thirty -A one

      #每隔[1小时]
      0 */1 * * * app_cli crontab -U minute -M every -I one -A zero
      1 */1 * * * app_cli crontab -U minute -M every -I one -A one

      #每天一次
      10 1 * * * app_cli crontab -U day -M once -A zero
      11 1 * * * app_cli crontab -U day -M once -A one
      12 1 * * * app_cli crontab -U day -M once -A two
      13 1 * * * app_cli crontab -U day -M once -A three
      14 1 * * * app_cli crontab -U day -M once -A four
      15 1 * * * app_cli crontab -U day -M once -A five
      16 1 * * * app_cli crontab -U day -M once -A six
      17 1 * * * app_cli crontab -U day -M once -A seven
      20 2 * * * app_cli crontab -U day -M once -A eight
      30 3 * * * app_cli crontab -U day -M once -A nine
      40 4 * * * app_cli crontab -U day -M once -A ten

      #每天两次次
      30 3,13 * * * app_cli crontab -U day -M every -A zero
      31 3,13 * * * app_cli crontab -U day -M every -A one
      33 3,13 * * * app_cli crontab -U day -M every -A two
