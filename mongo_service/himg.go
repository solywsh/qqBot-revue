package mongo_service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

//// 接收数据用
//type Lolicon struct {
//	Error string        `json:"error"`
//	Data  []LoliconData `json:"data"`
//}

type LoliconData struct {
	Pid        int      `json:"pid"`
	P          int      `json:"p"`
	Uid        int      `json:"uid"`
	Title      string   `json:"title"`
	Author     string   `json:"author"`
	R18        bool     `json:"r18"`
	Width      int      `json:"width"`
	Height     int      `json:"height"`
	Tags       []string `json:"tags"`
	Ext        string   `json:"ext"`
	UploadDate int64    `json:"uploadDate"`
	Urls       struct {
		Original string `json:"original"`
	} `json:"urls"`
}

func (m *Mongo) GetLoLiCon(r18 bool) (*LoliconData, error) {
	//DbCollection = DbClient.Database("revue").Collection("himg")
	DbCollection := m.Client.Database("revue").Collection("himg")
	var result []LoliconData
	aggregate, err := DbCollection.Aggregate(context.TODO(), []bson.M{
		{"$match": bson.M{"r18": r18}}, {"$sample": bson.M{"size": 1}},
	})
	if err != nil {
		fmt.Println("aggregate error:", err)
		return nil, err
	}
	err = aggregate.All(context.TODO(), &result)
	if err != nil {
		fmt.Println("Decode error:", err)
		return nil, err
	}
	return &result[0], nil
}

//func main() {
//	if mon, ok := NewMongo(); ok {
//		fmt.Println(mon.GetLoLiCon(true))
//	}
//}
