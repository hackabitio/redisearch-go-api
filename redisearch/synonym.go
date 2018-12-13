package redisearch

import (	
	"github.com/gomodule/redigo/redis"
)


// Add synonym group to the index
func (i *Client) SynAdd(term, syn string) (bool, error) {
	conn := i.pool.Get()
	defer conn.Close()

	args := redis.Args{i.name}
	args = append(args, term, syn)

	err := conn.Send("FT.SYNADD", args...)

	if err != nil {
		return false, err
	}

	return true, nil
}

type SynGroup struct {
	Term		string
	Synonym	string
}
// Add synonym group to the index
func (i *Client) SynBulk(synGroup []SynGroup) (bool, error) {
	conn := i.pool.Get()
	defer conn.Close()

	args := redis.Args{i.name}

	for _, group := range synGroup {
		args = append(args, group.Term, group.Synonym)
	
		err := conn.Send("FT.SYNADD", args...)
	
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
