package test

import (
	"fmt"
	_ "github.com/segmentio/kafka-go"
	"goskeleton/app/global/variable"
	_ "goskeleton/bootstrap"
	"testing"
)

func TestConfig(t *testing.T) {
	fmt.Println(variable.ConfigKafkaV2Yml.GetString("KafkaV2.Producer.Addr"))
}
