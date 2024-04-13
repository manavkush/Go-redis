package client

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"testing"
)

func TestClient(t *testing.T) {

	client, err := NewClient("localhost:5001")
	if err != nil {
		slog.Error("Error creation new client", "err", err)
	}
	for i := range 5 {
		fmt.Println("SET =>", fmt.Sprintf("bar-%d", i))
		if err := client.Set(context.Background(), fmt.Sprintf("foo-%d", i), fmt.Sprintf("bar-%d", i)); err != nil {
			log.Fatal(err)
		}
		val, err := client.Get(context.Background(), fmt.Sprintf("foo-%d", i))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("GET =>", val)
	}
}
