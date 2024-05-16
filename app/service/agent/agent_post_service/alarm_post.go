package agent_post_service

import (
	"context"
	"fmt"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/kafka_go"
	"goskeleton/app/utils/kafka_go/producer_kafka"
)

func CreateAlarmFactory() *AlarmPost {
	return &AlarmPost{
		queryClient: producer_kafka.CreateProducerFactory(),
	}
}

type AlarmPost struct {
	queryClient *producer_kafka.ProducerClient
}

// ToQueryBYComm 格式化参数后推送给队列
func (crd *AlarmPost) ToQueryBYComm(orgId int64, agentId, inputJson string) error {
	var params []alarmComm
	if err := json_helper.Json2Any(inputJson, &params); err != nil {
		return err
	}
	now := date_helper.CurDatetime()
	fmt.Println("params========================================")
	fmt.Println(json_helper.Encode(params))
	for _, item := range params {

	}

	return nil
}

// 发送进程数据
func (crd *AlarmPost) sendProcess(model *alarm_models.AlarmProcessModel) error {

	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeProcessPutMaster,
		Data: json_helper.Encode(&model),
	}
	crd.queryClient.ProduceProcess(context.Background(), msg)

	return nil
}
