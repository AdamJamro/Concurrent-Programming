package models

import (
	"fmt"
	"grid-travelers-v1/config"
)

func (t *TraceSequenceType) PrintTrace() {
	for i := 0; i < t.len; i++ {
		fmt.Printf(
			"%.9f %d %d %d %c\n",
			float64(t.data[i].timeStamp)/1e9,
			t.data[i].id,
			t.data[i].pos.x,
			t.data[i].pos.y,
			t.data[i].symbol,
		)
	}
}

func (t *TraceSequenceType) add(trace TraceType) error {
	if t.len >= config.MaxSteps {
		return fmt.Errorf("error: TraceSequenceType is full")
	}
	t.data[t.len] = trace
	t.len++
	return nil
}
