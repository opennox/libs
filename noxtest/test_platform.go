package noxtest

import (
	"time"
)

// MockPlatform is a mock Platform implementation to allow test using ticks
type MockPlatform struct {
	T time.Duration
}

func (p *MockPlatform) Ticks() time.Duration {
	return p.T
}

func (p *MockPlatform) Sleep(dt time.Duration) {
	p.T += dt
}

func (p *MockPlatform) TimeSeed() int64 {
	return 0
}

func (p *MockPlatform) RandInt() int {
	return 0
}

func (p *MockPlatform) RandSeed(seed int64) {
}

func (p *MockPlatform) RandSeedTime() {
}
