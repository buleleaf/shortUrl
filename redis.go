package main

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/mattheath/base62"
)

const (
	// URLIDKEY is golbal counter
	URLDKEY = "next.url.id"
	//
	ShortlinkKey       = "shortlink:%s:url"
	URLHashKey         = "urlhash:%s:url"
	ShortlinkDetailKey = "shortlink:%s:detail"
)

type RedisCli struct {
	Cli *redis.Client
}

type URLDetail struct {
	URL                 string        `json:"url"`
	CreateAt            string        `json:"create_at"`
	ExpirationInMinutes time.Duration `json:"expiration_in_minutes"`
}

func NewRedisCli(addr string, passwd string, db int) *RedisCli {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if _, err := rdb.Ping().Result(); err != nil {
		panic(err)
	}
	return &RedisCli{Cli: rdb}
}

func (r *RedisCli) Shorten(url string, exp int64) (string, error) {
	hash := toSha1(url)
	d, err := r.Cli.Get(fmt.Sprintf(URLHashKey, hash)).Result()
	// fmt.Printf("%v\n", fmt.Sprintf(URLHashKey, hash))
	// fmt.Printf("%v", d)

	if err == redis.Nil {
		// not existed, nothing to do
	} else if err != nil {
		return "", err
	} else {
		if d == "{}" {
		} else {
			return d, nil
		}
	}
	err = r.Cli.Incr(URLDKEY).Err()
	if err != nil {
		return "", err
	}

	// encode global counter to base64
	id, err := r.Cli.Get(URLDKEY).Int64()
	if err != nil {
		return "", nil
	}
	eid := base62.EncodeInt64(id)

	err = r.Cli.Set(fmt.Sprintf(ShortlinkKey, eid), url,
		time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}

	err = r.Cli.Set(fmt.Sprintf(URLHashKey, hash), eid, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}
	detail, err := json.Marshal(
		&URLDetail{
			URL:                 url,
			CreateAt:            time.Now().String(),
			ExpirationInMinutes: time.Duration(exp),
		})
	if err != nil {
		return "", err
	}
	fmt.Printf("%v", detail)
	err = r.Cli.Set(fmt.Sprintf(ShortlinkDetailKey, eid), detail,
		time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}

	return eid, nil
}

func (r *RedisCli) ShortlinkInfo(eid string) (interface{}, error) {
	d, err := r.Cli.Get(fmt.Sprintf(ShortlinkDetailKey, eid)).Result()
	if err == redis.Nil {
		return "", StatusError{404, errors.New("Unknown short URL")}
	} else if err != nil {
	} else {
		var detail interface{}
		if err := json.Unmarshal([]byte(d), &detail); err != nil {
			return "", err
		} else {
			return detail, nil
		}
	}
	return "", err
}

func (r *RedisCli) Unshorten(eid string) (string, error) {
	url, err := r.Cli.Get(fmt.Sprintf(ShortlinkKey, eid)).Result()
	if err == redis.Nil {
		return "", StatusError{404, err}
	} else if err != nil {
		return "", err
	} else {
		return url, nil
	}
}

func toSha1(str string) string {
	sha := sha1.New()
	return string(sha.Sum([]byte(str)))
}
