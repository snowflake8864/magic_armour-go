package test

import (
	"fmt"
	"goskeleton/app/utils/gomail_factory"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"goskeleton/command/console/alarm_cmd"
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
	alarm_cmd.CreateCmdFactory().EmailNotice()
}

func TestGetAlarmCycle(t *testing.T) {
	var alarmList []*gomail_factory.AlarmCycleAlarm
	alarmList = append(alarmList, &gomail_factory.AlarmCycleAlarm{
		Name: "进程告警",
		Add:  10,
	}, &gomail_factory.AlarmCycleAlarm{
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
	fmt.Println(subject)
	fmt.Println("=============================")
	fmt.Println(context)
}
