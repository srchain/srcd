/**
 * Read the configuration file
 *
 * @copyright           (C) 2014  widuu
 * @lastmodify          2014-2-22
 * @website		http://www.widuu.com
 *
 */

package goini

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Config struct {
	filepath string                       //your ini file path directory+file
	conf     map[string]map[string]string //configuration information slice
}

//Create an empty configuration file
func SetConfig(filepath string) (*Config, error) {
	_, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}
	c := new(Config)
	c.filepath = filepath
	c.conf = make(map[string]map[string]string)

	return c, nil
}

//To obtain corresponding value of the key values
func (c *Config) GetValue(section, name string) (string, error) {
	c.ReadList()
	conf := c.ReadList()

	for key, value := range conf {
		if key == section {
			return value[name], nil
		}
	}

	return "", errors.New("no value")
}

//List all the configuration file
func (c *Config) ReadList() map[string]map[string]string {
	file, err := os.Open(c.filepath)
	if err != nil {
		CheckErr(err)
	}
	defer file.Close()
	var section string
	buf := bufio.NewReader(file)
	for {
		l, err := buf.ReadString('\n')
		line := strings.TrimSpace(l)
		if err != nil {
			if err != io.EOF {
				CheckErr(err)
			}
			if len(line) == 0 {
				break
			}
		}
		switch {
		case len(line) == 0:
		case string(line[0]) == "#": //增加配置文件备注
		case line[0] == '[' && line[len(line)-1] == ']':
			section = strings.TrimSpace(line[1 : len(line)-1])
			c.conf[section] = make(map[string]string)
		default:
			i := strings.IndexAny(line, "=")
			value := strings.TrimSpace(line[i+1 : len(line)])
			c.conf[section][strings.TrimSpace(line[0:i])] = value
		}

	}

	return c.conf
}

func CheckErr(err error) string {
	if err != nil {
		return fmt.Sprintf("Error is :'%s'", err.Error())
	}
	return "Notfound this error"
}
