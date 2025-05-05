package models

import "grid-travelers-v4/config"

type FieldSemaphore struct {
	semaphore chan struct{}
}

type GridFieldSemaphores struct {
	fields [][]FieldSemaphore
}

func (gf *GridFieldSemaphores) InitGridFields() {
	gf.fields = make([][]FieldSemaphore, config.GridWidth)
	for i := range gf.fields {
		gf.fields[i] = make([]FieldSemaphore, config.GridHeight)
	}

	for i := range gf.fields {
		for j := range gf.fields[i] {
			gf.fields[i][j].semaphore = make(chan struct{}, config.GridFieldCapacity)
			for k := 0; k < config.GridFieldCapacity; k++ {
				gf.fields[i][j].semaphore <- struct{}{}
			}
		}
	}
}

func (gf *GridFieldSemaphores) at(x int, y int) chan struct{} {
	return gf.fields[x][y].semaphore
}
