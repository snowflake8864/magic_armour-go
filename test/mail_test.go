package test

import (
	"fmt"
	"goskeleton/app/utils/gomail_factory"
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
