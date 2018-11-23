package cg

import (
	"encoding/json"
	"errors"
	"ipc"
	"reflect"
	"sync"
)

var _ ipc.Server = &CenterServer{}

type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

type CenterServer struct {
	servers map[string]ipc.Server
	players []*Player
	//rooms   []*Room
	mutex   sync.RWMutex
	methods map[string]*reflect.Value
}

func NewCenterServer() *CenterServer {
	result := &CenterServer{}
	result.servers = make(map[string]ipc.Server, 0)
	result.players = make([]*Player, 0)

	//扫描方法
	typ := reflect.TypeOf(result)
	val := reflect.ValueOf(result)
	numMethod := typ.NumMethod()
	for i := 0; i < numMethod; i++ {
		methodName := typ.Method(i).Name
		method := val.MethodByName(methodName)
		result.methods[methodName] = &method
	}

	return result
}

//新增一个玩家
func (server *CenterServer) AddPlayer(params string) error {
	player := NewPlayer()

	err := json.Unmarshal([]byte(params), &player)
	if err != nil {
		return err
	}

	server.mutex.Lock()
	defer server.mutex.Unlock()

	server.players = append(server.players, player)

	return nil
}

func (server *CenterServer) RemovePlayer(params string) error {

	server.mutex.Lock()
	defer server.mutex.Unlock()

	for i, v := range server.players {
		if v.Name == params {
			if len(server.players) == 1 {
				server.players = make([]*Player, 0)
			} else if i == len(server.players)-1 {
				server.players = server.players[:i]
			} else if i == 0 {
				server.players = server.players[1:]
			} else {
				server.players = append(server.players[:i], server.players[i+1:]...)
			}
		}
		return nil
	}

	return errors.New("PlayerNotFound")
}

func (server *CenterServer) ListPlayer(params string) (players string, err error) {
	server.mutex.RLock()
	defer server.mutex.RUnlock()

	if len(server.players) > 0 {
		b, _ := json.Marshal(server.players)
		players = string(b)
	} else {
		err = errors.New("NoPlayerOnline")
	}
	return
}

func (server *CenterServer) Broadcast(params string) error {
	var message Message
	err := json.Unmarshal([]byte(params), &message)
	if err != nil {
		return err
	}

	server.mutex.Lock()
	defer server.mutex.Unlock()

	if len(server.players) > 0 {
		for _, player := range server.players {
			player.mq <- &message
		}
	} else {
		err = errors.New("NoPlayerOnline")
	}
	return err
}

func (server *CenterServer) Handle(method, params string) *ipc.Response {
	if method == "Handle" {
		return nil
	}
	//寻找方法
	methodFunc, exist := server.methods[method]
	if !exist {
		return &ipc.Response{Code: errors.New("MethodNotExist").Error()}
	}
	//call 方法
	callResult := methodFunc.Call([]reflect.Value{reflect.ValueOf(params)})
	err := callResult[0].Interface().(error)
	if err != nil {
		return &ipc.Response{Code: err.Error()}
	}
	return &ipc.Response{Code: "200"}
}

func (server *CenterServer) Name() string {
	return "CenterServer"
}
