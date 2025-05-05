package models

import (
	"context"
	"fmt"
	"grid-travelers-v4/config"
	"strconv"
	"time"
)

//type FieldSemaphore struct {
//	semaphore chan struct{}
//}
//
//type GridFieldServer struct {
//	fields [][]FieldSemaphore
//}

const (
	ResponseCodeSuccess = iota
	ResponseCodeError   // request has failed
	ResponseCodeTimeout // request has timed out
	ResponseCodeNoPath  // wild traveler has no space to move over
)

const (
	Empty = iota
	Wildlife
	GridTraveler
	Ambush
	Dead
)

//type TileContent struct {
//	contentType int
//	contentRef  *Traveler
//}

type ServerState struct {
	content int
	id      int
}

type ResponseType struct {
	ResponseCode    int
	ChangeTimestamp time.Time
}

type Requester struct {
	pos  Position
	id   int
	kind int
}

type RequestType struct {
	requester   Requester
	destination Position
	timestamp   time.Time
	timeout     time.Duration
	response    chan ResponseType
}

type ForcedMoveRequestType struct {
	destination Position
	timestamp   time.Time
	response    chan ResponseType
}

type RequestChannelType struct {
	channel chan RequestType
}

type TileRequestChannels struct {
	tiles [config.GridHeight][config.GridWidth]RequestChannelType
}

type ForcedMoveRequestChannelType struct {
	channel chan ForcedMoveRequestType
}

type ForcedMoveChannels struct {
	tiles [config.GridHeight][config.GridWidth]ForcedMoveRequestChannelType
}

func (c *TileRequestChannels) InitTileRequestChannels() {
	for x := 0; x < config.GridWidth; x++ {
		for y := 0; y < config.GridHeight; y++ {
			c.tiles[y][x].channel = make(chan RequestType)
		}
	}
}

func (c *ForcedMoveChannels) InitForcedMoveChannels() {
	for x := 0; x < config.GridWidth; x++ {
		for y := 0; y < config.GridHeight; y++ {
			c.tiles[y][x].channel = make(chan ForcedMoveRequestType)
		}
	}
}

func InitTileServers(tileRequestChannels *TileRequestChannels, forcedMoveChannels *ForcedMoveChannels, ctx context.Context) {
	for x := 0; x < config.GridWidth; x++ {
		for y := 0; y < config.GridHeight; y++ {
			go func() {
				err := TileServer(x, y, tileRequestChannels, forcedMoveChannels, ctx)
				if err != nil {
					fmt.Println("TileServer: ", err)
				}
			}()
		}
	}
}

func handleLeaveRequest(x int, y int, state *ServerState, request *RequestType, requests *TileRequestChannels) error {
	tileContent := state.content

	// tileContent needs to be moved onto a new tile

	if (request.destination == Position{x: x, y: y}) {
		request.response <- ResponseType{
			ResponseCode:    ResponseCodeError,
			ChangeTimestamp: time.Now(),
		}
		return fmt.Errorf("error: cannot make moves in place")
	}

	if tileContent == Dead {
		request.response <- ResponseType{
			ResponseCode:    ResponseCodeError,
			ChangeTimestamp: time.Now(),
		}
		return fmt.Errorf("error: cannot move dead content")
	}

	neighborResponseChannel := make(chan ResponseType)
	neighborRequestChannel := requests.tiles[request.destination.y][request.destination.x].channel

	timeoutDuration := max(time.Duration(0), request.timeout-time.Since(request.timestamp)) + config.MinTimeout
	timeout := time.After(timeoutDuration)
	log(x, y, *state, "asking for a permit from server(", request.destination.x, ",", request.destination.y, ") with timeout:", timeoutDuration)
	select {
	case neighborRequestChannel <- RequestType{
		requester:   request.requester,
		destination: request.destination,
		timestamp:   request.timestamp,
		timeout:     request.timeout,
		response:    neighborResponseChannel,
	}:
		neighborResponse := <-neighborResponseChannel
		switch neighborResponse.ResponseCode {
		case ResponseCodeSuccess:
			state.content = Empty
			state.id = -1
			log(x, y, *state, "traveler has left the tile", x, y, "neighbor response:", neighborResponse, "timestamp:", neighborResponse.ChangeTimestamp.Nanosecond())
			request.response <- neighborResponse
		case ResponseCodeError:
			log(x, y, *state, "neighbor response: Error")
			request.response <- neighborResponse
		case ResponseCodeTimeout:
			log(x, y, *state, "neighbor response: Timeout")
			request.response <- neighborResponse
		default:
			log(x, y, *state, "neighbor response:", neighborResponse)
			request.response <- neighborResponse
		}

	case <-timeout:
		log(x, y, *state, "neighbor request handling timed out", x, y)
		request.response <- ResponseType{
			ResponseCode:    ResponseCodeTimeout,
			ChangeTimestamp: time.Now(),
		}
	}
	return nil
}

// defines how the tile server behaves if it is empty
func handleEnteringEmptyTile(x int, y int, state *ServerState, request *RequestType, requests *TileRequestChannels) error {
	//originalPos := t.GetPosition()

	//request.requester.pos.x = x
	//request.requester.pos.y = y

	//_ = request.requester.LogState() // discard overflow error
	state.content = request.requester.kind
	state.id = request.requester.id
	time.Sleep(time.Microsecond) // for animation purposes
	timestamp := time.Now()
	log(x, y, *state, "handleEnteringEmptyTile: traveler", request.requester.id, "has been put on tile", x, y, "timestamp: ", time.Now().Nanosecond())

	request.response <- ResponseType{
		ResponseCode:    ResponseCodeSuccess,
		ChangeTimestamp: timestamp,
	}

	return nil
}

// defines how the tile server behaves if it has traveler on it
func handleEnteringOnTraveler(x int, y int, state *ServerState, request *RequestType, requests *TileRequestChannels) error {
	// requester is trying to move to a tile occupied by traveler
	request.response <- ResponseType{
		ResponseCode:    ResponseCodeError,
		ChangeTimestamp: time.Now(),
	}
	return nil
}

func log(x int, y int, state ServerState, args ...interface{}) {
	var contentString string
	switch state.content {
	case Empty:
		contentString = "EMPTY"
	case Wildlife:
		contentString = "WILDLIFE"
	case GridTraveler:
		contentString = "TRAVELER"
	case Dead:
		contentString = "DEAD"
	default:
		contentString = "UNKNOWN"
	}

	msg := fmt.Sprint(strconv.Itoa(time.Now().Nanosecond()), " server(", x, y, "), content:", contentString, ", id:", state.id, ", msg: ", args)
	//fmt.Fprintf(os.Stderr, "%s\n", msg)
	fmt.Printf("%s\n", msg)
}

func handleEnteringOnWildlife(x int, y int, state *ServerState, request *RequestType, requests *TileRequestChannels, forcedMoveChannels *ForcedMoveChannels) error {
	//request.response <- ResponseType{
	//	ResponseCode: ResponseCodeError,
	//}
	//return nil

	timeout := time.After(max(config.MinTimeout, request.timeout-(time.Since(request.timestamp))))
	var wildlifeDestination []Position = getNeighbours(Position{x: x, y: y}, request.requester.pos.x, request.requester.pos.y)
	for _, dest := range wildlifeDestination {
		queryTimeout := time.After(config.MinTimeout) // don't wait too long

		neighborResponseChannel := make(chan ResponseType)
		neighborRequestChannel := requests.tiles[dest.y][dest.x].channel
		select {
		case neighborRequestChannel <- RequestType{
			requester: Requester{
				pos:  Position{x, y},
				id:   state.id,
				kind: state.content,
			},
			destination: Position{x: dest.x, y: dest.y},
			timestamp:   request.timestamp,
			timeout:     request.timeout,
			response:    neighborResponseChannel,
		}:
			neighborResponse := <-neighborResponseChannel
			switch neighborResponse.ResponseCode {
			case ResponseCodeSuccess:
				//accept requester
				state.content = request.requester.kind
				state.id = request.requester.id
				request.requester.pos = Position{x: x, y: y}
				//_ = request.requester.LogState()

				forcedResponseChannel := make(chan ResponseType)
				forcedMoveChannels.tiles[y][x].channel <- ForcedMoveRequestType{
					destination: Position{x: dest.x, y: dest.y},
					timestamp:   neighborResponse.ChangeTimestamp,
					response:    forcedResponseChannel,
				}
				forcedResponse := <-forcedResponseChannel
				if forcedResponse.ResponseCode != ResponseCodeSuccess {
					log(x, y, *state, "error: forced move failed")
					request.response <- ResponseType{
						ResponseCode:    ResponseCodeError,
						ChangeTimestamp: time.Now(),
					}
					return fmt.Errorf("error: forced move failed")
				}

				time.Sleep(time.Microsecond) // for animation purposes
				request.response <- ResponseType{
					ResponseCode:    ResponseCodeSuccess,
					ChangeTimestamp: time.Now(),
				}
				return nil
			default:
				continue // try next destination
			}
		case <-queryTimeout:
			continue
		case <-timeout:
			request.response <- ResponseType{
				ResponseCode:    ResponseCodeTimeout,
				ChangeTimestamp: time.Now(),
			}
			return nil
		}
	}
	request.response <- ResponseType{
		ResponseCode:    ResponseCodeTimeout,
		ChangeTimestamp: time.Now(),
	}
	return nil
}

func TileServer(x int, y int, requestChannels *TileRequestChannels, forcedMoveChannels *ForcedMoveChannels, ctx context.Context) error {
	state := ServerState{
		content: Empty,
		id:      -1,
	}

	for {
		log(x, y, state, "waiting for request")

		select {
		case request := <-requestChannels.tiles[y][x].channel:

			log(x, y, state, "handling request:", request.requester.id, request.destination, request.timeout)

			if (request.destination.y == config.GridHeight) && (request.destination.x == config.GridWidth) {
				if state.id != request.requester.id {
					return fmt.Errorf("error: traveler %d is trying to leave the board, but he is not on this tile!", request.requester.id)
				}
				log(x, y, state, "requester", request.requester.id, "has left the board")
				state.content = Empty
				state.id = -1
				continue
			}

			//timestamp := time.Now()
			if state.content == Empty {
				log(x, y, state, "trying to put traveler", request.requester.id, request.requester.pos, " on tile", request.destination)
				err := handleEnteringEmptyTile(x, y, &state, &request, requestChannels)
				if err != nil || state.content == Empty {
					log(x, y, state, "error: content:", state.content,
						"\nerror: traveler", request.requester.id, " has not been put on tile", x, y,
						"\nerror msg:", err,
					)
				}
				continue
			}

			// assert content != nil

			//if request.requester.pos == state.content.pos &&
			//	request.requester.id != content.id {
			//	request.response <- ResponseType{
			//		ResponseCode: ResponseCodeError,
			//	}
			//	fmt.Println("error: two travellers were found on the same tile")
			//	continue
			//}

			if request.requester.id == state.id {
				_ = handleLeaveRequest(x, y, &state, &request, requestChannels)
				continue
			}

			if (request.destination != Position{x: x, y: y} && request.requester.kind != Wildlife) {
				request.response <- ResponseType{
					ResponseCode:    ResponseCodeError,
					ChangeTimestamp: time.Now(),
				}
				log(x, y, state, "requester asked for action wrong server")
				return fmt.Errorf("error: server cannot handle request that does not concern it, "+
					"requesterId: %d, position: (%d, %d), destination: (%d, %d), tile: (%d, %d)",
					request.requester.id,
					request.requester.pos.x, request.requester.pos.y,
					request.destination.x, request.destination.y,
					x, y,
				)
			}

			// now we are sure it's a request to enter this tile

			//if content.alive == false {
			//	request.response <- ResponseType{
			//		ResponseCode: ResponseCodeError,
			//	}
			//	fmt.Printf("error: cannot move into (%d, %d) the tile has dead content\n", x, y)
			//	continue
			//}

			var err error
			if state.content == Ambush {
				return fmt.Errorf("error: ambush not implemented")
				//err = tileServerAmbush(x, y, &content, request)
			} else if state.content == GridTraveler || request.requester.kind == Wildlife {
				err = handleEnteringOnTraveler(x, y, &state, &request, requestChannels)
			} else if state.content == Wildlife {
				err = handleEnteringOnWildlife(x, y, &state, &request, requestChannels, forcedMoveChannels)
			} else {
				log(x, y, state, "error: unknown content type")
			}
			if err != nil {
				return err
			}
		case <-ctx.Done():
			//fmt.Println("TileServer: ", x, y, " has been stopped")
			return nil
		}
	}
}

// TODO: delete
//func (gf *GridFieldSemaphores) at(x int, y int) chan struct{} {
//	return gf.fields[x][y].semaphore
//}
