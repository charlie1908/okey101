package Core

//go get github.com/go-redis/redis/
//go get github.com/go-redis/redis/v8
//go clean -i github.com/go-redis/redis/
//go get github.com/ulule/limiter/v3
//go get github.com/ulule/limiter/v3/drivers/store/redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	redisStore "github.com/ulule/limiter/v3/drivers/store/redis"
	"log"
	"okey101/Model"
	shared "okey101/Shared"
	"time"
)

type RedisClient struct {
	c *redis.Client
}

var client = &RedisClient{}
var ctx = context.Background()

func GetRedisClient() *RedisClient {

	//Test SqlDatabase
	//database, _ := Decrypt(shared.Config.SQLURL, shared.Config.SECRETKEY)
	//encdatabase, _ := Encrypt("sqlserver://sa:AspNet55@192.168.50.173:1433?database=Customer&encrypt=disable&connection+timeout=30", shared.Config.SECRETKEY)
	//fmt.Println(database)
	//fmt.Println(encdatabase)

	redisPassword, _ := Decrypt(shared.Config.REDISPASSWORD, shared.Config.SECRETKEY)
	redisUrl, _ := Decrypt(shared.Config.REDISURL, shared.Config.SECRETKEY)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPassword, // no password set
		DB:       0,             // use default DB
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("Unable to connect to redis " + err.Error())
	}
	client.c = rdb
	return client
}

func (client *RedisClient) SetKey(key string, value interface{}, expiration time.Duration) error {
	cacheData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = client.c.Set(ctx, key, cacheData, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (client *RedisClient) GetKey(key string, src interface{}) error {
	val, err := client.c.Get(ctx, key).Result()
	if err == redis.Nil || err != nil {
		return err
	}

	err = json.Unmarshal([]byte(val), &src)
	if err != nil {
		return err
	}
	return nil
}

func (client *RedisClient) GetLifeTimeMinutes(key string) (int64, error) {
	ttl, err := client.c.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if ttl <= 0 {
		return 0, nil // Key yoksa ya da süresizse 0 döndür
	}
	return int64(ttl.Minutes()), nil
}

func GenerateRedisKey(userName string, isRefreshToken bool) string {
	if isRefreshToken {
		return userName + ":RefreshToken"
	}
	return userName + ":Token"
}

// Oda genel durumu (oyunun tamamı için global durum)
func GenerateRoomStateRedisKey(roomID string) string {
	return fmt.Sprintf("room:%s:game_state", roomID)
}

// Oyuncunun özel (gizli) durumu – sadece kendi elindeki taşlar
func GeneratePlayerPrivateStateRedisKey(roomID, userID string) string {
	return fmt.Sprintf("room:%s:player:%s:private", roomID, userID)
}

// Oyuncunun herkese açık durumu – diğer oyuncuların görebileceği bilgiler
func GeneratePlayerPublicStateRedisKey(roomID, userID string) string {
	return fmt.Sprintf("room:%s:player:%s:public", roomID, userID)
}

func GetGlobalLimiter() *limiter.Limiter {
	var client = GetRedisClient()
	rate, err := limiter.NewRateFromFormatted("100-M") // Dakikada 100 istek
	//rate, err := limiter.NewRateFromFormatted("10-S") //Saniyede 10 istek
	if err != nil {
		log.Fatal(err)
	}

	store, err := redisStore.NewStoreWithOptions(client.c, limiter.StoreOptions{
		Prefix:   "rate_limiter",
		MaxRetry: 3,
	})
	if err != nil {
		log.Fatal(err)
	}

	globalLimiter := limiter.New(store, rate)
	return globalLimiter
}

// Tüm verileri Redis'e kaydeder
func SaveGameToRedis(roomState *Model.RoomState, privateStates []Model.PlayerPrivateState, publicStates []Model.PlayerPublicState) error {
	client := GetRedisClient()

	// Oda genel durumu
	if err := client.SetKey(GenerateRoomStateRedisKey(roomState.RoomID), roomState, 30*time.Minute); err != nil {
		return err
	}

	// Public durumlar ayrı ayrı tutulur
	for _, pub := range publicStates {
		pubKey := GeneratePlayerPublicStateRedisKey(roomState.RoomID, pub.UserID)
		if err := client.SetKey(pubKey, pub, 30*time.Minute); err != nil {
			return err
		}
	}

	// Private durumlar
	for _, pvt := range privateStates {
		pvtKey := GeneratePlayerPrivateStateRedisKey(roomState.RoomID, pvt.UserID)
		if err := client.SetKey(pvtKey, pvt, 30*time.Minute); err != nil {
			return err
		}
	}

	return nil
}

// Oyuna reconnect eden bir kullanıcı için state’leri geri döner
func LoadGameForPlayer(roomID, userID string) (*Model.RoomState, *Model.PlayerPrivateState, []Model.PlayerPublicState, error) {
	client := GetRedisClient()

	var room Model.RoomState
	if err := client.GetKey(GenerateRoomStateRedisKey(roomID), &room); err != nil {
		return nil, nil, nil, err
	}

	// Oyuncunun kendi private state'i
	var private Model.PlayerPrivateState
	if err := client.GetKey(GeneratePlayerPrivateStateRedisKey(roomID, userID), &private); err != nil {
		return nil, nil, nil, err
	}

	// Diğer oyuncuların public state'leri
	var publicStates []Model.PlayerPublicState
	for _, basic := range room.Players { // Artık PlayerBasicInfo var
		var ps Model.PlayerPublicState
		if err := client.GetKey(GeneratePlayerPublicStateRedisKey(roomID, basic.UserID), &ps); err == nil {
			publicStates = append(publicStates, ps)
		}
	}

	return &room, &private, publicStates, nil
}
