package helpers

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis/v8"
)

type RedisHelper struct {
	RedisClient *redis.Client
}

func NewRedisHelper(redisClient *redis.Client) RedisHelper {
	return RedisHelper{RedisClient: redisClient}
}

func (helper RedisHelper) IsKeyExist(key string) bool {
	count := helper.RedisClient.Exists(context.Background(), key).Val()
	if count == 0 {
		return false
	}

	return true
}

func (helper RedisHelper) QueueingOutlet(message *kafka.Message) {

	queue := "queue"
	key := "transactions"
	marshaledMessage := ToInterface(message.Value)

	tx := func(tx *redis.Tx) (err error) {

		// check is key exist
		count := tx.Exists(context.Background(), key).Val()

		if count == 0 { // if not exist

			// create data
			data := map[string]interface{}{
				queue: []interface{}{
					marshaledMessage,
				},
			}
			if err = tx.Set(context.Background(), key, ToByte(data), 0).Err(); err != nil {
				return err
			}

		} else { // if exist

			// get data
			lists, err := tx.Get(context.Background(), key).Bytes()
			if err != nil {
				return err
			}

			// marshalling lists & map insertion
			marshaledLists := ToInterface(lists)
			mappedLists := marshaledLists.(map[string]interface{})

			// append queue
			queues := mappedLists[queue].([]interface{})
			queues = append(queues, marshaledMessage)

			// insert into queue
			mappedLists[queue] = queues

			// insert into redis
			if err = tx.Set(context.Background(), key, ToByte(mappedLists), 0).Err(); err != nil {
				return err
			}

		}

		return err
	}

	if err := helper.RedisClient.Watch(context.Background(), tx, key); err != nil {
		fmt.Println("error redis transaction :", err.Error())
	}
}
