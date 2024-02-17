package Services

import (
	"216/internal/orchestrator/Database"
	"216/internal/orchestrator/Entities"
	"216/internal/orchestrator/Services"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func fetchTask() (Entities.ArithmeticExpressions, int) {
	resp, err := http.Get("http://localhost:8181/api/v1/task")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	data := Entities.ArithmeticExpressions{}

	_ = json.Unmarshal(body, &data)
	return data, resp.StatusCode
}

func fetchParts(id uuid.UUID) (Entities.ArithmeticExpressions, int) {
	resp, err := http.Get("http://localhost:8181/api/v1/task/" + id.String())
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	data := Entities.ArithmeticExpressions{}
	//Database.Instance.Model(&Entities.ArithmeticExpressions{}).Where("id = ?", data.ID).Update("status", Entities.PROGRESS)
	_ = json.Unmarshal(body, &data)
	return data, resp.StatusCode
}

func PollAPI(interval time.Duration, numberComputing int, prefix string) {
	var Operand1, Operand2 float64

	for {
		data, code := fetchTask()
		//log.Println(code)
		//log.Println(1)
		if code == 200 {
			dataPart, codepart := fetchParts(data.ID)
			if codepart == 200 {
				//Database.Instance.Model(&Entities.ArithmeticExpressions{}).Where("id = ?", dataPart.ID).Update("status", Entities.PROGRESS)

				Database.Instance.Model(&Entities.ComputingResource{}).Where("name = ?", prefix+" "+strconv.Itoa(numberComputing)).Update("task", dataPart.ID)
				//log.Println(dataPart)
				//log.Println(code)
				//NullUuid, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")
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
				Database.Instance.Model(&Entities.ComputingResource{}).Where("name = ?", prefix+" "+strconv.Itoa(numberComputing)).Update("task_str", calckString)
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
				Database.Instance.Model(&Entities.ComputingResource{}).Where("name = ?", prefix+" "+strconv.Itoa(numberComputing)).Update("task_str", "")
				Database.Instance.Model(&Entities.ComputingResource{}).Where("name = ?", prefix+" "+strconv.Itoa(numberComputing)).Update("task", nil)

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
			}
		}

		time.Sleep(interval)
	}
}
