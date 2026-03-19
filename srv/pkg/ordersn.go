package pkg

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

func OrderSn() string {
	timeStr := time.Now().Format("20060102141512")
	uuidStr := uuid.New().String()

	return fmt.Sprintf("%s%s", timeStr, uuidStr[:8])
}
