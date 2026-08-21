package enum_port

const (
	IpTypeIpv4 string = "ipv4" //Ip类型
	IpTypeIpv6 string = "ipv6" //Ip类型

	ProtocolTcp string = "tcp" //协议

	AlarmLevel1 int = 1 //告警权重
	AlarmLevel2 int = 2 //告警权重
	AlarmLevel3 int = 3 //告警权重
	AlarmLevel4 int = 4 //告警权重
	AlarmLevel5 int = 5 //告警权重

	Ipv4BehalfAll       string = "0.0.0.0"                                //代表所有ipv4地址
	Ipv4BehalfLocal     string = "127.0.0.1"                              //Ipv4的回路地址
	Ipv6BehalfAll       string = "::"                                     //代表所有ipv6地址
	Ipv6BehalfLocalFull string = "000:0000:0000:0000:0000:0000:0000:0001" //Ipv6的回路完整地址
	Ipv6BehalfLocalAbbr string = "::1"                                    //Ipv6的回路简写地址

	ForwardTypeFixed    int = 0 //转发端口类型: 固定端口
	ForwardTypeOriginal int = 1 //转发端口类型: 原端口
)

var EnumMapInt = map[string]map[int]string{
	"alarm_level": {
		AlarmLevel1: "1级",
		AlarmLevel2: "2级",
		AlarmLevel3: "3级",
		AlarmLevel4: "4级",
		AlarmLevel5: "5级",
	},
	"alarm_level_en": {
		AlarmLevel1: "Level 1",
		AlarmLevel2: "Level 2",
		AlarmLevel3: "Level 3",
		AlarmLevel4: "Level 4",
		AlarmLevel5: "Level 5",
	},
	"forward_type": {
		ForwardTypeFixed:    "指定端口",
		ForwardTypeOriginal: "原端口转发",
	},
	"forward_type_en": {
		ForwardTypeFixed:    "Fixed Port",
		ForwardTypeOriginal: "Original Port Forwarding",
	},
}

var EnumMapStr = map[string]map[string]string{
	"ip_type": {
		IpTypeIpv4: "ipv4",
		IpTypeIpv6: "ipv6",
	},
	"protocol": {
		ProtocolTcp: "tcp",
	},
}

var WebMapInt = map[string][]int{
	"alarm_level": {
		AlarmLevel1,
		AlarmLevel2,
		AlarmLevel3,
		AlarmLevel4,
		AlarmLevel5,
	},
	"forward_type": {
		ForwardTypeFixed,
		ForwardTypeOriginal,
	},
}

var WebMapStr = map[string][]string{
	"ip_type": {
		IpTypeIpv4,
		IpTypeIpv6,
	},
	"protocol": {
		ProtocolTcp,
	},
}
