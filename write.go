package config

import (
	"bytes"
	"os"
	// "strings"
)

func SaveConfig(c *Config, filename string) (err error) {
	var f *os.File
	if f, err = os.Create(filename); err != nil {
		return err
	}
	buffer := bytes.NewBuffer(nil)
	for group, data := range c.data {
		if group != DEFAULT_GROUP {
			if _, err := buffer.WriteString("[" + group + "]" + lineBreak); err != nil {
				return err
			}
		}

		for key, val := range data {
			if key != " " {
				keyName := key
				if key[0] == '#' {
					keyName = "-"
				}

				// if strings.Contains(keyName, `=`) ||

				value := `"` + val + `"`

				if _, err := buffer.WriteString(keyName + " = " + value + lineBreak); err != nil {
					return err
				}

			}
		}
		if _, err := buffer.WriteString(lineBreak); err != nil {
			return err
		}
	}

	buffer.WriteTo(f)
	f.Close()
	return nil
}
