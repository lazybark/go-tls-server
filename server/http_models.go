package server

import "time"

var (
	//ResultError represents error with status, code and human-readable error
	//
	//Format: `{"success":false,"code":"%v","error":"%s"}`
	ResultError = `{"success":false,"code":"%v","error":"%s"}`

	//ResultString represents one-string result.
	//Useful in case of one-string result from a method
	ResultString = `{"success":true,"result":"%v"}`

	ResultJSON = `{"success":true,"result":%v}`
)

type ServerStatsOutput struct {
	Received    int       `json:"bytes_received"`
	Sent        int       `json:"bytes_sent"`
	Errors      int       `json:"total_errors"`
	Connections int       `json:"connections"`
	Started     time.Time `json:"server_started"`
}
