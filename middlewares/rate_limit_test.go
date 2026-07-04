package middlewares

import (
	"testing"
	"time"
)

func TestInitRateLimiter(t *testing.T) {
	err := InitRateLimiter(":memory:")
	if err != nil {
		t.Fatalf("InitRateLimiter failed: %v", err)
	}

	if loginLimiter == nil {
		t.Fatal("loginLimiter should not be nil after init")
	}

	if loginLimiter.maxTries != 5 {
		t.Errorf("maxTries should be 5, got %d", loginLimiter.maxTries)
	}
}

func TestAPIRateLimitStruct(t *testing.T) {
	rl := &apiRateLimit{
		tokens:     10,
		maxTokens:  10,
		rate:       time.Second,
		lastAccess: time.Now(),
	}

	if rl.tokens != 10 {
		t.Errorf("tokens should be 10, got %d", rl.tokens)
	}

	if rl.maxTokens != 10 {
		t.Errorf("maxTokens should be 10, got %d", rl.maxTokens)
	}
}

func TestShutdown(t *testing.T) {
	err := InitRateLimiter(":memory:")
	if err != nil {
		t.Fatalf("InitRateLimiter failed: %v", err)
	}

	Shutdown()

	time.Sleep(100 * time.Millisecond)

	if shutdownCancel == nil {
		t.Log("Shutdown cancel function cleared")
	}
}