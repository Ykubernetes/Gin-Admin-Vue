package conf

import (
	"fmt"
	"strconv"
)

type Redis struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Password    string `yaml:"password"`
	SelectDb    int    `yaml:"selectDb"`
	PolSize     int    `yaml:"PolSize"`
	MinIdleConn int    `yaml:"minIdleConn"`
}

func (r *Redis) GetHost() string {
	return fmt.Sprintf(r.Host + ":" + strconv.Itoa(r.Port))
}
