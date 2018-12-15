package redisearch

import (
	"github.com/gomodule/redigo/redis"
)

// Returns the full contents of a document.
func (i *Client) Get(id string) (doc map[string]interface{}, err error) {
	conn := i.pool.Get()
	defer conn.Close()

	args := redis.Args{i.name}
	args = append(args, id)

	res, err := redis.Values(conn.Do("FT.GET", args...))
	if err != nil {
		return nil, err
	}
	// Initiate return document interface
	d := map[string]interface{}{}
	// Prepare the document to be returned
	for i := 0; i < len(res); i += 2 {
		prop := string(res[i].([]byte))
		var val interface{}
		switch v := res[i+1].(type) {
		case []byte:
			val = string(v)
		default:
			val = v
		}
		// Assign fields to their values
		d[prop] = val
	}
	
	return d, nil
}