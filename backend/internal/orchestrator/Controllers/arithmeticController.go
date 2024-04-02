package Controllers

import (
	"216/internal/orchestrator/Database"
	"216/internal/orchestrator/Entities"
	"216/internal/orchestrator/Services"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

func checkIfRequestExists(requestId string) bool {
	var expression Entities.ArithmeticExpressions
	result := Database.Instance.Where("request_id like ?", requestId).First(&expression)
	if result.RowsAffected == 0 {
		return false
	}
	return true
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
func AddArithmeticExpressions(w http.ResponseWriter, r *http.Request) {
	var expression Entities.ArithmeticExpressions
	requestId := r.Header.Get("X-Request-ID")

	w.Header().Set("Content-Type", "application/json")
	if checkIfRequestExists(requestId) {
		Database.Instance.Where("request_id=?", requestId).First(&expression)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expression)
		return
	}

	json.NewDecoder(r.Body).Decode(&expression)
	if len(expression.ExpressionString) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(requestId) == 0 {
		requestId = GetMD5Hash(expression.ExpressionString)
	}
	expression.RequestId = &requestId

	if !Services.IsValidExpression(expression.ExpressionString) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	postfix := Services.ToPostfix(expression.ExpressionString) //Services.ToPostfix(expression.ExpressionString)

	log.Println("Postfix notation:", postfix)

	uidM, _ := uuid.NewV4()
	expression.ID = uidM
	expression.StrValue = postfix
	ow, _ := Services.ExtractTokenID(r)
	//var user Entities.User
	//Database.Instance.First(&user, ow)
	expression.Owner = ow

	part := []Entities.ArithmeticExpressions{}
	order := 0
	for _, v := range Services.Mathslice(postfix) {

		log.Println("'", v.Operation, "'")

		partExpression := Entities.ArithmeticExpressions{
			ID:               v.Id,
			ExpressionString: fmt.Sprintf("%.2f", v.Operand1) + v.Operation + fmt.Sprintf("%.2f", v.Operand2),
			Status:           Entities.WHAIT,
			Operation:        v.Operation,
			Operand1:         v.Operand1,
			Operand2:         v.Operand2,
			StrValue:         v.Polsk,
			OrderField:       order,
		}
		if v.Previous != nil {
			partExpression.Previous = &v.Previous.Id
		}

		if v.Next != nil {
			partExpression.Next = &v.Next.Id
		}
		partExpression.Owner = ow
		part = append(part, partExpression)
		order++
	}

	log.Println(part)
	expression.ExpressionPart = part

	Database.Instance.Create(&expression)

	json.NewEncoder(w).Encode(expression)
}

func ArithmeticExpressionsList(w http.ResponseWriter, r *http.Request) {
	var expressions []Entities.ArithmeticExpressions

	log.Println(Services.ExtractTokenID(r))
	//
	uid, _ := Services.ExtractTokenID(r)
	Database.Instance.Preload("PreviousExpression").Preload("NextExpression").Preload("ExpressionPart").Where("parent is null").Where("owner=?", uid).Order("added_at desc").Find(&expressions)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expressions)
}

func checkIfExpressionExists(expressionId uuid.UUID) bool {
	var expression Entities.ArithmeticExpressions
	result := Database.Instance.First(&expression, expressionId)
	if result.RowsAffected == 0 {
		return false
	}
	return true
}

func ArithmeticExpression(w http.ResponseWriter, r *http.Request) {
	expressionId := uuid.FromStringOrNil(mux.Vars(r)["id"])
	if checkIfExpressionExists(expressionId) == false {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Expression Not Found!")
		return
	}
	var expression Entities.ArithmeticExpressions
	Database.Instance.First(&expression, expressionId)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expression.Result)
}

func AvailableArithmeticOperations(w http.ResponseWriter, r *http.Request) {

	var operations []Entities.ArithmeticOperation
	Database.Instance.Find(&operations)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(operations)
}

func TaskForExecution(w http.ResponseWriter, r *http.Request) {
	var expression Entities.ArithmeticExpressions
	result := Database.Instance.Where("parent IS NULL AND status = ?", Entities.WHAIT).First(&expression)

	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expression)
}

func TaskForExecutionPart(w http.ResponseWriter, r *http.Request) {
	expressionId := uuid.FromStringOrNil(mux.Vars(r)["id"])
	log.Println(expressionId)
	var expression Entities.ArithmeticExpressions
	result := Database.Instance.Preload("NextExpression").Preload("PreviousExpression").Where("parent=? AND status = ? AND previous is null AND next is null", expressionId, Entities.WHAIT).First(&expression)
	log.Println("tut proverka", expression)
	if result.RowsAffected == 0 {
		var r []string
		Database.Instance.Raw("SELECT id From arithmetic_expressions WHERE parent = ? AND status=?", expressionId, Entities.SUCCESS).Scan(&r)

		log.Println(r, "результаты")
		result1 := Database.Instance.Preload("NextExpression").Preload("PreviousExpression").Where("(parent=? AND status = ?) and ((next is null and previous IN ?) OR (previous is null and  next IN ?) OR (previous IN ? and  next IN ?))",
			expressionId,
			Entities.WHAIT,
			r,
			r,
			r,
			r,
		).Order("order_field asc").Limit(1).Find(&expression)
		log.Println("po resusltatu", expression)
		if result1.RowsAffected == 0 {
			w.WriteHeader(http.StatusNotFound)
		}
		//w.WriteHeader(http.StatusNotFound)
	}
	Database.Instance.Model(&Entities.ArithmeticExpressions{}).Where("id = ?", expression.ID).Update("status", Entities.PROGRESS)
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expression)
}

func SetResult(w http.ResponseWriter, r *http.Request) {
	expressionId := uuid.FromStringOrNil(mux.Vars(r)["id"])
	if checkIfExpressionExists(expressionId) == false {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Expression Not Found!")
		return
	}

	var expression Entities.ArithmeticExpressions
	Database.Instance.First(&expression, expressionId)
	json.NewDecoder(r.Body).Decode(&expression)

	expression.SetResult(expression.Result)

	Database.Instance.Save(&expression)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expression)
}

func checkIfArithmeticOperationExists(operationId uuid.UUID) bool {
	var operation Entities.ArithmeticOperation
	result := Database.Instance.First(&operation, operationId)
	if result.RowsAffected == 0 {
		return false
	}
	return true
}

func SetLeadTimeToArithmeticOperation(w http.ResponseWriter, r *http.Request) {
	operationId := uuid.FromStringOrNil(mux.Vars(r)["id"])
	//log.Println(operationId)
	//log.Println(mux.Vars(r)["id"])
	if checkIfArithmeticOperationExists(operationId) == false {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Expression Not Found!")
		return
	}

	var operation Entities.ArithmeticOperation
	Database.Instance.First(&operation, operationId)
	json.NewDecoder(r.Body).Decode(&operation)

	Database.Instance.Save(&operation)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(operation)
}
