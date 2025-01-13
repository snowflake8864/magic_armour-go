## 清理磁盘日志
10 1 * * * /home/veda/clear_disk.sh >> /var/log/armour_disk_clear.log 2>&1

## 每隔[1分钟], 已按顺序使用
# 威胁情报缓存同步默认状态（未同步）数据到威胁情报中心
*/1 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I one -A zero >> /var/log/armour_crontab.log 2>&1
# 终端离线
*/1 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I one -A one  >> /var/log/armour_crontab.log 2>&1
*/1 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I one -A two >> /var/log/armour_crontab.log 2>&1
*/1 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I one -A three  >> /var/log/armour_crontab.log 2>&1
*/1 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I one -A four >> /var/log/armour_crontab.log 2>&1

## 每隔[2分钟], 已按顺序使用
# 动态封堵已过期的自动解封
*/2 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I two -A zero >> /var/log/armour_crontab.log 2>&1
# 终端安全状态
*/2 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I two -A one >> /var/log/armour_crontab.log 2>&1
*/2 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I two -A two >> /var/log/armour_crontab.log 2>&1
*/2 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I two -A three >> /var/log/armour_crontab.log 2>&1

## 每隔[3分钟], 已按顺序使用
# 从进程或模块添加任务
*/3 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I three -A zero  >> /var/log/armour_crontab.log 2>&1
# 服务器状态监控
*/3 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I three -A one  >> /var/log/armour_crontab.log 2>&1
# 威胁情报同步数据到威胁情报中心
*/3 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I three -A two >> /var/log/armour_crontab.log 2>&1
*/3 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I three -A three >> /var/log/armour_crontab.log 2>&1
*/3 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I three -A four >> /var/log/armour_crontab.log 2>&1

## 每隔[5分钟], 已按顺序使用
# 同步推送数据到幻势
*/5 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I five -A zero >> /var/log/armour_crontab.log 2>&1
*/5 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I five -A one >> /var/log/armour_crontab.log 2>&1
*/5 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I five -A two >> /var/log/armour_crontab.log 2>&1

## 每隔[10分钟], 已按顺序使用
# 统计告警今天数据
*/10 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I ten -A zero >> /var/log/armour_crontab.log 2>&1
*/10 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I ten -A one >> /var/log/armour_crontab.log 2>&1

## 每隔[15分钟], 已按顺序使用
# 统计企业今天数据
*/15 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I fifteen -A zero >> /var/log/armour_crontab.log 2>&1
*/15 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I fifteen -A one >> /var/log/armour_crontab.log 2>&1

## 每隔[30分钟], 已按顺序使用
*/30 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I thirty -A zero >> /var/log/armour_crontab.log 2>&1
*/30 * * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U minute -I thirty -A one >> /var/log/armour_crontab.log 2>&1

## 每隔[1小时], 已按顺序使用
# 威胁情报缓存修复进程模块
0 */1 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U hour -I one -A zero >> /var/log/armour_crontab.log 2>&1
1 */1 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U hour -I one -A one >> /var/log/armour_crontab.log 2>&1

## 每天一次, 已按顺序使用
# 每天清理一次数据
10 1 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A zero >> /var/log/armour_crontab.log 2>&1
# 订单自动过期
11 1 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A one >> /var/log/armour_crontab.log 2>&1
# 终端授权自动过期
12 1 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A two >> /var/log/armour_crontab.log 2>&1
# 统计企业昨天数据
13 1 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A three >> /var/log/armour_crontab.log 2>&1
# 统计告警昨天数据
14 1 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A four >> /var/log/armour_crontab.log 2>&1
# 删除威胁情报缓存
15 1 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A five >> /var/log/armour_crontab.log 2>&1
# 威胁情报缓存同步未知状态数据到威胁情报中心
16 1 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A six >> /var/log/armour_crontab.log 2>&1
# 告警通知
20 2 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A seven >> /var/log/armour_crontab.log 2>&1
21 2 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A eight >> /var/log/armour_crontab.log 2>&1
30 3 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A nine >> /var/log/armour_crontab.log 2>&1
40 4 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M once -A ten >> /var/log/armour_crontab.log 2>&1

## 每天两次, 已按顺序使用
# 从进程或模块添加到威胁情报缓存
30 3,13 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M every -A zero >> /var/log/armour_crontab.log 2>&1
31 3,13 * * * cd /home/wwwroot/armour_go && ./app_cli.linux64 crontab -U day -M every -A one >> /var/log/armour_crontab.log 2>&1

