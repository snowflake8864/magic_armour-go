package gomail_factory

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/helpers/float_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"math"
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
	AlarmList  []AlarmCycleInfo
}

type AlarmCycleInfo struct {
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

type AlarmTimely struct {
	Date          string
	AgentId       string
	AgentIp       string
	AgentHostname string
	AlarmLogin    map[string]int //口令破解告警
}

func (k *Template) GetAlarmTimely(param *AlarmTimely) (subject, context string) {
	subject = variable.ConfigYml.GetString("AppCustomConfig.Email.AlarmTimely.Subject")
	title := variable.ConfigYml.GetString("AppCustomConfig.Email.AlarmTimely.Title")

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
    /* 中间告警信息开始 */
    .alarm-ul{
      margin-bottom: 30px;
    }
    .alarm-ul li {
      list-style: none;
      background: #FFFFFF;
      box-shadow: 0px 0px 8px 0px rgba(0,0,0,0.11);
      border-radius: 13px;
      padding: 18px 20px 15px;
      margin-bottom: 20px;
    }
    .alarm-ul li .card-tit{
      font-size: 14px;
      font-family: MicrosoftYaHei;
      color: #343434;
      line-height: 14px;
      line-height: 20px;
    }
    .alarm-ul li .card-tit img{
      float: left;
      margin: -2px 10px 0 0;
    }
    .alarm-ul .agent-div{
      margin: 15px 0 0 32px;
    }
    .alarm-ul .agent-div .agent-tit{
      font-weight: bold;
      color: #49505A;
      margin-bottom: 7px;
    }
    .alarm-ul .agent-div .pl{
      font-size: 12px;
      color: #343434;
      margin-bottom: 7px;
    }
    .alarm-ul .agent-div .pl p{
      padding-left: 47px;
    }
    .alarm-ul .agent-div .pl span{
      color: #A6A6A8;
      float: left;
    }
    /* 中间告警信息结束 */
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
	   <p class="tit">` + title + `</p>
       <ul class="alarm-ul">
`
	if len(param.AlarmLogin) > 0 {
		context += `
		 	<li>
          		<div class="card-tit">
					<img src="https://hj.veda.com/static/img/email_gauge2.png" width="23px" />
					您的终端在` + param.Date + `连续触发了口令破解告警！
				  </div>
				  <div class="agent-div">
					<div class="agent-tit">终端信息</div>
					<div class="pl"><span>IP地址：</span>` + param.AgentIp + `</div>
					<div class="pl"><span>主机名：</span>` + param.AgentHostname + `</div>
					<div class="pl"><span>设备ID：</span>` + param.AgentId + `</div>
					<div class="pl">
					  <span>攻击IP：</span>
`
		for ip, total := range param.AlarmLogin {
			context += `
					  <p>` + ip + `【攻击 ` + int_helper.IntToStr(total) + `次】</p>
`
		}

		context += `
					</div>
				  </div>
        	</li>
`
	}
	context += `
		</ul>
		<div class="danger"><span></span>请您及时处理安全风险。</div>
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

type HardwareTimely struct {
	Date          string
	AgentId       string
	AgentIp       string
	AgentHostname string
	UsageList     []*HardwareTimelyInfo
}

type HardwareTimelyInfo struct {
	Name      string
	Threshold int
	Usage     float64
}

func (k *Template) GetHardwareTimely(param *HardwareTimely) (subject, context string) {
	subject = variable.ConfigYml.GetString("AppCustomConfig.Email.HardwareTimely.Subject")
	title := variable.ConfigYml.GetString("AppCustomConfig.Email.HardwareTimely.Title")

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
    /* 中间表盘样式开始 */
    .gauge-ul{
      margin-bottom: 30px;
    }
    .gauge-ul li{
      list-style: none;
      width: 632px;
      height: 162px;
      background: #FFFFFF;
      box-shadow: 0px 0px 8px 0px rgba(0,0,0,0.11);
      border-radius: 13px;
      padding: 20px;
      margin-bottom: 20px;
    }
    .gauge-ul li .card-tit{
      font-size: 14px;
      font-family: MicrosoftYaHei;
      color: #343434;
      line-height: 14px;
    }
    .gauge-ul li .card-tit img{
      float: left;
      margin: -6px 10px 0 0;
    }
    .gauge-div{
      margin-top: 20px;
      padding: 0 16px;
    }
    .gauge-img{
      width: 146px;
      height: 121px;
      background: url('https://hj.veda.com/static/img/email_gauge.png') no-repeat;
      position: relative;
      float: left;
    }
    .gauge-img img{
      position: absolute;
      left: 50%;
      top: 32px;
      margin: 0 0 0 -4.5px;
      transform-origin: bottom;
      -ms-transform-origin: bottom;
      -webkit-transform-origin: bottom;
    }
    .gauge-top {
      font-size: 12px;
      color: #FF3B30;
      text-align: center;
      padding-top: 88px;
    }
    .gauge-bot{
      font-size: 12px;
      color: #999999;
      text-align: center;
      padding-top: 3px;
    }
    .gauge-txt{
      padding: 19px 0 0 172px;
    }
    .gauge-txt .title{
      font-size: 12px;
      font-weight: bold;
      color: #49505A;
      margin-bottom: 10px;
    }
    .gauge-txt p{
      font-size: 12px;
      color: #343434;
      margin-bottom: 7px;
    }
    .gauge-txt p span{
      color: #A6A6A8;
    }
    /* 中间表盘样式结束 */
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
	   <p class="tit">` + title + `</p>
       <ul class="gauge-ul">
`
	if len(param.UsageList) > 0 {
		for _, item := range param.UsageList {
			rate := float_helper.Add(float_helper.Mul(item.Usage, 2.59, 2), 233, 2)
			rotate := float_helper.Float64ToStr(math.Round(rate))

			context += `
		 	<li>
			  <p class="card-tit"><img src="https://hj.veda.com/static/img/email_gauge2.png" width="23px" />
` + param.Date + `终端触发了 ` + item.Name + ` ` + int_helper.IntToStr(item.Threshold) + `%</p>
			  <div class="gauge-div">
				<div class="gauge-img">
				  <img src="https://hj.veda.com/static/img/email_gauge1.png" alt="" width="9px"
style="transform: rotate(` + rotate + `deg); -ms-transform: rotate` + rotate + `deg); -webkit-transform: rotate(` + rotate + `deg);">
				  <p class="gauge-top">` + int_helper.IntToStr(item.Threshold) + `%</p>
				  <p class="gauge-bot">` + item.Name + `</p>
				</div>
				<div class="gauge-txt">
				  <div class="title">终端信息</div>
				  <p><span>IP地址：</span>` + param.AgentIp + `</p>
				  <p><span>主机名：</span>` + param.AgentHostname + `</p>
				  <p><span>设备ID：</span>` + param.AgentId + `</p>
				</div>
			  </div>
			</li>
`
		}
	}
	context += `
		</ul>
		<div class="danger"><span></span>请您及时处理安全风险。</div>
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
