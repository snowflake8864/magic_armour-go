package threaten_service

import (
	"context"
	"errors"
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/model/agent/threaten_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/excel_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/map_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/kafka_go"
	"goskeleton/app/utils/kafka_go/producer_kafka"
	"strings"
)

func CreateImportFactory() *ThreatenImport {
	return &ThreatenImport{
		masterModel: threaten_models.CreateCacheFactory(),
		sheetName:   "Worksheet",
		colsLen:     5,
	}
}

type ThreatenImport struct {
	masterModel *threaten_models.CacheModel
	sheetName   string
	colsLen     int
}

// ParseExcel 解析Excel
func (crd *ThreatenImport) ParseExcel(filename string) error {
	//创建新 Excel 文件
	excelHandler, err := excel_helper.NewExcelHandlerFromFile(filename)
	if err != nil {
		return fmt.Errorf("read file: %v", err)
	}
	defer excelHandler.Close()

	// 获取工作表中指定单元格的值
	if cellVal, err1 := excelHandler.GetCellValue(crd.sheetName, "A1"); err1 != nil {
		return fmt.Errorf("read file: %v", err1)
	} else if cellVal != "HASH" {
		return errors.New("文件格式不正确，不能导入")
	}
	if cellVal, err1 := excelHandler.GetCellValue(crd.sheetName, "B1"); err1 != nil {
		return fmt.Errorf("read file: %v", err1)
	} else if cellVal != "类型" {
		return errors.New("文件格式不正确，不能导入")
	}

	//获取所有行数据
	rowsValues, err2 := excelHandler.GetRowsValues(crd.sheetName)
	if err2 != nil {
		return fmt.Errorf("get rows: %v", err)
	}
	if len(rowsValues) <= 1 {
		return nil
	}

	//属性映射
	levelMap := map_helper.FlipIntAndString(enums.Get1stIsStrReturnMapIntOfStr("level", enum_threaten.EnumMap))
	now := date_helper.CurDatetime()
	var insList []threaten_models.CacheModel
	var rewriteList []threaten_models.CacheModel

	//处理行数据
	for i, rowValues := range rowsValues {
		//忽略前1行
		if i < 1 {
			continue
		}
		//补足
		if colRealLen := len(rowValues); colRealLen < crd.colsLen {
			tmp := make([]string, crd.colsLen-colRealLen)
			rowValues = append(rowValues, tmp...)
		}

		//初始化数据结构
		model := threaten_models.CacheModel{
			Source:      enum_threaten.SourceImport,
			DataFrom:    enum_threaten.DataFromUnknown,
			CheckAt:     now,
			CheckSource: enum_threaten.CheckSourceVeda,
			CreatedAt:   now,
		}

		// hash不能为空
		if string_helper.IsEmpty(rowValues[0]) {
			continue
		} else {
			model.Id = strings.ToLower(rowValues[0])
		}
		// 类型是否正确
		if val, isOk := levelMap[rowValues[1]]; isOk {
			model.Level = val
		} else {
			continue
		}

		insList = append(insList, model)
		if enum_threaten.IsWhite(model.Level) || enum_threaten.IsBlack(model.Level) {
			rewriteList = append(rewriteList, model)
		}
	}

	//先处理要插入的
	total := len(insList)
	if total == 0 {
		return nil
	}
	// 循环切割
	chunkSize := 500
	for i := 0; i < total; i += chunkSize {
		end := i + chunkSize
		if end > total {
			end = total
		}

		chunks := insList[i:end]
		if !crd.masterModel.CustomizeImportInsert(&chunks) {
			return errors.New("写入失败")
		}
	}

	//处理需要重写
	if len(rewriteList) > 0 {
		ctx := context.Background()
		producerFact := producer_kafka.CreateProducerFactory()
		for _, model := range rewriteList {
			producerFact.ProduceThreaten(ctx, &kafka_go.MessageStruct{
				Type: kafka_go.TypeThreatenSyncRewrite,
				Data: json_helper.Encode(model),
			})
		}

	}

	return nil
}
