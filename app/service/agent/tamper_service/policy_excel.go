package tamper_service

import (
	"errors"
	"fmt"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/enums/enum_tamper"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/model/agent/tamper_models"
	"goskeleton/app/service/upload_file"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/excel_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/map_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/md5_encrypt"
	"strings"
)

func CreatePolicyExcelFactory(orgId int64, agentId string, lang string) *PolicyExcel {
	return &PolicyExcel{
		masterModel: tamper_models.CreatePolicyFactory(),
		orgId:       orgId,
		agentId:     agentId,
		sheetName:   "Sheet1",
		sep:         "|",
		colsLen:     9,
		processMap:  make(map[string]string),
		lang:        lang,
	}
}

type PolicyExcel struct {
	masterModel *tamper_models.PolicyModel
	orgId       int64
	agentId     string
	sheetName   string
	sep         string
	colsLen     int
	processMap  map[string]string

	lang string //多语言支持
}

func (crd *PolicyExcel) generateExportPath() string {

	fullDir, _ := upload_file.GenerateYearMonthPath(upload_file.GetDefaultPath())
	name := md5_encrypt.MD5(fmt.Sprintf("%s_%d", crd.agentId, variable.SnowFlake.GetId()))

	return fullDir + name + ".xlsx"
}

func (crd *PolicyExcel) getPolicies() []*tamper_models.PolicyModel {
	condition := tamper_models.SearchConditions{
		OrgId:   crd.orgId,
		AgentId: crd.agentId,
	}
	return tamper_models.CreatePolicyFactory().GetListByConditions(condition)
}

func (crd *PolicyExcel) getTop() string {
	if crd.lang == consts.LanguageZh {
		return "" +
			"类型：单个文件或文件夹 如：文件夹\n" +
			"名称：当前文件或文件夹路径 如：D:\\Backup\\Downloads\n" +
			"保护中：读取|写入|删除|重命名|新建 如：读取|写入\n" +
			"包含文件类型：多个以竖线隔开 如: gif|jpg\n" +
			"排除文件类型：多个以竖线隔开 如: gif|jpg\n" +
			"排除文件夹：多个以竖线隔开 如: www|test\n" +
			"继承关系：0否 1是 如：是\n" +
			"信任进程：多个hash以竖线隔开如: aa|bb"
	} else {
		return "" +
			"Type: Single file or folder e.g., Folder\n" +
			"Name: Current file or folder path e.g., D:\\Backup\\Downloads\n" +
			"Protected Operations: Read|Write|Delete|Rename|Create e.g., Read|Write\n" +
			"Included File Types: Multiple separated by | e.g., gif|jpg\n" +
			"Excluded File Types: Multiple separated by | e.g., gif|jpg\n" +
			"Excluded Folders: Multiple separated by | e.g., www|test\n" +
			"Inheritance: 0=No 1=Yes e.g., Yes\n" +
			"Trusted Processes: Multiple hashes separated by | e.g., aa|bb"
	}
}

func (crd *PolicyExcel) getHeader() []interface{} {
	if crd.lang == consts.LanguageZh {
		return []interface{}{
			"类型",
			"名称",
			"保护中",
			"包含文件类型",
			"排除文件类型",
			"排除文件夹",
			"继承关系",
			"信任进程",
			"备注",
		}
	} else {
		return []interface{}{
			"Type",
			"Name",
			"Protected Operations",
			"Included File Types",
			"Excluded File Types",
			"Excluded Folders",
			"Inheritance",
			"Trusted Processes",
			"Remarks",
		}
	}
}

// Export 导出
func (crd *PolicyExcel) Export() (string, error) {
	//创建新 Excel 文件
	excelHandler := excel_helper.NewExcelHandler()
	defer excelHandler.Close()
	// 初始化默认样式
	_ = excelHandler.InitDefaultStyles()

	//顶部
	top := crd.getTop()
	if err := excelHandler.MergeCells(crd.sheetName, 1, 1, 1, crd.colsLen); err != nil {
		return "", fmt.Errorf("merge cells fail: %v", err)
	}
	if err := excelHandler.SetCellValue(crd.sheetName, "A1", top); err != nil {
		return "", fmt.Errorf("set top fail: %v", err)
	}
	_ = excelHandler.SetRowHeight(crd.sheetName, 1, 160)

	//头部
	header := crd.getHeader()
	if err := excelHandler.SetHeader(crd.sheetName, "A2", "I2", header); err != nil {
		return "", fmt.Errorf("set header fail: %v", err)
	}
	_ = excelHandler.SetColWidth(crd.sheetName, "A", "A", 20)
	_ = excelHandler.SetColWidth(crd.sheetName, "B", "B", 50)
	_ = excelHandler.SetColWidth(crd.sheetName, "C", "F", 30)
	_ = excelHandler.SetColWidth(crd.sheetName, "G", "G", 15)
	_ = excelHandler.SetColWidth(crd.sheetName, "H", "H", 50)
	_ = excelHandler.SetColWidth(crd.sheetName, "I", "I", 30)

	//策略数据
	models := crd.getPolicies()
	if len(models) >= 0 {
		var rowsData [][]interface{}
		for _, model := range models {

			var processes []tamper_models.Process
			_ = json_helper.ToAny(model.Processes, &processes)
			var processHashes []string
			if len(processes) > 0 {
				for _, proc := range processes {
					processHashes = append(processHashes, proc.Hash)
				}
			}

			rowsData = append(rowsData, []interface{}{
				enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_tamper.EnumMapInt),
				model.Dir,
				strings.Join(enum_tamper.GetPowerLabels(model.ProtectRw), crd.sep),
				model.IncludeSuffix,
				model.ExcludeSuffix,
				model.ExcludeFolder,
				enums.Get1stIsStrAnd2ndIsIntReturnStr("is_extend", model.IsExtend, enum_tamper.EnumMapInt),
				strings.Join(processHashes, crd.sep),
				model.Remark,
			})
		}

		if err := excelHandler.SetRowsValues(crd.sheetName, 1, 3, rowsData); err != nil {
			return "", fmt.Errorf("set rows fail: %v", err)
		}
	}

	// 保存文件
	path := crd.generateExportPath()
	if err := excelHandler.SaveAs(path); err != nil {
		return "", fmt.Errorf("save path fail: %v", err)
	}

	return path, nil
}

// Import 导入
func (crd *PolicyExcel) getProcesses(hashes string) string {
	result := []tamper_models.Process{}
	if string_helper.IsEmpty(hashes) {
		return json_helper.Encode(result)
	}

	//需求去查询的
	var queryList []string
	//转成数组格式
	hashList := strings.Split(hashes, crd.sep)
	for _, hash := range hashList {
		if val, isOK := crd.processMap[hash]; !isOK {
			queryList = append(queryList, hash)
			crd.processMap[hash] = ""
		} else if !string_helper.IsEmpty(val) {
			result = append(result, tamper_models.Process{
				Hash: hash,
				Dir:  val,
			})
		}
	}

	//查询进程白名单
	if len(queryList) > 0 {
		condition := proc_models.SearchConditions{
			OrgId:      crd.orgId,
			AgentId:    crd.agentId,
			ProcHashes: queryList,
			Status:     enum_process.ProcStatusWhite,
		}
		models := proc_models.CreateProcessFactory().GetListByConditions(condition)
		if models != nil && len(models) > 0 {
			for _, model := range models {
				crd.processMap[model.ProcHash] = model.Dir

				result = append(result, tamper_models.Process{
					Hash: model.ProcHash,
					Dir:  model.Dir,
				})
			}
		}
	}

	return json_helper.Encode(result)
}

// Import 导入
func (crd *PolicyExcel) Import(fullPath string, userId int64, userName string) error {
	//创建新 Excel 文件
	excelHandler, err := excel_helper.NewExcelHandlerFromFile(fullPath)
	if err != nil {
		return fmt.Errorf("read file: %v", err)
	}
	defer excelHandler.Close()

	// 获取工作表中指定单元格的值
	if cellVal, err1 := excelHandler.GetCellValue(crd.sheetName, "A2"); err1 != nil {
		return fmt.Errorf("read file: %v", err1)
	} else if cellVal != "类型" || cellVal != "Type" {
		return errors.New(my_errors.FilesFormatIncorrect)
	}

	//获取所有行数据
	rowsValues, err2 := excelHandler.GetRowsValues(crd.sheetName)
	if err2 != nil {
		return fmt.Errorf("get rows: %v", err)
	}
	if len(rowsValues) <= 2 {
		return nil
	}

	//获取现有策略
	policyMap := make(map[string]struct{})
	models := crd.getPolicies()
	if len(models) >= 0 {
		for _, model := range models {
			policyMap[model.Dir] = struct{}{}
		}
	}

	//属性映射
	typeMap := map_helper.FlipIntAndString(enums.Get1stIsStrReturnMapIntOfStr("type", enum_tamper.EnumMapInt))
	powerMap := map_helper.FlipIntAndString(enums.Get1stIsStrReturnMapIntOfStr("power_folder", enum_tamper.EnumMapInt))
	extendMap := map_helper.FlipIntAndString(enums.Get1stIsStrReturnMapIntOfStr("is_extend", enum_tamper.EnumMapInt))
	now := date_helper.CurDatetime()
	var insList []tamper_models.PolicyModel
	//处理行数据
	for i, rowValues := range rowsValues {
		//忽略前2行
		if i < 2 {
			continue
		}
		//补足
		if colRealLen := len(rowValues); colRealLen < crd.colsLen {
			tmp := make([]string, crd.colsLen-colRealLen)
			rowValues = append(rowValues, tmp...)
		}

		//初始化数据结构
		model := tamper_models.PolicyModel{
			Id:            variable.SnowFlake.GetId(),
			OrgId:         crd.orgId,
			AgentId:       crd.agentId,
			IncludeSuffix: rowValues[3],
			ExcludeSuffix: rowValues[4],
			ExcludeFolder: rowValues[4],
			OpenSwitch:    enums.SwitchOff,
			Remark:        rowValues[8],
			CreateId:      userId,
			CreateName:    userName,
			CreatedAt:     now,
			UpdateId:      userId,
			UpdateName:    userName,
			UpdatedAt:     now,
		}

		// 类型是否正确
		if val, isOk := typeMap[rowValues[0]]; isOk {
			model.Type = val
		} else {
			continue
		}
		//路径是否重复
		dir := path_helper.ToStandard(rowValues[1])
		if _, isOk := policyMap[dir]; isOk {
			continue
		} else {
			model.Dir = dir
		}
		// 保护权限是否正确
		protectRw := 0
		for _, val := range strings.Split(rowValues[2], crd.sep) {
			if v, isOk := powerMap[val]; isOk {
				protectRw += v
			}
		}
		if protectRw == 0 {
			continue
		} else {
			model.ProtectRw = protectRw
		}
		// 继承是否正确
		if val, isOk := extendMap[rowValues[6]]; isOk {
			model.IsExtend = val
		} else {
			continue
		}
		//信任进程
		model.Processes = crd.getProcesses(rowValues[7])

		insList = append(insList, model)
	}

	if len(insList) > 0 {
		if !crd.masterModel.InsertBatch(&insList) {
			return errors.New("batch insert fail")
		}
	}

	return nil
}
