package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_port"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/port_models"
	"goskeleton/app/service/agent/agent_post_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
	"strings"
)

// PortController 端口相关
type PortController struct {
	ControllerBase
}

// PutWorking 上传业务端口
func (u *PortController) PutWorking(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//关闭任务
	tasks := []int{enum_task.TypePortWorkingUpload}
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "service")
	err := agent_post_service.CreatePortFactory().SaveWorking(loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)

}

// GetVirtual 下载虚开端口
func (u *PortController) GetVirtual(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//关闭任务
	tasks := []int{enum_task.TypePortVirtualPolicyDownload}
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	//win xp 不支持的功能
	if loginInfo.OsType == enum_software.OsTypeWinXp {
		response.SuccessSimplify(ctx)
		return
	}
	openSwitch := enums.SwitchOn
	result := []*PortResp{} //需初始化
	//全局策略
	globalModels := port_models.CreateVirtualGlobalFactory().GetListByConditions(port_models.SearchConditions{
		OrgId:      loginInfo.OrgId,
		DeptIds:    u.getDeptIdsByAgent(loginInfo.OrgId, loginInfo.AgentId),
		OpenSwitch: &openSwitch,
	})
	if globalModels != nil && len(globalModels) > 0 {
		for _, model := range globalModels {
			policies := u.format(model.Type, model.Protocol, model.SourceIp, model.SourcePort,
				model.ForwardIp, model.ForwardPort, model.ForwardType, model.AlarmLevel)
			for _, policy := range policies {
				result = append(result, policy)
			}
		}
	}

	//终端策略
	policyModels := port_models.CreateVirtualPolicyFactory().GetListByConditions(port_models.SearchConditions{
		OrgId:      loginInfo.OrgId,
		AgentId:    loginInfo.AgentId,
		OpenSwitch: &openSwitch,
	})
	if policyModels != nil && len(policyModels) > 0 {
		for _, model := range policyModels {
			policies := u.format(model.Type, model.Protocol, model.SourceIp, model.SourcePort,
				model.ForwardIp, model.ForwardPort, model.ForwardType, model.AlarmLevel)
			for _, policy := range policies {
				result = append(result, policy)
			}
		}
	}

	response.SuccessReturnData(ctx, result)
}

type PortResp struct {
	Id             int    `json:"id"`
	Type           string `json:"type"`
	Protocol       string `json:"protocol"`
	SourceIp       string `json:"source_ip"`
	SourcePort     string `json:"source_port"`
	SourcePortType int    `json:"source_port_type"` //虚拟端口类型：1单个端口 2范围
	ForwardIp      string `json:"dest_ip"`
	ForwardPort    string `json:"dest_port"`
	ForwardType    int    `json:"dest_port_type"`
	AlarmLevel     int    `json:"alarm_level"`
}

func (u *PortController) format(pType, protocol, sourceIp, sourcePort, forwardIp, forwardPort string, forwardType, alarmLevel int) (list []*PortResp) {

	//源端口:单端口(80);多端口(80,90);范围(80-90)
	sourcePortArr := strings.Split(sourcePort, ",")
	for _, sPort := range sourcePortArr {
		sourcePortType := 1
		if strings.Contains(sPort, "-") {
			sourcePortType = 2
		}
		if string_helper.IsEmpty(forwardIp) {
			forwardPort = ""
		}
		if forwardType == enum_port.ForwardTypeOriginal {
			forwardPort = ""
		}

		list = append(list, &PortResp{
			Type:           pType,
			Protocol:       protocol,
			SourceIp:       sourceIp,
			SourcePort:     sPort,
			ForwardIp:      forwardIp,
			ForwardPort:    forwardPort,
			ForwardType:    forwardType,
			AlarmLevel:     alarmLevel,
			SourcePortType: sourcePortType,
		})
	}

	return list
}
