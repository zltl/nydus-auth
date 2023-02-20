package id

import (
	"testing"
	"time"

	"log"
)

func TestCurTime(t *testing.T) {
	cur := time.Now()
	log.Printf("cur: %v, mili: %v\n", cur, time.Now().UnixMilli())
}
