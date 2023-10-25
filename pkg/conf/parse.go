// Package conf is parsing yaml, json, toml configuration files to go struct.
package conf

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Parse configuration files to struct, including yaml, toml, json, etc., and turn on listening for configuration file changes if fs is not empty
func Parse(configFile string, obj interface{}, fs ...func()) error {
	confFileAbs, err := filepath.Abs(configFile)
	if err != nil {
		return err
	}

	filePathStr, filename := filepath.Split(confFileAbs)
	ext := strings.TrimLeft(path.Ext(filename), ".")
	filename = strings.ReplaceAll(filename, "."+ext, "") // excluding suffix names

	viper.AddConfigPath(filePathStr) // path
	viper.SetConfigName(filename)    // file name
	viper.SetConfigType(ext)         // get the configuration type from the file name
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(obj)
	if err != nil {
		return err
	}

	if len(fs) > 0 {
		watchConfig(obj, fs...)
	}

	return nil
}

// listening for profile updates
func watchConfig(obj interface{}, fs ...func()) {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		err := viper.Unmarshal(obj)
		if err != nil {
			fmt.Println("viper.Unmarshal error: ", err)
		} else {
			for _, f := range fs {
				f()
			}
		}
	})
}

// Show print configuration information (hide sensitive fields)
func Show(obj interface{}, fields ...string) string {
	var out string

	data, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		fmt.Println("json.MarshalIndent error: ", err)
		return ""
	}

	buf := bufio.NewReader(bytes.NewReader(data))
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			break
		}
		fields = append(fields, `"dsn"`, `"password"`, `"pwd"`)

		out += hideSensitiveFields(line, fields...)
	}

	return out
}

func hideSensitiveFields(line string, fields ...string) string {
	for _, field := range fields {
		if strings.Contains(line, field) {
			index := strings.Index(line, field)
			if strings.Contains(line, "@") && strings.Contains(line, ":") {
				return replaceDSN(line)
			}
			return fmt.Sprintf("%s: \"******\",\n", line[:index+len(field)])
		}
	}

	// replace dsn
	if strings.Contains(line, "@") && strings.Contains(line, ":") {
		return replaceDSN(line)
	}

	return line
}

// replace dsn password
func replaceDSN(str string) string {
	data := []byte(str)
	start, end := 0, 0
	for k, v := range data {
		if v == ':' {
			start = k
		}
		if v == '@' {
			end = k
			break
		}
	}

	if start >= end {
		return str
	}

	return fmt.Sprintf("%s******%s", data[:start+1], data[end:])
}
