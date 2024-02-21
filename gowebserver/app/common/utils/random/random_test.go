package random

import (
	"testing"
	"time"
)

func TestGenerateSubId(t *testing.T) {
	for num := 0; num < 10; num++ {
		t.Log(GenerateSubId(6))
		time.Sleep(1 * time.Millisecond)
	}
}

func TestGenValidateCode(t *testing.T) {
	for num := 0; num < 10; num++ {
		t.Log(GenValidateCode(6))
		time.Sleep(1 * time.Millisecond)
	}
}

func TestGenerateU64(t *testing.T) {
	for num := 0; num < 10; num++ {
		t.Log(GenerateU64())
	}
}
