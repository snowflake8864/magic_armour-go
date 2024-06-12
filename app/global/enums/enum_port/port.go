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
)

var EnumMapInt = map[string]map[int]string{
	"alarm_level": {
		AlarmLevel1: "1级",
		AlarmLevel2: "2级",
		AlarmLevel3: "3级",
		AlarmLevel4: "4级",
		AlarmLevel5: "5级",
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
