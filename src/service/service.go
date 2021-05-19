package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	model2 "test3/src/model"
	pro "test3/src/proto"

	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Client = GetRedisClient()
var MongoClient = GetMongoClient()
var Collection = MongoClient.Database("test").Collection("User")

var tm = 24 * time.Hour

// GetRedisClient 和redis客户端建立连接
func GetRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println(pong, err)
	}
	return client
}

// GetMongoClient 与mongodb客户端建立连接
func GetMongoClient() *mongo.Client {
	opts := options.Client().ApplyURI("mongodb://localhost:27017")

	//连接数据库
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Fatal(err)
	}

	//判断服务是不是可用
	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected succeed")

	return client
}

// Log 用户注册/登陆
func Log(UserID string) (User model2.User) {
	var Result model2.User
	err := Collection.FindOne(context.TODO(), bson.M{"uid": UserID}).Decode(&Result)
	if err != nil {
		fmt.Println("User 不存在")
		fmt.Println("ERROR : ", err)
		user := model2.User{UID: UserID}
		_, err := Collection.InsertOne(context.TODO(), user)
		if err != nil {
			log.Fatal(err)
		}
		return user
	}

	return Result
}

// CreateCode 创建礼品码
func CreateCode(GiftCodeInfo model2.GiftCodeInfo) string {
	code := Code()
	GiftCodeInfo.GiftCode = code
	text, _ := json.Marshal(&GiftCodeInfo)
	Client.Set(code, text, tm)
	return code
}

// Code 生成8位礼品码
func Code() string {
	rand.Seed(time.Now().UnixNano())
	var letters = []byte("0123456789abcdefghijklmnopqrstuvwxyzABDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	GiftCode := string(b)
	return GiftCode
}

// Inquire 查询礼品码信息
func Inquire(GiftCode string) (GiftCodeInfo model2.GiftCodeInfo) {
	value, _ := Client.Get(GiftCode).Result()
	if value == "" {
		fmt.Println("礼品码不存在")
		return
	}

	errJson := json.Unmarshal([]byte(value), &GiftCodeInfo)
	if errJson != nil {
		fmt.Println("Unmarshal Error : ", errJson)
	}
	return GiftCodeInfo
}

// Verify 验证用户输入的礼品码
func Verify(UserID string, GiftCode string) (GiftInfo model2.GiftInfo) {
	var Result model2.User
	err := Collection.FindOne(context.TODO(), bson.M{"uid": UserID}).Decode(&Result)
	if err != nil {
		fmt.Println("用户不存在")
		return
	}

	value, _ := Client.Get(GiftCode).Result()
	if value == "" {
		fmt.Println("礼品码不存在")
		return
	}
	var GiftCodeInfo = model2.GiftCodeInfo{}
	errJson := json.Unmarshal([]byte(value), &GiftCodeInfo)
	if errJson != nil {
		fmt.Println("Unmarshal Error")
		return
	}
	if GiftCodeInfo.Times == 0 {
		fmt.Println("已被领取完")
		return
	}

	GiftCodeInfo.Times--
	GiftCodeInfo.TimesHasReceived++
	GiftCodeInfo.ListReceived.UsersReceived += UserID + ";"
	GiftCodeInfo.ListReceived.DateReceived += "----" + UserID + ":" + time.Now().String() + ";"

	text, _ := json.Marshal(&GiftCodeInfo)
	Client.Set(GiftCode, text, tm)

	Collection.UpdateOne(context.TODO(), bson.M{"uid": UserID}, bson.M{"$inc": bson.M{"goldcoin": GiftCodeInfo.GiftText.GoldCoin, "diamond": GiftCodeInfo.GiftText.Diamond}})

	Convert(GiftCodeInfo.GiftText, UserID)
	return GiftCodeInfo.GiftText
}

func Convert(GiftInfo model2.GiftInfo, UserID string) {
	var Reward = pro.GeneralReward{}
	change := make(map[uint32]uint64)
	change[model2.GoldCoin] = uint64(GiftInfo.GoldCoin)
	change[model2.Diamond] = uint64(GiftInfo.Diamond)
	Reward.Changes = change

	balance := make(map[uint32]uint64)
	var Result model2.User
	err := Collection.FindOne(context.TODO(), bson.M{"uid": UserID}).Decode(&Result)
	if err != nil {
		return
	}
	balance[model2.GoldCoin] = uint64(Result.GoldCoin)
	balance[model2.Diamond] = uint64(Result.Diamond)
	Reward.Balance = balance

	out, err := proto.Marshal(&Reward)
	fmt.Println("------------out:", out)
	if err != nil {
		fmt.Println("Error : ", err)
	}

	User := pro.GeneralReward{}
	proto.Unmarshal(out, &User)
	fmt.Println("decode : ", User)

}
