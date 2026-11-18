package main


type Check struct {
	Name string `json:"name"`

	Slo  float64
	Rise int
	Fall int

	Polled    int
	Suceeded  int

	Healthy   bool   `json:"healthy"`
	Enabled   bool   `json:"enabled"`
	IsDefault bool   `json:"is_default"`
	Statement string `json:"statement"`

	Every     time.Duration
	Threshold time.Duration

	LastFailedAt  time.Time `json:"last_failed_at"`
	LastCheckedAt time.Time `json:"last_checked_at"`
}

type SharedDefaultConfig struct {
	Checks    []Check
	Threshold Threshold
}

type Defaults interface {
	SharedDefaultChecks() []Check
	SharedDefaultThreshold() Threshold
}

func (sc *SharedDefaultConfig) SharedDefaultChecks() []Check {
	return sc.Checks
}

func (sc *SharedDefaultConfig) SharedDefaultThreshold() Threshold {
	return sc.Threshold
}


startchecks

customchecks

defaultchecks

registerchecks
