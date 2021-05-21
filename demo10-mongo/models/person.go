package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

// Location is a GeoJSON type.
type Location struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

type Point struct {
	// ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name"`
	Age     int             `json:"age"`
	City     string             `json:"city"`
	Location Location           `json:"location"`
}

type IPoint struct {
	Dis float64 `json:"dist"`
	Point Point
}

const (
	DBName = "sai"
	CollectionName = "persons"
	Key = "location"
)

func GetClient() *mgo {
	opt := options.Client().ApplyURI("mongodb://root:11111@localhost:27017")
	opt.SetLocalThreshold(3 * time.Second)     //只使用与mongo操作耗时小于3秒的
	opt.SetMaxConnIdleTime(5 * time.Second)    //指定连接可以保持空闲的最大毫秒数
	opt.SetMaxPoolSize(200)                    //使用最大的连接数
	opt.SetReadConcern(readconcern.Majority()) //指定查询应返回实例的最新数据确认为，已写入副本集中的大多数成员

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), opt)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return &mgo{client}
}

type mgo struct {
	client *mongo.Client
}

func(mgo *mgo) Start() {
	collection := mgo.client.Database(DBName).Collection(CollectionName)
	collection.Drop(context.TODO())

	// 设置存储格式 2dsphere
	collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.M{Key: "2dsphere"},
	})

	a := Point{"王二", 18, "杭州", Location{"Point", []float64{120.185614,30.300738}}}
	b := Point{"张三", 25, "杭州", Location{"Point", []float64{120.094778,30.310217}}}
	c := Point{"小晴", 35, "绍兴", Location{"Point", []float64{120.603847,30.054237}}}
	d := Point{"李四", 34, "杭州", Location{"Point", []float64{120.110893,30.207849}}}
	e := Point{"小明", 24, "北京", Location{"Point", []float64{116.435721,39.914031}}}
	f := Point{"吴六", 25, "杭州", Location{"Point", []float64{120.126443,30.33084}}}
	h := Point{"于一", 23, "杭州", Location{"Point", []float64{120.28132,30.184083}}}
	j := Point{"小七", 14, "杭州", Location{"Point", []float64{119.73926,30.247639}}}

	// 单条插入
	insertResult, err := collection.InsertOne(context.TODO(), a)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	ps := []interface{}{b, c, d, e, f, h, j}

	// 批量插入
	insertManyResult, err := collection.InsertMany(context.TODO(), ps)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
}

func (mgo *mgo) Near() {
	collection := mgo.client.Database(DBName).Collection(CollectionName)
	cur, err := collection.Find(context.TODO(), bson.D{
		{Key, bson.D{
				{"$near", bson.D{
					{
						"$geometry", Location{
							"Point",
							[]float64{120.110893,30.2078490},
						},
					},
					{"$maxDistance", 15000},
				}},
			}},
	})

	if err != nil {
		fmt.Println(err)
		return
	}
	var results []Point

	for cur.Next(context.TODO()) {
		var elem Point
		err := cur.Decode(&elem)
		fmt.Println(elem)
		fmt.Println(cur)
		if err != nil {
			fmt.Println("Could not decode Point")
			return
		}

		results = append(results, elem)
	}
	fmt.Println("查找到", len(results))
}

func(mgo *mgo) Close() {
	mgo.client.Disconnect(context.TODO())
	fmt.Println("!!!Disconnect MongoDB")
}
