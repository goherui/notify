package pkg

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func OrderSn() string {
	timeStr := time.Now().Format("20060102150405")
	uuidStr := uuid.New().String()[:4]
	return fmt.Sprintf("%s%s", timeStr, uuidStr)
}
