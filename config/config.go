package config

import "time"

const (
	LogDir                       = "./logs"
	LogFileFormat                = "2006-01-02 15:04:05"
	WebAuthMaxAttempts           = 10
	WebAuthAttemptResetDuration  = 5 * time.Minute
	WebAuthSessionDurationInHour = 7 * 24
	WebServerPort                = "35353"
)
