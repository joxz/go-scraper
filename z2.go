package main

import (
	"time"
)

type ZscalerRegion struct {
	ZRegion string `json:"region"`
	Element string `json:"-"`
}

type ZscalerNode struct {
	ZscalerRegion
	Cidr string `json:"ip_prefix"`
}

type Zscaler struct {
	Created  time.Time     `json:"created"`
	Prefixes []ZscalerNode `json:"prefixes"`
}
