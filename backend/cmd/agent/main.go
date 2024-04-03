package main

import (
	"216/internal/Redis"
	"216/internal/Types"
	"216/internal/orchestrator/Database"
	"216/internal/orchestrator/Entities"
	pb "216/proto"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gofrs/uuid"
	"github.com/joho/godotenv"
)

var mu sync.Mutex

//	func setValueInRedis(client *redis.Client, wg *sync.WaitGroup, key string, value interface{}) {
//		defer wg.Done()
//		mu.Lock()
//		defer mu.Unlock()
//
//		//jsonString, err := json.Marshal(value)
//
//		err := client.HSet(key, "", 0).Err()
//		if err != nil {
//			panic(err)
//		}
//		fmt.Printf("Value set for key %s successfully\n", key)
//	}
//
//	func GetFromRedis(client *redis.Client, wg *sync.WaitGroup, key string) []interface{} {
//		defer wg.Done()
//		mu.Lock()
//		defer mu.Unlock()
//
//		return client.Get(key)
//	}
func HeartBeat(ID uuid.UUID) {
	val, err := Redis.Client.Get("calculators").Result()
	if err != nil {
		log.Println(err)
	}
	var calculators []Types.Calculator
	if err := json.Unmarshal([]byte(val), &calculators); err != nil {
		log.Println(err)
	}

	for i := range calculators {
		if calculators[i].Id == ID {
			calculators[i].HeartBeat = int(time.Now().Unix())
			break
		}
	}

	jsonArray, err := json.Marshal(calculators)
	if err != nil {
		log.Println(err)
	}

	err = Redis.Client.Set("calculators", jsonArray, 0).Err()
	log.Println("Обновили", ID)
	if err != nil {
		log.Println(err)
	}
}

var calculators []Types.Calculator

func calculate(wg *sync.WaitGroup, uuidCalculator uuid.UUID) {
	//var m = []Types.Calculator{}
	defer wg.Done()
	mu.Lock()
	defer mu.Unlock()
	var compuers []Types.Calculator

	exists, err := Redis.Client.Exists("calculators").Result()
	if err != nil {
		log.Fatal(err)
	}
	if exists == 0 {

		calculators = append(compuers, Types.Calculator{Id: uuidCalculator, Name: "name", Status: Types.ALIVE, HeartBeat: int(time.Now().Unix())})

		jsonArray, err := json.Marshal(calculators)
		if err != nil {
			log.Fatal(err)
		}
		err = Redis.Client.Set("calculators", jsonArray, 0).Err()
		if err != nil {
			log.Fatal(err)
		}
	} else {

		val, err := Redis.Client.Get("calculators").Result()
		if err != nil {
			log.Fatal(err)
		}

		if err := json.Unmarshal([]byte(val), &calculators); err != nil {
			log.Fatal(err)
		}

		calculators = append(calculators, Types.Calculator{Id: uuidCalculator, Name: "name", Status: Types.ALIVE, HeartBeat: int(time.Now().Unix())})
		jsonArray, err := json.Marshal(calculators)
		if err != nil {
			log.Fatal(err)
		}
		err = Redis.Client.Set("calculators", jsonArray, 0).Err()
		if err != nil {
			log.Fatal(err)
		}
	}

}

func computingResource(i int, prefix string) {
	uuidCalculator, _ := uuid.NewV4()
	Database.Instance.Create(Entities.ComputingResource{
		Id:        uuidCalculator,
		Name:      prefix + " " + strconv.Itoa(i),
		Task:      nil,
		TaskStr:   nil,
		Status:    Types.ALIVE,
		HeartBeat: int(time.Now().Unix()),
	})
}

//	type SafeDBQuery struct {
//		db    *gorm.DB
//		mutex sync.Mutex
//	}
//
//	func NewSafeDBQuery(db *gorm.DB) *SafeDBQuery {
//		return &SafeDBQuery{db: db}
//	}
//
//	func (q *SafeDBQuery) GetTask() (Entities.ArithmeticExpressions, error) {
//		q.mutex.Lock()
//		defer q.mutex.Unlock()
//
//		var expression Entities.ArithmeticExpressions
//		result := q.db.Find(&expression)
//		return expression, result.Error
//	}

func loadF(i string) {

	host := "grpc-server"
	port := "5000"

	addr := fmt.Sprintf("%s:%s", host, port)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		os.Exit(1)
	}
	//for {
	grpcClient := pb.NewExpressionClient(conn)
	//var expression Entities.ArithmeticExpressions
	//result := Database.Instance.Preload("ExpressionPart").Preload("PreviousExpression").Preload("NextExpression").Where("parent IS NULL AND status = ?", Entities.WHAIT).First(&expression)
	//	if result.RowsAffected > 0 {
	expression, err := grpcClient.Do(context.TODO(), &pb.Request{Messgae: i})
	if err != nil {
		log.Println("failed invoked Expression: ", err)
	}
	fmt.Println("Expression: ", expression.Message)
	//	}
	//time.Sleep(interval)
	//}
}

func main() {
	err := godotenv.Load(".env")
	Prefix, err := rand.Int(rand.Reader, big.NewInt(80000))
	log.Println(Prefix)
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	Database.Connect(os.Getenv("DB_CONNECTION_STRING"))
	//Redis.InitRedis(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"))
	//Redis.Client.Del("calculators")
	Database.Instance.Model(&Entities.ArithmeticExpressions{}).Where("status = ?", Entities.PROGRESS).Update("status", Entities.WHAIT)

	//defer conn.Close()

	numberOfComputers, err := strconv.Atoi(os.Getenv("NUMBER_OF_COMPUTERS"))
	if err != nil {
		log.Println("Не верно задао кол-во вычислителей")
	}

	var wg sync.WaitGroup
	ch := make(chan string)
	//parts := make(chan int, numberOfComputers)
	//wg.Add(numberOfComputers)
	for i := 0; i < numberOfComputers; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			var PartsCount int
			for {
				Database.Instance.Raw("SELECT COUNT(*) FROM arithmetic_expressions WHERE parent IS NULL AND status = ?", Entities.WHAIT).Scan(&PartsCount)
				if PartsCount > 0 {
					loadF(Prefix.String() + " " + strconv.Itoa(i))
				}
				time.Sleep(time.Duration(i+1) * time.Second)
			}

			//fmt.Println(i)
			//Services.PollAPI(time.Duration(i)*time.Second, i, Prefix.String())

		}(i)
		//uuidCalculator, _ := uuid.NewV4()

		//go Services.PollAPI(5*time.Second, i, Prefix.String())
		go computingResource(i, Prefix.String())
		//go calculate(&wg, uuidCalculator)
		//wg.Done()
	}

	for {
		select {
		case result := <-ch:
			println(result)
		case <-time.After(5 * time.Second):
			for _, v := range calculators {
				HeartBeat(v.Id)
			}
			for i := 0; i < numberOfComputers; i++ {
				Database.Instance.Model(&Entities.ComputingResource{}).Where("name = ?", Prefix.String()+" "+strconv.Itoa(i)).Update("heart_beat", int(time.Now().Unix()))
			}
		}
	}
	wg.Wait()
	//wg.Wait()
}
