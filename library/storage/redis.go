package storage

import (
	"github.com/go-redis/redis"
)

type Redis struct {
	Address string

	Password string

	client *redis.Client
}

func NewRedis(address, password string) *Redis {
	rs := &Redis{
		Address:  address,
		Password: password,
	}

	rs.client = redis.NewClient(&redis.Options{
		Addr:     rs.Address,
		Password: rs.Password,
	})

	return rs
}

func (r *Redis) Set(key string, checksum []byte) error {
	cmd := r.client.Set(key, checksum, 0)

	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (r *Redis) Get(key string) ([]byte, error) {
	cmd := r.client.Get(key)

	err := cmd.Err()
	if err != nil {
		return nil, err
	}

	return cmd.Bytes()
}

func (r *Redis) Delete(key string) error {
	cmd := r.client.Del(key)

	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}
