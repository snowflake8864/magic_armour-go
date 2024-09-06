package test

import (
	"fmt"
	"goskeleton/app/utils/gomail_factory"
	"goskeleton/app/utils/helpers/date_helper"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"testing"
)

func TestGomail(t *testing.T) {

	fact := gomail_factory.CreateMailerFactory()

	body := "Hello <b>Gomail!</b> This is an <a href='https://github.com/go-gomail/gomail'>HTML</a> email."

	if err := fact.SendHtmlMsg("sub@veda.com", "幻甲终端动态防御系统-安全告警通知", body); err != nil {
		fmt.Println(err.Error())
	}

	if err := fact.SendHtmlMsg("sub@veda.com", "幻甲终端动态防御系统-安全告警通知", body); err != nil {
		fmt.Println(err.Error())
	}

}

func TestEmailNotice(t *testing.T) {
	addr := "sub@veda.com"
	fact := gomail_factory.CreateMailerFactory()

	/*var alarmList []*gomail_factory.AlarmCycleAlarm
	alarmList = append(alarmList, &gomail_factory.AlarmCycleAlarm{
		Name: "进程告警",
		Add:  10,
	}, &gomail_factory.AlarmCycleInfo{
		Name: "模块告警",
		Add:  53,
	})

	param := gomail_factory.AlarmCycle{
		DateStart:  "2024-08-30 00:00:00",
		DateEnd:    "2024-09-06 00:00:00",
		AlarmTotal: 63,
		AlarmList:  alarmList,
	}
	subject, context := gomail_factory.CreateTemplateFactory().GetAlarmCycle(&param)

	if err := fact.SendHtmlMsg(addr, subject, context); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("=============")
	param1 := gomail_factory.AlarmCycleInfo{
		DateStart:  "2024-08-30 00:00:00",
		DateEnd:    "2024-09-06 00:00:00",
		AlarmTotal: 0,
		AlarmList:  nil,
	}
	subject1, context1 := gomail_factory.CreateTemplateFactory().GetAlarmCycle(&param1)
	if err := fact.SendHtmlMsg(addr, subject1, context1); err != nil {
		fmt.Println(err.Error())
	}*/

	/*param2 := gomail_factory.AlarmTimely{
		Date:          date_helper.CurDatetime(),
		AgentId:       "8ac94fd732cfe3628922262747607f67",
		AgentIp:       "192.168.4.6",
		AgentHostname: "local.name",
		AlarmLogin:    map[string]int{"192.168.4.6": 12, "192.168.4.7": 30},
	}
	subject2, context2 := gomail_factory.CreateTemplateFactory().GetAlarmTimely(&param2)
	if err := fact.SendHtmlMsg(addr, subject2, context2); err != nil {
		fmt.Println(err.Error())
	}*/

	var usageList []*gomail_factory.HardwareTimelyInfo
	usageList = append(usageList, &gomail_factory.HardwareTimelyInfo{
		Name:      "CPU阈值",
		Threshold: 80,
		Usage:     90.55,
	}, &gomail_factory.HardwareTimelyInfo{
		Name:      "内存阈值",
		Threshold: 70,
		Usage:     85.34,
	}, &gomail_factory.HardwareTimelyInfo{
		Name:      "硬盘阈值",
		Threshold: 60,
		Usage:     65.99,
	})
	param3 := gomail_factory.HardwareTimely{
		Date:          date_helper.CurDatetime(),
		AgentId:       "8ac94fd732cfe3628922262747607f67",
		AgentIp:       "192.168.4.6",
		AgentHostname: "local.name",
		UsageList:     usageList,
	}
	subject3, context3 := gomail_factory.CreateTemplateFactory().GetHardwareTimely(&param3)
	if err := fact.SendHtmlMsg(addr, subject3, context3); err != nil {
		fmt.Println(err.Error())
	}
}
