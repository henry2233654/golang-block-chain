package di

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang-block-chain/web/cores"
	"log"
	"reflect"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func detectFunc(fn interface{}) (inType, outType reflect.Type) {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		panic("fn is not a function")
	}

	if fnType.NumIn() != 3 || fnType.NumOut() != 0 {
		panic("function is not func(context.context, <-chan T1, chan<- T2)")
	}
	if fnType.In(0) != reflect.TypeOf((*context.Context)(nil)).Elem() {
		panic("function is not func(context.context, <-chan T1, chan<- T2)")
	}
	if fnType.In(1).Kind() != reflect.Chan {
		panic("function is not func(context.context, <-chan T1, chan<- T2)")
	}
	if fnType.In(2).Kind() != reflect.Chan {
		panic("function is not func(context.context, <-chan T1, chan<- T2)")
	}
	inType = fnType.In(1).Elem()
	outType = fnType.In(2).Elem()
	return
}

func AutoWsDi(fn interface{}) gin.HandlerFunc {
	inType, outType := detectFunc(fn)
	return func(c *gin.Context) {
		conn, err := cores.NewWebsocketConnection(c)
		defer conn.Close()
		//errCh := make(chan error)
		if err != nil {
			return
		}
		inChValue := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, inType), 0)
		outChValue := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, outType), 0)
		ctx, cancel := context.WithCancel(c.Request.Context())
		defer cancel()
		done := make(chan struct{})
		go outPump(ctx, conn, outChValue, done)
		go inPump(ctx, conn, inChValue, done)
		go func() {
			_ = reflect.ValueOf(fn).Call([]reflect.Value{reflect.ValueOf(ctx), inChValue, outChValue})
			done <- struct{}{}
		}()
		<-done
	}
}

func inPump(c context.Context, conn *websocket.Conn, inChValue reflect.Value, done chan struct{}) {
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	inType := inChValue.Type().Elem()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			//break
			log.Println(err)
			done <- struct{}{}
			return
		}
		dataPtr := reflect.New(inType).Interface()
		log.Println(string(message))
		if err := json.Unmarshal(message, dataPtr); err != nil {
			log.Println(err)
			done <- struct{}{}
			return
		}
		data := reflect.ValueOf(dataPtr).Elem().Interface()
		inChValue.Send(reflect.ValueOf(data)) //inCh <- data
	}
}

func outPump(c context.Context, conn *websocket.Conn, outChValue reflect.Value, done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	cases := []reflect.SelectCase{
		{Dir: reflect.SelectRecv, Chan: outChValue},
		{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ticker.C)},
		{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(c.Done())},
	}
	for {
		switch index, value, _ := reflect.Select(cases); index {
		case 0: //out := <-outCh
			out := value.Interface()
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			data, err := json.Marshal(out)
			if err != nil {
				log.Println(err)
				done <- struct{}{}
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
				done <- struct{}{}
				return
			}
		case 1: //<-ticker.C
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
				done <- struct{}{}
				return
			}
		case 2: //<-c.Done():
			return
		}

		//select {
		//case out := <-outCh:
		//	conn.SetWriteDeadline(time.Now().Add(writeWait))
		//	data, _ := json.Marshal(out)
		//	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		//		log.Println(err)
		//		done <- struct{}{}
		//		return
		//	}
		//case <-ticker.C:
		//	conn.SetWriteDeadline(time.Now().Add(writeWait))
		//	if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
		//		log.Println(err)
		//		done <- struct{}{}
		//		return
		//	}
		//case <-c.Done():
		//	return
		//}
	}
}
