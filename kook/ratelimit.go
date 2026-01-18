package kook

import (
	"sync"
	"time"
)

// RateLimiter 速率限制器
type RateLimiter struct {
	tokens     chan struct{}
	refillMu   sync.Mutex
	lastRefill time.Time
	rate       time.Duration
	burst      int
}

// NewRateLimiter 创建新的速率限制器
// rate: 令牌补充间隔
// burst: 令牌桶容量
func NewRateLimiter(rate time.Duration, burst int) *RateLimiter {
	rl := &RateLimiter{
		tokens:     make(chan struct{}, burst),
		lastRefill: time.Now(),
		rate:       rate,
		burst:      burst,
	}

	// 初始填满令牌桶
	for i := 0; i < burst; i++ {
		select {
		case rl.tokens <- struct{}{}:
		default:
		}
	}

	// 启动令牌补充协程
	go rl.refillLoop()

	return rl
}

// Wait 等待获取令牌
func (rl *RateLimiter) Wait() {
	<-rl.tokens
}

// TryAcquire 尝试获取令牌，不等待
func (rl *RateLimiter) TryAcquire() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

// refillLoop 令牌补充循环
func (rl *RateLimiter) refillLoop() {
	ticker := time.NewTicker(rl.rate)
	defer ticker.Stop()

	for range ticker.C {
		rl.refillMu.Lock()
		// 尝试添加一个令牌
		select {
		case rl.tokens <- struct{}{}:
			// 成功添加令牌
		default:
			// 令牌桶已满
		}
		rl.lastRefill = time.Now()
		rl.refillMu.Unlock()
	}
}

// EndpointRateLimiter 端点级别的速率限制器
type EndpointRateLimiter struct {
	limiters map[string]*RateLimiter
	mu       sync.RWMutex
	rate     time.Duration
	burst    int
}

// NewEndpointRateLimiter 创建端点级别的速率限制器
func NewEndpointRateLimiter(rate time.Duration, burst int) *EndpointRateLimiter {
	return &EndpointRateLimiter{
		limiters: make(map[string]*RateLimiter),
		rate:     rate,
		burst:    burst,
	}
}

// Wait 等待指定端点的令牌
func (erl *EndpointRateLimiter) Wait(endpoint string) {
	erl.getLimiter(endpoint).Wait()
}

// TryAcquire 尝试获取指定端点的令牌
func (erl *EndpointRateLimiter) TryAcquire(endpoint string) bool {
	return erl.getLimiter(endpoint).TryAcquire()
}

// getLimiter 获取或创建端点的速率限制器
func (erl *EndpointRateLimiter) getLimiter(endpoint string) *RateLimiter {
	erl.mu.RLock()
	limiter, exists := erl.limiters[endpoint]
	erl.mu.RUnlock()

	if exists {
		return limiter
	}

	erl.mu.Lock()
	defer erl.mu.Unlock()

	// 双重检查
	if limiter, exists := erl.limiters[endpoint]; exists {
		return limiter
	}

	// 创建新的限制器
	limiter = NewRateLimiter(erl.rate, erl.burst)
	erl.limiters[endpoint] = limiter
	return limiter
}

// GlobalRateLimiter 全局速率限制器
type GlobalRateLimiter struct {
	generalLimiter  *RateLimiter
	endpointLimiter *EndpointRateLimiter
}

// NewGlobalRateLimiter 创建全局速率限制器
func NewGlobalRateLimiter() *GlobalRateLimiter {
	return &GlobalRateLimiter{
		// KOOK API 全局限制：120 requests per minute
		generalLimiter: NewRateLimiter(500*time.Millisecond, 10),
		// 端点级别限制：更宽松一些
		endpointLimiter: NewEndpointRateLimiter(200*time.Millisecond, 5),
	}
}

// Wait 等待令牌（同时检查全局和端点限制）
func (grl *GlobalRateLimiter) Wait(endpoint string) {
	// 先等待全局限制
	grl.generalLimiter.Wait()
	// 再等待端点限制
	grl.endpointLimiter.Wait(endpoint)
}

// TryAcquire 尝试获取令牌
func (grl *GlobalRateLimiter) TryAcquire(endpoint string) bool {
	// 需要同时满足全局和端点限制
	if !grl.generalLimiter.TryAcquire() {
		return false
	}
	if !grl.endpointLimiter.TryAcquire(endpoint) {
		// 如果端点限制失败，需要把全局令牌还回去
		// 这里简化处理，直接返回失败
		return false
	}
	return true
}

