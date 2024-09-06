package gomail_factory

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/helpers/float_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateTemplateFactory() *Template {
	return &Template{
		ImgHead:      variable.ConfigYml.GetString("AppCustomConfig.Email.ImgHead"),
		ImgFoot:      variable.ConfigYml.GetString("AppCustomConfig.Email.ImgFoot"),
		ImgWechat:    variable.ConfigYml.GetString("AppCustomConfig.Email.ImgWechat"),
		FootAbout:    variable.ConfigYml.GetString("AppCustomConfig.Email.FootAbout"),
		FootInfo:     variable.ConfigYml.GetString("AppCustomConfig.Email.FootInfo"),
		FootLinkHref: variable.ConfigYml.GetString("AppCustomConfig.Email.FootLinkHref"),
		FootLinkTxt:  variable.ConfigYml.GetString("AppCustomConfig.Email.FootLinkTxt"),
	}
}

// 结构体
type Template struct {
	ImgHead      string
	ImgFoot      string
	ImgWechat    string
	FootAbout    string
	FootInfo     string
	FootLinkHref string
	FootLinkTxt  string
}

type AlarmCycle struct {
	DateStart  string
	DateEnd    string
	AlarmTotal int
	AlarmList  []*AlarmCycleAlarm
}

type AlarmCycleAlarm struct {
	Name string
	Add  int
}

func (k *Template) GetAlarmCycle(param *AlarmCycle) (subject, context string) {
	subject = variable.ConfigYml.GetString("AppCustomConfig.Email.AlarmCycle.Subject")
	titleNoData := variable.ConfigYml.GetString("AppCustomConfig.Email.AlarmCycle.TitleNoData")
	titleHaveData := variable.ConfigYml.GetString("AppCustomConfig.Email.AlarmCycle.TitleHaveData")

	//head + <body><div class="center">
	context = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title></title>
	<style type="text/css">
    *{
      margin: 0;
      padding: 0;
    }
		html,body {
      color: #343434;
      font-size: 14px;
    }
    .center{
      width: 720px;
      margin: 0 auto;
      background: #fff;
    }
    .txt-c{
      margin: 0 24px;
    }
    .tit{
      font-size: 16px;
      font-weight: bold;
      padding: 30px 0;
      margin-bottom: 24px;
      border-bottom: 1px solid #EBEBEB;
    }
    .blue{
      color: #267EFF;
    }
    .ul{
      width: 100%;
      overflow: hidden;
    }
    .ul .li{
      width: 100%;
      list-style: none;
      line-height: 24px;
      margin: 24px 0;
      overflow: hidden;
    }
    .line-bg{
      float: left;
      width: 66%;
      height: 12px;
      background: #EFF3F6;
      border-radius: 8px;
      margin-top: 7px;
    }
    .total{
      float: right;
      font-size: 12px;
    }
    .num{
      float: left;
      width: 24px;
      height: 24px;
      line-height: 24px;
      text-align: center;
      background: url('https://hj.veda.com/static/img/email_li-bg.png') no-repeat center;
      margin-right: 16px;
    }
    .line{
      height: 12px;
      background: #267EFF;
      border-radius: 8px;
    }
    .txt{
      float: left;
      width: 100px;
      border: none;
      background: none;
    }
    .danger{
      width: 100%;
      height: 40px;
      line-height: 40px;
      text-align: center;
      background: #fee8ea;
      font-weight: bold;
      color: #F5222D;
      margin-bottom: 38px;
    }
    .danger span{
      display: inline-block;
      width: 16px;
      height: 15px;
      background: url('https://hj.veda.com/static/img/email_de.png') no-repeat;
      margin-right: 12px;
    }
    .about{
      font-size: 16px;
      font-weight: bold;
      margin: 44px 0 18px;
      text-align: center;
    }
    .info{
      text-indent: 24px;
      line-height: 20px;
    }
    .link{
      text-align: center;
      margin: 32px 0 44px;
    }
    .link a{
      display: inline-block;
      width: 160px;
      height: 40px;
      line-height: 40px;
      text-align: center;
      background: #267EFF;
      border-radius: 4px;
      font-size: 12px;
      color: #FFFFFF;
      text-decoration: none;
    }
    .nodata{
      display: block;
      text-align: center;
      margin: 39px 0;
    }
	</style>
</head>
<body>
  <div class="center">
`

	if !string_helper.IsEmpty(k.ImgHead) {
		context += `
  	<img src="` + k.ImgHead + `" alt="" width="720px">
`
	}

	/********* 内容 S *********/
	context += `
  	<div class="txt-c">
`
	if param.AlarmTotal == 0 {
		context += `
 		<p class="tit">` + titleNoData + `</p>
	    <p>` + param.DateStart + ` 至 ` + param.DateEnd + ` 无安全问题</p>
		<img src="https://hj.veda.com/static/img/email_no.png" alt="">
		<b class="blue nodata">暂无风险</b>
`
	} else {
		context += `
 		<p class="tit">` + titleHaveData + `</p>
		<p>` + param.DateStart + ` 至 ` + param.DateEnd + ` 内产生了 <span class="blue">` + int_helper.IntToStr(param.AlarmTotal) + `</span> 条告警：</p>
		<div class="ul">
`
		for key, alarm := range param.AlarmList {
			context += `
			 <div class="li">
				<span class="num">` + int_helper.IntToStr(key+1) + `</span>
				<span class="txt">` + alarm.Name + `</span>
				<div class="line-bg"><div class="line" style="width: ` + float_helper.Float64ToStr(float_helper.Div(alarm.Add, param.AlarmTotal, 2)) + `%;"></div></div>
				<span class="total">` + int_helper.IntToStr(alarm.Add) + `</span>
			</div>
`
		}
		context += `
		</div>
		<div class="danger"><span></span>请您及时处理安全风险。</div>
`
	}
	context += `
 	 </div>
`
	/********* 内容 E *********/

	if !string_helper.IsEmpty(k.ImgFoot) {
		context += `
  	<img src="` + k.ImgFoot + `" alt="" width="720px">
`
	}

	if !string_helper.IsEmpty(k.FootAbout) {
		context += `
	 <div class="txt-c">
      <p class="about">` + k.FootAbout + `</p>
      <p class="info">` + k.FootInfo + `</p>
      <div class="link">
        <a href="` + k.FootLinkHref + `" target="_blank" style="color:#fff;text-decoration: none;">` + k.FootLinkTxt + `</a>
      </div>
`
		if !string_helper.IsEmpty(k.ImgWechat) {
			context += `
      <div style='text-align: center'>
        <img src="` + k.ImgWechat + `" alt="" width="160px">
      </div>
      <p style="text-align: center; margin-top: 10px;">扫一扫</p>
      <p style="text-align: center; padding-bottom: 50px;">关注微信公众号</p>
`
		}
		context += `
    </div>
`
	}

	//head + <body><div class="center"> 的呼应
	context += `
  </div>
</body>
</html>
`
	return subject, context
}
