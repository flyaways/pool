package pool

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

var (
	errClosed   = errors.New("pool is closed")
	errInvalid  = errors.New("invalid config")
	errRejected = errors.New("connection is nil. rejecting")
	errTargets  = errors.New("targets server is empty")
)

func init() {
	rand.NewSource(time.Now().UnixNano())
}

//Options pool options
type Options struct {
	lock sync.RWMutex
	//targets node
	targets *[]string
	//targets channel
	input chan *[]string

	//InitTargets init targets
	InitTargets []string
	// init connection
	InitCap int
	// max connections
	MaxCap       int
	DialTimeout  time.Duration
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Input is the input channel
func (o *Options) Input() chan<- *[]string {
	return o.input
}

// update targets
func (o *Options) update() {
	//init targets
	o.targets = &o.InitTargets

	go func() {
		for targets := range o.input {
			if targets == nil {
				continue
			}

			o.lock.Lock()
			o.targets = targets
			o.lock.Unlock()
		}
	}()

}

// NewOptions returns a new newOptions instance with sane defaults.
func NewOptions() *Options {
	o := &Options{}
	o.InitCap = 5
	o.MaxCap = 100
	o.DialTimeout = 5 * time.Second
	o.ReadTimeout = 5 * time.Second
	o.WriteTimeout = 5 * time.Second
	o.IdleTimeout = 60 * time.Second
	return o
}

// validate checks a Config instance.
func (o *Options) validate() error {
	if o.InitTargets == nil ||
		o.InitCap <= 0 ||
		o.MaxCap <= 0 ||
		o.InitCap > o.MaxCap ||
		o.DialTimeout == 0 ||
		o.ReadTimeout == 0 ||
		o.WriteTimeout == 0 {
		return errInvalid
	}
	return nil
}

//nextTarget next target implement load balance
func (o *Options) nextTarget() string {
	o.lock.RLock()
	defer o.lock.RUnlock()

	tlen := len(*o.targets)
	if tlen <= 0 {
		return ""
	}

	//rand server
	return (*o.targets)[rand.Int()%tlen]
}
