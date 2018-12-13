package redisearch

import (
	"fmt"
	
	"github.com/gomodule/redigo/redis"
)


// Add synonym group to the index
func (i *Client) SynAdd(term, syn string) (bool, error) {
	conn := i.pool.Get()
	defer conn.Close()

	args := redis.Args{i.name}
	args = append(args, term, syn)

	res, err := conn.Do("FT.SYNADD", args...)
	fmt.Printf("Type: %T, value: %v", res, res)
	if err != nil {
		return false, err
	}

	return true, nil
}
