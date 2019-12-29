package helpers

import (
	"golang.org/x/time/rate"
	"sync"
)

type PhoneNumberRateLimiters struct {
	phoneNumbers map[string]*rate.Limiter
	mu           *sync.RWMutex
	r            rate.Limit
	b            int
}

func NewPhoneNumberRateLimiters(r float64, b int) *PhoneNumberRateLimiters {
	i := &PhoneNumberRateLimiters{
		phoneNumbers: make(map[string]*rate.Limiter),
		mu:           &sync.RWMutex{},
		r:            rate.Limit(r),
		b:            b,
	}

	return i
}

func (l *PhoneNumberRateLimiters) addPhoneNumber(phoneNumber string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	limiter := rate.NewLimiter(l.r, l.b)
	l.phoneNumbers[phoneNumber] = limiter

	return limiter
}

func (l *PhoneNumberRateLimiters) GetLimiter(phoneNumber string) *rate.Limiter {
	l.mu.Lock()
	limiter, exists := l.phoneNumbers[phoneNumber]
	if !exists {
		l.mu.Unlock()
		return l.addPhoneNumber(phoneNumber)
	}

	l.mu.Unlock()
	return limiter
}
