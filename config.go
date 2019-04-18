package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func ParseConf() error {
	conf, err := ioutil.ReadFile("config.txt")
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`(?s)\/\*(.+?)\*\/`)
	conf = re.ReplaceAll(conf, []byte(""))

	kk := strings.Replace(string(conf), "\\", "\\\\", -1)
	conf = []byte(kk)
	err = json.Unmarshal(conf, &Config)
	if err != nil {
		conf = []byte(kk)
		if conf[0] == 239 {
			conf = conf[3:]
		}
		err = json.Unmarshal(conf, &Config)
	}
	cnfg_stat, _ := os.Stat("config.txt")
	Config_time = cnfg_stat.ModTime().Unix()

	return err
}

