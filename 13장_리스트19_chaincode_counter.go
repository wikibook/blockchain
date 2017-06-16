package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type CounterChaincode struct {
}

// 카운터 정보
type Counter struct {
	Name string `json:"name"`
	Counts uint64 `json:"counts"`
}

const numOfCounters int = 3

// 카운터 정보의 초기 값을 설정
func (cc *CounterChaincode) Init(stub *shim.ChaincodeStub, function string, args []string)([]byte, error) {
	var counters [numOfCounters]Counter
	var countersBytes [numOfCounters][]byte

	// 카운터 정보 생성
	counters[0] = Counter{Name: "Office Worker", Counts: 0}
	counters[1] = Counter{Name: "Home Worker", Counts: 0}
	counters[2] = Counter{Name: "Student", Counts: 0}

	// 카운터 정보를 월드 스테이트에 추가
	for i := 0; i < len(counters); i++ {
		// JSON 형식으로 변경
		countersBytes[i], _ = json.Marshal(counters[i])
		// 월드 스테이트에 추가
		stub.PutState(strconv.Itoa(i), countersBytes[i])
	}
	return nil, nil
}
// 카운터 정보 갱신

func (cc *CounterChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	// function 이름으로 핸들링
	if function == "countUp" {
		// 카운트 증가 수행
		return cc.countUp(stub, args)
	}
	return nil, errors.New("Received unknown function")
}

//　카운터 정보 참조
func (cc *CounterChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	// function 이름으로 핸들링
	if function == "refresh" {
		// 카운터 정보 취득
		return cc.getCounters(stub, args)
	}
	return nil, errors.New("Received unknown function")
}
// 카운트 증가 수행
func (cc *CounterChaincode) countUp(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	// 월드 스테이트로부터 선택된 카운터 정보 취득
	counterId := args[0]
	counterJson, _ := stub.GetState(counterId)

	// 취득한 JSON 형식 정보를 Counter로 변환
	counter := Counter{}
	json.Unmarshal(counterJson, &counter)

	// 카운트 증가
	counter.Counts++

	// 월드 스테이트에 변경 후의 값을 추가
	counterJson, _ = json.Marshal(counter)
	stub.PutState(counterId, counterJson)

	return nil, nil
}
// 카운터 정보 취득
func (cc *CounterChaincode) getCounters(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var counters [numOfCounters]Counter
	var countersBytes [numOfCounters][]byte
	for i := 0; i < len(counters); i++ {
		// 카운터 정보를 월드 스테이트로부터 취득
		countersBytes[i], _ = stub.GetState(strconv.Itoa(i))

		// 취득한 JSON 형식 정보를 Counter로 변환
		counters[i] = Counter{}
		json.Unmarshal(countersBytes[i], &counters[i])
	}

	//json 형식으로 변환
	return json.Marshal(counters)
}

// Validating Peer에 연결해 체인 코드를 실행
func main() {
	err := shim.Start(new(CounterChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}