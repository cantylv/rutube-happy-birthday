package consumer

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/IBM/sarama"
	"github.com/cantylv/service-happy-birthday/internal/entity"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Init() sarama.Consumer {
	logger := zap.Must(zap.NewProduction())
	cons, err := setupConsumer()
	if err != nil {
		logger.Fatal("error while set up kafka consumer", zap.Error(err))
	}
	go runEngine(cons)
	return cons
}

func setupConsumer() (sarama.Consumer, error) {
	kafkaAddress := fmt.Sprintf("%s:%d", viper.GetString("kafka.host"), viper.GetUint16("kafka.port"))
	consumer, err := sarama.NewConsumer([]string{kafkaAddress}, sarama.NewConfig())
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

func runEngine(cons sarama.Consumer) {
	logger := zap.Must(zap.NewProduction())
	partitions, err := cons.Partitions(viper.GetString("kafka.topic"))
	if err != nil {
		logger.Fatal(err.Error())
	}

	// Открываем файл для записи (если файл не существует, он будет создан)
	file, err := os.Create("./notification.json")
	if err != nil {
		logger.Error("error creating file:", zap.Error(err))
		return
	}
	defer file.Close()

	var wg sync.WaitGroup
	var mtx sync.RWMutex
	for _, partition := range partitions {
		partitionConsumer, err := cons.ConsumePartition(viper.GetString("kafka.topic"), partition, sarama.OffsetNewest)
		if err != nil {
			logger.Fatal(err.Error())
		}
		wg.Add(1)
		go listen(partitionConsumer, file, &wg, &mtx)
	}
	wg.Wait()
}

func listen(partitionConsumer sarama.PartitionConsumer, file *os.File, wgParent *sync.WaitGroup, mtx *sync.RWMutex) {
	logger := zap.Must(zap.NewProduction()).Sugar()
	logger.Infof("Listen kafka topic %s", viper.GetString("kafka.topic"))
	defer wgParent.Done()
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		var segment entity.Notification
		err := json.Unmarshal(msg.Value, &segment)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		fmt.Println("NOTIFICATION RECEIVE")
		fmt.Println(segment)

		jsonData, err := json.MarshalIndent(segment, "", "    ")
		if err != nil {
			logger.Error("error marshalling data:", zap.Error(err))
			continue
		}

		mtx.Lock()
		_, err = file.Write(jsonData)
		mtx.Unlock()
		if err != nil {
			logger.Error("error writing to file:", zap.Error(err))
			return
		}
	}
}
