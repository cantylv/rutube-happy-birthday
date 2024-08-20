package producer

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// SetupProduces - function that creates new produces for Apache Kafka Cluster
func Init() sarama.SyncProducer {
	logger := zap.Must(zap.NewProduction())
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	address := fmt.Sprintf("%s:%d", viper.GetString("kafka.host"), viper.GetUint16("kafka.port"))
	producer, err := sarama.NewSyncProducer([]string{address}, config)
	if err != nil {
		logger.Fatal(err.Error())
	}
	return producer
}
