package v1

import (
	"time"
)

type DataServer struct {
	UUID       int64
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	Port       int32     `json:"port"`
	AccessKey  string    `json:"accessKey"`
	SecretKey  string    `json:"secretKey"`
	Token      string    `json:"token"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	LastSeenAt time.Time `json:"lastSeenAt"`
	Parked     bool      `json:"Parked"`
}

type DataServerList struct {
	Items []*DataServer `json:"items"`
}
