package main

import (
	"216/internal/orchestrator/Database"
	"216/internal/orchestrator/Entities"
	"216/internal/orchestrator/Services"
	pb "216/proto"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type Server struct {
	pb.ExpressionServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Do(
	ctx context.Context,
	in *pb.Request,
) (*pb.Response, error) {
	log.Println("invoked Expression: ", in)

	var expression Entities.ArithmeticExpressions
	result := Database.Instance.Preload("ExpressionPart").Preload("PreviousExpression").Preload("NextExpression").Where("parent IS NULL AND status = ?", Entities.WHAIT).First(&expression)

	var dataPart Entities.ArithmeticExpressions

	var Operand1, Operand2 float64

	for _, v := range expression.ExpressionPart {
		if v.Status == Entities.WHAIT {
			// fmt.Println(v)
			// fmt.Println(v)
			// fmt.Println(v)
			var dataPart1 Entities.ArithmeticExpressions
			Database.Instance.Preload("ExpressionPart").Preload("PreviousExpression").Preload("NextExpression").First(&dataPart1, v.ID)
			if dataPart1.Status != Entities.WHAIT {
				continue
			}
			fmt.Println(dataPart1)
			if dataPart1.Next == nil && dataPart1.Previous == nil {
				dataPart = dataPart1
				Database.Instance.Model(&Entities.ArithmeticExpressions{}).Where("id = ?", dataPart.ID).Update("status", Entities.PROGRESS)
				Database.Instance.Model(&Entities.ComputingResource{}).Where("name = ?", in.Messgae).Update("task", dataPart.ID)
				break
			}
			if dataPart1.Next == nil && dataPart1.Previous != nil {
				if dataPart1.PreviousExpression.Status == Entities.SUCCESS {
					dataPart = dataPart1
					Database.Instance.Model(&Entities.ArithmeticExpressions{}).Where("id = ?", dataPart.ID).Update("status", Entities.PROGRESS)
					Database.Instance.Model(&Entities.ComputingResource{}).Where("name = ?", in.Messgae).Update("task", dataPart.ID)
					break
				}
			}
			if dataPart1.Previous == nil && dataPart1.Next != nil {
				if dataPart1.NextExpression.Status == Entities.SUCCESS {
					dataPart = dataPart1
					Database.Instance.Model(&Entities.ArithmeticExpressions{}).Where("id = ?", dataPart.ID).Update("status", Entities.PROGRESS)
					Database.Instance.Model(&Entities.ComputingResource{}).Where("name = ?", in.Messgae).Update("task", dataPart.ID)
					break
				}
			}
			if dataPart1.Previous != nil && dataPart1.Next != nil {
				if dataPart1.NextExpression.Status == Entities.SUCCESS && dataPart1.PreviousExpression.Status == Entities.SUCCESS {
					dataPart = dataPart1
					Database.Instance.Model(&Entities.ArithmeticExpressions{}).Where("id = ?", dataPart.ID).Update("status", Entities.PROGRESS)
					Database.Instance.Model(&Entities.ComputingResource{}).Where("name = ?", in.Messgae).Update("task", dataPart.ID)
					break
				}
			}

			//continue
		}
	}

	//}
	log.Println(dataPart.NextExpression)
	calck := dataPart.StrValue
	calckString := ""
	if dataPart.Previous != nil {

		//Database.Instance.Raw("SELECT result FROM arithmetic_expressions WHERE id = ?", dataPart.Previous).Scan(&Operand1)
		if len(dataPart.Operation) > 0 {
			Operand1, _ = strconv.ParseFloat(dataPart.PreviousExpression.Result, 64)
		} else {
			Operand1 = dataPart.Operand1
		}

		//=
		log.Println("previous not null")
		//log.Println(getResult(dataPart.Previous))
	}
	if dataPart.Next != nil {
		//Database.Instance.Raw("SELECT result FROM arithmetic_expressions WHERE id = ?", dataPart.Next).Scan(&Operand2)
		if len(dataPart.Operation) > 0 {
			Operand2, _ = strconv.ParseFloat(dataPart.NextExpression.Result, 64)
		} else {
			Operand2 = dataPart.Operand1
		}
		log.Println("next not null")
	}

	if dataPart.Previous != nil && dataPart.Next == nil {
		log.Println("not nill")
		calck = fmt.Sprintf("%.2f", Operand1) + " " + fmt.Sprintf("%.2f", dataPart.Operand1) + " " + dataPart.Operation
		calckString = fmt.Sprintf("%.2f", Operand1) + " " + dataPart.Operation + " " + fmt.Sprintf("%.2f", dataPart.Operand1)
	} else if dataPart.Previous == nil && dataPart.Next == nil {
		log.Println("nill nill")
		calck = fmt.Sprintf("%.2f", dataPart.Operand1) + " " + fmt.Sprintf("%.2f", dataPart.Operand2) + " " + dataPart.Operation
		calckString = fmt.Sprintf("%.2f", dataPart.Operand1) + " " + dataPart.Operation + " " + fmt.Sprintf("%.2f", dataPart.Operand2)
	} else if dataPart.Previous != nil && dataPart.Next != nil {
		log.Println("not not")
		calck = fmt.Sprintf("%.2f", Operand1) + " " + fmt.Sprintf("%.2f", Operand2) + " " + dataPart.Operation
		calckString = fmt.Sprintf("%.2f", Operand1) + " " + dataPart.Operation + " " + fmt.Sprintf("%.2f", Operand2)
	} else if dataPart.Previous == nil && dataPart.Next != nil {
		log.Println("nill not")
		calck = fmt.Sprintf("%.2f", dataPart.Operand1) + " " + fmt.Sprintf("%.2f", Operand2) + " " + dataPart.Operation
		calckString = fmt.Sprintf("%.2f", dataPart.Operand1) + " " + dataPart.Operation + " " + fmt.Sprintf("%.2f", Operand2)
	}
	log.Println(calck)
	log.Println(calckString)
	Database.Instance.Model(&Entities.ComputingResource{}).Where("name = ?", in.Messgae).Update("task_str", calckString)
	//	dataPart.Status = Entities.PROGRESS
	//	dataPart.Result = fmt.Sprintf("%f", Services.EvaluatePostfix(dataPart.StrValue))
	//	log.Println(dataPart)
	//Database.Instance.Save(dataPart)
	//log.Println(calck)
	//time.Sleep(5)
	calckResult := fmt.Sprintf("%.2f", Services.EvaluatePostfix(calck))
	//log.Println(calckResult, "Proverka")
	if calckResult == "+Inf" || calckResult == "NaN" || calckResult == "-Inf" {
		calckResult = "Ошибка"
	}
	Database.Instance.Model(&Entities.ArithmeticExpressions{}).Where("id = ?", dataPart.ID).Updates(Entities.ArithmeticExpressions{ExpressionString: calckString, Status: Entities.SUCCESS, FinishedAt: int(time.Now().Unix()), Result: calckResult})
	Database.Instance.Model(&Entities.ComputingResource{}).Where("name = ?", in.Messgae).Update("task_str", "")
	Database.Instance.Model(&Entities.ComputingResource{}).Where("name = ?", in.Messgae).Update("task", nil)

	var PartsCount int
	Database.Instance.Raw("SELECT COUNT(*) FROM arithmetic_expressions WHERE parent = ?", dataPart.Parent).Scan(&PartsCount)
	var PartsCountSuccess int
	Database.Instance.Raw("SELECT COUNT(*) FROM arithmetic_expressions WHERE parent = ? AND status=?", dataPart.Parent, Entities.SUCCESS).Scan(&PartsCountSuccess)
	//log.Println(PartsCount, "fggfdfg")
	//log.Println(PartsCountSuccess)
	if PartsCount == PartsCountSuccess {
		Database.Instance.Model(&Entities.ArithmeticExpressions{}).Where("id = ?", dataPart.Parent).Updates(Entities.ArithmeticExpressions{Status: Entities.SUCCESS, FinishedAt: int(time.Now().Unix()), Result: calckResult})
		//log.Println(calckResult)
	}
	result1 := Database.Instance.Where("parent IS NULL AND status = ?", Entities.WHAIT).First(&expression)

	if result1.RowsAffected > 0 {
		s.Do(ctx,
			in)
	}

	if result.RowsAffected == 0 {

	}
	log.Println(calckResult)

	return &pb.Response{
		Message: calckResult,
	}, nil
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	Database.Connect(os.Getenv("DB_CONNECTION_STRING"))

	host := os.Getenv("GRPC_LISTEN_HOST") //"0.0.0.0"
	port := os.Getenv("GRPC_PORT")        //"5000"

	addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Println("error starting tcp listener: ", err)
		os.Exit(1)
	}

	log.Println("tcp listener started at port: ", port)
	grpcServer := grpc.NewServer()
	expressionServiceServer := NewServer()

	pb.RegisterExpressionServer(grpcServer, expressionServiceServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Println("error serving grpc: ", err)
		os.Exit(1)
	}
}
