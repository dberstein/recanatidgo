package rate

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	ginratelimit "github.com/ljahier/gin-ratelimit"
)

type rateLimiter struct {
	tb *ginratelimit.TokenBucket
}

func NewRateLimiter(rateSpec string) (*rateLimiter, error) {
	rateTmp := strings.Split(rateSpec, "/")
	rate, err := strconv.Atoi(rateTmp[0])
	if err != nil {
		return nil, err
	}
	if rate <= 0 {
		return nil, errors.New("rate has to be greater than zero")
	}

	ttl, err := time.ParseDuration(rateTmp[1])
	if err != nil {
		return nil, err
	}
	if ttl <= 0 {
		return nil, errors.New("rate's TTL has to be greater than zero")
	}

	log.Printf("Rate limiter: %d per %s", rate, ttl)
	return &rateLimiter{
		tb: ginratelimit.NewTokenBucket(rate, ttl),
	}, nil
}

func (rl *rateLimiter) GetTokenBucket() *ginratelimit.TokenBucket {
	return rl.tb
}
