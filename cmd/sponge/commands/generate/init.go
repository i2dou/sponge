package generate

import (
	"embed"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/i2dou/sponge/pkg/gofile"
	"github.com/i2dou/sponge/pkg/replacer"
)

const warnSymbol = "⚠ "

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Replacers replacer name
var Replacers = map[string]replacer.Replacer{}

// Template information
type Template struct {
	Name     string
	FS       embed.FS
	FilePath string
}

// Init initializing the template
func Init(name string, filepath string) error {
	// determine if the template file exists, if not, prompt to initialize first
	if !gofile.IsExists(filepath) {
		if isShowCommand() {
			return nil
		}
		return fmt.Errorf("%s not yet initialized, run the command \"sponge init\"", warnSymbol)
	}

	var err error
	if _, ok := Replacers[name]; ok {
		panic(fmt.Sprintf("template name \"%s\" already exists", name))
	}
	Replacers[name], err = replacer.New(filepath)
	if err != nil {
		return err
	}

	return nil
}

// InitFS initializing th FS templates
func InitFS(name string, filepath string, fs embed.FS) {
	var err error
	if _, ok := Replacers[name]; ok {
		panic(fmt.Sprintf("template name \"%s\" already exists", name))
	}
	Replacers[name], err = replacer.NewFS(filepath, fs)
	if err != nil {
		panic(err)
	}
}

func isShowCommand() bool {
	l := len(os.Args)

	// sponge
	if l == 1 {
		return true
	}

	// sponge init or sponge -h
	if l == 2 {
		if os.Args[1] == "init" || os.Args[1] == "-h" {
			return true
		}
		return false
	}
	if l > 2 {
		return strings.Contains(strings.Join(os.Args[:3], ""), "init")
	}

	return false
}
