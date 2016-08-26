package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

/*
 a valid command line format includes: commands [parameters] [flags].
 parameters and flags maybe empty, the kinds of fields only be : string, numbers,array,bool

 example: agent ping hello world -t=3
 command: ping
 parameter: hello world
 flag: t

 mission file json format:
 {
	"reservedcmd":"ping",
	"reservedpara":"hello world",
	"t":3
 }
or
 {
	"reservedcmd":"daemon start",
	"c":true,
	"r":false
 }
or
{
	"reservedcmd":"ping",
	"reservedpara":"hello world",
	"tlscert":"./cert/client-cert.pem",
	"tlskey":"./cert/client-key.pem",
	"D":true,
	"host":"tcp://127.0.0.1:3567",
	"norsa":true
 }
 or
{
	"reservedcmd":"ping",
	"reservedpara":1234,
	"tlscert":"./cert/client-cert.pem",
	"tlskey":"./cert/client-key.pem",
	"D":true,
	"host":"tcp://127.0.0.1:3567",
	"norsa":true
 }
or
{
	"reservedcmd":"ping",
	"reservedpara":["hello" ,2.5],
	"tlscert":"./cert/client- cert.pem",
	"tlskey":"./cert/client-key.pem",
	"D":true,
	"host":["tcp://127.0.0.1:3567","127.0.0. 1",3567],
	"norsa":true
 }
***todo: multiple missions. like:

[
	{
		"reservedcmd":"ping",
		"reservedpara":[1,2],
		"t":3
	},
	{
	"reservedcmd":"daemon start",
	"c":true,
	"r":false
 	}
]
*/

const (
	mission_command_field_name = "reservedcmd"
	mission_para_field_name    = "reservedpara"
)

type mission map[string]interface{}

func NewMission() mission {
	m := mission(make(map[string]interface{}))
	return m
}
func (m mission) ToArgs() ([]string, error) {
	args := []string{}

	if cmd, exist := m[mission_command_field_name]; exist {
		if cmdstr, ok := cmd.(string); ok {
			args = append(args, cmdstr)
		} else {
			return nil, fmt.Errorf(`reservedcmd is not string type`)
		}
	} else {
		return nil, fmt.Errorf(`reservedcmd must exist`)
	}

	if para, exist := m[mission_para_field_name]; exist {
		val := reflect.ValueOf(para)
		valtype := val.Type().Kind()

		if valtype == reflect.Array || valtype == reflect.Slice {
			for idx := 0; idx != val.Len(); idx++ {
				args = append(args, fmt.Sprintf("%v", val.Index(idx).Interface()))
			}
		} else {
			args = append(args, fmt.Sprintf("%v", para))
		}
	}

	for k, v := range m {
		if k == mission_command_field_name || k == mission_para_field_name {
			continue
		}
		sign := "-"

		if len(k) > 1 {
			sign = "--"
		}

		val := reflect.ValueOf(v)
		valtype := val.Type().Kind()

		if valtype == reflect.Array || valtype == reflect.Slice {
			tmpslice := []string{}
			for idx := 0; idx != val.Len(); idx++ {
				tmpslice = append(tmpslice, fmt.Sprintf("%v", val.Index(idx).Interface()))
			}
			v = strings.Join(tmpslice, ",")
		}

		if str, ok := v.(string); ok {
			if strings.ContainsAny(str, "\t\n\v\f\r ") {
				v = fmt.Sprintf(`"%s"`, str)
			}
		}
		args = append(args, fmt.Sprintf("%s%s=%v", sign, k, v))
	}

	return args, nil
}
func (m mission) Read(name string) error {
	buf, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, &m)
	return err
}
