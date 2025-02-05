package cmd

import (
	"bytes"
	"errors"
	"github.com/saichler/shared/go/share/interfaces"
	"reflect"
	"strconv"
	"strings"
)

type Command interface {
	Name() string
	Help() string
	Run(log interfaces.ILogger) string
}

func argNames(cmd Command) string {
	v := reflect.ValueOf(cmd).Elem()
	buff := bytes.Buffer{}
	buff.WriteString("Available arguments:\n")
	for i := 0; i < v.NumField(); i++ {
		name := v.Type().Field(i).Name
		kind := v.Field(i).Kind().String()
		buff.WriteString("   ")
		buff.WriteString(name)
		buff.WriteString(" - ")
		buff.WriteString(kind)
		buff.WriteString("\n")
	}
	return buff.String()
}

func fillArgs(cmd Command, osargs []string) error {
	v := reflect.ValueOf(cmd).Elem()
	args := osargs[2:]
	for _, arg := range args {
		split := strings.Split(arg, "=")
		if len(split) != 2 {
			return errors.New("Invalid argument format: " + arg + " for command " + osargs[1])
		}
		fld := v.FieldByName(split[0])
		if !fld.IsValid() {
			return errors.New("Invalid argument name: " + split[0] + " for command " + osargs[1] + "\n" + argNames(cmd))
		}
		if fld.Kind() == reflect.String {
			fld.Set(reflect.ValueOf(split[1]))
		} else if fld.Kind() == reflect.Int {
			i, err := strconv.Atoi(split[1])
			if err != nil {
				return errors.New("Invalid argument int: " + split[1] + " for command " + osargs[1])
			}
			fld.Set(reflect.ValueOf(i))
		}
	}
	return nil
}
