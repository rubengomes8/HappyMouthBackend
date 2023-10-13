package recipe

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

type RecipeEvent struct {
	Content string `json:"content"`
}

func ProduceRecipeEvent(producer sarama.SyncProducer, key, content, topic string) (int32, int64, error) {
	recipeEvent := RecipeEvent{
		Content: content,
	}
	eventBytes, err := json.Marshal(&recipeEvent)
	if err != nil {
		return 0, 0, err
	}

	event := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(eventBytes),
	}

	partition, offset, err := producer.SendMessage(event)
	if err != nil {
		return 0, 0, err
	}

	return partition, offset, nil
}
