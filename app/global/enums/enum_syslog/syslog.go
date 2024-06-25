package enum_syslog

import (
	"goskeleton/app/global/enums"
	"strings"
)

const (
	LogTypeDns     int = 1 //日志类型：dns
	LogTypeOuter   int = 2 //日志类型：对外通讯
	LogTypeInner   int = 3 //日志类型：对内通讯
	LogTypeProcess int = 4 //日志类型：进程启动
	LogTypeLogin   int = 5 //日志类型：登录

	LoginRdp string = "RDP" //登录类型： 远程桌面协议
	LoginSmb string = "SMB" //登录类型： 一种局域网文件共享传输协议
	LoginSsh string = "SSH" //登录类型： 一种网络协议，用于计算机之间的加密登录
)

var EnumMapInt = map[string]map[int]string{
	"login_status": {
		enums.IsYes: "成功",
		enums.IsNot: "失败",
	},
	"port_status": {
		enums.IsYes: "是",
		enums.IsNot: "否",
	},
	"proto": {
		0:   "ANY",
		1:   "ICMP",
		6:   "TCP",
		17:  "UDP",
		13:  "ARGUS",
		51:  "AH",
		61:  "any host internal protocol",
		63:  "any locol network",
		68:  "any distributed file system",
		93:  "AX.25",
		99:  "any private encryp scheme",
		104: "ARIS",
		107: "A/N",
		114: "any 0-hop protocol",
		10:  "BBN-RCC-MON",
		49:  "BNA",
		76:  "BR-SAT-MON",
		7:   "CBT",
		16:  "CHAOS",
		72:  "CPNX",
		73:  "CPHB",
		110: "Compaq-Peer",
		62:  "CFTP",
		126: "CRTP",
		127: "CRUDP",
		116: "DDX",
		19:  "DCN-MEAS",
		33:  "DCCP",
		37:  "DDP",
		86:  "DGP",
		8:   "EGP",
		14:  "EMCON",
		50:  "ESP",
		97:  "ETHERIP",
		98:  "ENCAP",
		88:  "EIGRP",
		125: "FIRE",
		133: "FC",
		47:  "GRE",
		3:   "GGP",
		100: "GMTP",
		20:  "HMP",
		4:   "IP",
		2:   "IGMP",
		9:   "IGP",
		28:  "IRTP",
		29:  "ISO-TP4",
		35:  "IDPR",
		40:  "IL",
		43:  "IPv6-Route",
		44:  "IPv6-Frag",
		45:  "IDRP",
		41:  "IPv6",
		52:  "I-NLSP",
		58:  "IPv6-ICMP",
		59:  "IPv6-NoNxt",
		60:  "IPv6-Opts",
		67:  "IPPC",
		71:  "IPCV",
		80:  "ISO-IP",
		94:  "IPIP",
		101: "IFMP",
		108: "IPComp",
		111: "IPX-in-IP",
		117: "IATP",
		124: "ISISover",
		129: "IPLT",
		38:  "IDPR-CMTP",
		65:  "KRYPTOLAN",
		25:  "LEAF-1",
		26:  "LEAF-2",
		115: "L2TP",
		91:  "LARP",
		18:  "MUX",
		31:  "MFE-NSP",
		32:  "MERIT-INP",
		48:  "MHRP",
		55:  "MOBILE",
		92:  "MTP",
		95:  "MICP",
		135: "MobilityHeader",
		137: "MPLS-in-IP",
		30:  "NETBLT",
		54:  "NARP",
		11:  "NVP-II",
		85:  "NSFNET-IGP",
		89:  "OSPFIGP",
		12:  "PUP",
		21:  "PRM",
		75:  "PVP",
		102: "PNNI",
		103: "PIM",
		113: "PGM",
		123: "PTP",
		131: "PIPE",
		106: "QNX",
		27:  "RDP",
		46:  "RSVP",
		66:  "RVD",
		134: "RSVP-E2E-IGNORE",
		5:   "ST",
		42:  "SDRP",
		53:  "SWIPE",
		57:  "SKIP",
		64:  "SAT-EXPAK",
		69:  "SAT-MON",
		82:  "SECURE-VMTP",
		90:  "Sprite-RPC",
		96:  "SCC-SP",
		105: "SCPS",
		77:  "SUN-ND",
		109: "SNP",
		118: "STP",
		119: "SRP",
		121: "SMP",
		122: "SM",
		128: "SSCOPMCE",
		130: "SPS",
		132: "SCTP",
		23:  "TRUNK-1",
		24:  "TRUNK-2",
		39:  "TP++",
		56:  "TLSP",
		84:  "TTP",
		87:  "TCF",
		120: "UTI",
		136: "UDPLite",
		70:  "VISA",
		81:  "VMTP",
		83:  "VINES",
		112: "VRRP",
		74:  "WSN",
		78:  "WB-MON",
		79:  "WB-EXPAK",
		15:  "XNET",
		22:  "XNS-IDP",
		36:  "XTP",
		34:  "3PC"},
}

var WebMapInt = map[string][]int{}

// CheckLoginType 校验登录类型
func CheckLoginType(s string) bool {

	arr := map[string]struct{}{
		LoginRdp: {},
		LoginSmb: {},
		LoginSsh: {},
	}

	s = strings.ToUpper(s)
	_, iSOK := arr[s]
	return iSOK
}
