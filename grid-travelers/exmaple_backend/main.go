package main

import (
	//"encoding/json"
	"fmt"
	"sync"

	//"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Grid struct {
	Cols int
	Rows int
}

type GridParams struct {
	Cols int `json:"cols"`
	Rows int `json:"rows"`
}

const TravelersAmount int = 10
const gridCols int = 10
const gridRows int = 10

var SGrid = Grid{Cols: gridCols, Rows: gridRows}

type Move int

const (
	Forward Move = iota
	Backward
	Left
	Right
)

func (m Move) ToVec() (int, int) {
	switch m {
	case Forward:
		return 0, 1
	case Backward:
		return 0, -1
	case Left:
		return -1, 0
	case Right:
		return 1, 0
	default:
		return 0, 0
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

//var redisClient = redis.NewClient(&redis.Options{
//	Addr: "localhost:6379",
//})

type Traveler struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type MyRandSrc struct {
	*rand.Rand
}

func (randomSrc MyRandSrc) getRandomColor() (int, int, int) {
	r := randomSrc.Intn(256)
	g := randomSrc.Intn(256)
	b := randomSrc.Intn(256)
	return r, g, b
}

func (randomSrc MyRandSrc) getRandomMove() (int, int) {
	return Move(rand.Intn(4)).ToVec()
}

type synchronizedSocketWriter struct {
	*websocket.Conn
	mutex sync.Mutex
}

func (s *synchronizedSocketWriter) WriteJSON(v interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.Conn.WriteJSON(v)
}

func main() {
	//http.Handle("/", http.FileServer(http.Dir("./public")))
	go http.HandleFunc("/ws", wsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	syncConn := &synchronizedSocketWriter{conn, sync.Mutex{}}
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	fmt.Println("SIMULATOR INITIALIZATION...")

	//Create a context for cancellation
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel() // Ensure cancellation when the handler exits
	//var wg sync.WaitGroup

	println("SENDING PARAMS...")

	params := GridParams{Rows: SGrid.Cols, Cols: SGrid.Rows}
	err = conn.WriteJSON(params)
	time.Sleep(time.Duration(5000) * time.Millisecond)
	if err != nil {
		log.Println("Failed sending the parameters:", err)
		return
	}

	//simulate travelers
	randSrc := MyRandSrc{rand.New(rand.NewSource(time.Now().UnixNano()))}
	travelers := populateGridStrategy(TravelersAmount)
	syncConn.WriteJSON(travelers)
	time.Sleep(time.Duration(3000) * time.Millisecond)

	fmt.Println("SIMULATOR RUNS...")

	for id, traveler := range travelers {
		go func() {
			for {
				time.Sleep(time.Duration(1000+randSrc.Intn(1500)) * time.Millisecond)

				traveler.takeRandomStep(randSrc)

				go syncConn.WriteJSON(map[string]Traveler{id: traveler})
				//fmt.Println("update: ", id, traveler)
			}
		}()
	}

	select {}
}

func populateGridStrategy(travelersAmount int) map[string]Traveler {
	travelersMapping := make(map[string]Traveler)
	for i := 0; i < min(SGrid.Rows, SGrid.Rows); i++ {
		travelersMapping[fmt.Sprintf("t%d", i)] = Traveler{
			X: i,
			Y: i,
			//Color: "red",
		}
	}
	return travelersMapping
}

func (p *Traveler) takeRandomStep(src MyRandSrc) {
	stepX, stepY := src.getRandomMove()
	p.moveBy(stepX, stepY)
}

func (p *Traveler) moveBy(stepX int, stepY int) {
	p.X += stepX
	if p.X < 0 {
		p.X = SGrid.Cols - 1
	}
	if p.X >= SGrid.Cols {
		p.X = 0
	}
	p.Y += stepY
	if p.Y < 0 {
		p.Y = SGrid.Rows - 1
	}
	if p.Y >= SGrid.Rows {
		p.Y = 0
	}
}
