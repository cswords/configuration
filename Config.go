package configuration

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config - A config represents the root node of the structures.
type Config struct {
	Server Server `yaml:"server"`
}

// Server - A server represents a concept like an HTTP virtual server having its own port and containing zero to multiple routing rules.
type Server struct {
	Port    string   `yaml:"port"`
	Routers []Router `yaml:"routers"`
}

// Router - A router contains a prefix-defined routing rule, and underlying middlewares and handlers.
type Router struct {
	Prefix      string       `yaml:"prefix"`
	Middlewares []Middleware `yaml:"middlewares"`
	Handlers    []Handler    `yaml:"handlers"`
}

// Middleware - A middlware represents an HTTP middleware having a type and its own configuration.
type Middleware struct {
	Type   string            `yaml:"type"`
	Config map[string]string `yaml:"config"`
}

// Handler - A handler represents an HTTP handler listening to a path and having a type and its own configuration.
type Handler struct {
	Path   string            `yaml:"path"`
	Type   string            `yaml:"type"`
	Config map[string]string `yaml:"config"`
}

func (c *Config) loadBinary(b []byte) *Config {

	err := yaml.Unmarshal(b, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

// LoadConfig use the given loader functions to retrieve a config data.
// Multiple loaders are tried one by one.
// FromLocal is an example loader function.
func LoadConfig(location string, loaders ...(func(string) []byte)) *Config {
	c := &(Config{})
	var b []byte
	if len(loaders) == 0 {
		loaders = []func(string) []byte{FromLocal}
	}
	for _, loader := range loaders {
		b = loader(location)
		if b != nil {
			break
		}
	}
	return c.loadBinary(b)
}

// FromLocal loads a configuration file on a local location into bytes.
func FromLocal(loc string) []byte {

	if !strings.HasPrefix(loc, "./") {
		return nil
	}

	folder, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		log.Println("load binary from local err", err)
		return nil
	}
	data, err := ioutil.ReadFile(folder + strings.TrimPrefix(loc, "."))
	if err != nil {
		log.Println("load binary from local  err", err)
		return nil
	}
	return data
}
