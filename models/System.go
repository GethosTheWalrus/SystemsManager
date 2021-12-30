package models

import "time"

type HostSystem struct {
	Id        int     `json:"id" gorm:"primary_key"`
	CpuId     int     `json:"cpu_id"`
	Cpu       Cpu     `json:"cpu" gorm:"foreignKey:CpuId"`
	MemoryId  int     `json:"memory_id"`
	Memory    Memory  `json:"memory" gorm:"foreignKey:MemoryId"`
	NetworkId int     `json:"network_id"`
	Network   Network `json:"network" gorm:"foreignKey:NetworkId"`
	GeneralId int     `json:"general_id"`
	General   General `json:"general" gorm:"foreignKey:GeneralId"`
}

type Cpu struct {
	Id     int     `json:"id" gorm:"primary_key"`
	Idle   float64 `json:"idle"`
	System float64 `json:"system"`
	Total  float64 `json:"total"`
	User   float64 `json:"user"`
}

type Memory struct {
	Id     int    `json:"id" gorm:"primary_key"`
	Cached uint64 `json:"cached"`
	Free   uint64 `json:"free"`
	Total  uint64 `json:"total"`
	Used   uint64 `json:"used"`
}

type Network struct {
	Id          int    `json:"id" gorm:"primary_key"`
	Hostname    string `json:"hostname"`
	PreferredIp string `json:"preferred_ip"`
}

type General struct {
	Id              int       `json:"id" gorm:"primary_key"`
	OperatingSystem string    `json:"operating_system"`
	LastSeen        time.Time `json:"last_seen"`
}
