package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)


type Znode struct {
	Created  time.Time `json:"created"`
	Prefixes []Zprefixes `json:"prefixes"`
}

type Zprefixes struct {
	Region   string `json:"region"`
	Hostname string `json:"hostname"`
	Location string `json:"location"`
	IPPrefix string `json:"ip_prefix"`
}

