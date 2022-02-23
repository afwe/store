package redistest

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func ExampleClient() {
	resp, err := http.Get("http://qt.gtimg.cn/q=sh600660")
	stringresp, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Print(string(stringresp))
	rdb := redis.NewClient(&redis.Options{
		Addr:     "123.60.208.60:6379",
		Password: "",
		DB:       0,
	})
	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	val3, err := rdb.SetEX(ctx, "key2", 0, time.Second*60).Result()

	if err == redis.Nil {
		fmt.Println("Nil")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("val3", val3)
	}
	err = rdb.RPush(ctx, "list", "1").Err()
	err = rdb.LPush(ctx, "list", "2").Err()
	err = rdb.LSet(ctx, "list", 0, "3").Err()
	Len, err := rdb.LLen(ctx, "list").Result()
	fmt.Println(Len)
	lists, err := rdb.LRange(ctx, "list", 0, Len-1).Result()
	fmt.Println(lists)
	result, err := rdb.BLPop(ctx, time.Second, "list").Result()
	fmt.Println(result)
	lists, err = rdb.LRange(ctx, "list", 0, Len-1).Result()
	fmt.Println(lists)
}
