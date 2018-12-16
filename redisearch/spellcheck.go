package redisearch

import (
	"strconv"

	"github.com/gomodule/redigo/redis"
)

type SpellCheck struct {
	Term 					string 			`json:"term"`
	Suggestions 	[]Suggs			`json:"suggestions"`
}

type Suggs struct {
	Suggestion 		string 			`json:"suggestion"`
	Distance 			float64 		`json:"distance"`
}

// SpellCheck, Performs spelling correction on a query,
// returning suggestions for misspelled terms
func (i *Client) SpellCheck(q *Query) (suggestions []SpellCheck, total int, err error) {
	conn := i.pool.Get()
	defer conn.Close()

	args := redis.Args{i.name}
	args = append(args, q.serialize()...)

	res, err := redis.Values(conn.Do("FT.SPELLCHECK", args...))
	if err != nil {
		return
	}

	// Initiate return document interface
	sugg := make([]SpellCheck, 0, len(res))
	// Prepare the document to be returned
		for _, v := range res {
			sugg = append(sugg, prepareSuggestions(v.([]interface{})))
		}

	return sugg, 1, err
}

// Prepare Spell Check suggestions
func prepareSuggestions(arr []interface{}) (SpellCheck) {
	sugg := SpellCheck{}
	sugg.Term = string(arr[1].([]uint8))
	
	for _, vv := range arr {
		switch vv.(type) {
		case []interface{}:
			s := Suggs{}
			for _, sv := range vv.([]interface{}) {
				s.Suggestion = string(sv.([]interface{})[1].([]uint8))
				sd, _ := strconv.ParseFloat(string(sv.([]interface{})[0].([]uint8)), 32)
				s.Distance = sd
				sugg.Suggestions = append(sugg.Suggestions, s)
			}
		}
	}
	
	return sugg
}