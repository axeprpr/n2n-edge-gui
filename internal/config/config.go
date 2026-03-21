package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	AddressModeDHCP   = "dhcp"
	AddressModeStatic = "static"
	fileName          = "n2n-gui.json"
	legacyININame     = "conf.ini"
)

type Config struct {
	Community     string `json:"community"`
	AddressMode   string `json:"addressMode"`
	Address       string `json:"address"`
	SupernodeHost string `json:"supernodeHost"`
	SupernodePort int    `json:"supernodePort"`
	MTU           int    `json:"mtu"`
	ExtraArgs     string `json:"extraArgs"`
}

func Default() Config {
	return Config{
		Community:     "community",
		AddressMode:   AddressModeDHCP,
		Address:       "",
		SupernodeHost: "127.0.0.1",
		SupernodePort: 7654,
		MTU:           1300,
		ExtraArgs:     "",
	}
}

func FilePath(baseDir string) string {
	return filepath.Join(baseDir, fileName)
}

func LegacyINIPath(baseDir string) string {
	return filepath.Join(baseDir, legacyININame)
}

func Load(baseDir string) (Config, error) {
	jsonPath := FilePath(baseDir)
	if _, err := os.Stat(jsonPath); err == nil {
		return loadJSON(jsonPath)
	}

	legacyPath := LegacyINIPath(baseDir)
	if _, err := os.Stat(legacyPath); err == nil {
		cfg, err := loadLegacyINI(legacyPath)
		if err != nil {
			return Config{}, err
		}
		return cfg, nil
	}

	return Default(), nil
}

func Save(baseDir string, cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	payload, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(FilePath(baseDir), append(payload, '\n'), 0o644)
}

func (c Config) Validate() error {
	c.Community = strings.TrimSpace(c.Community)
	c.AddressMode = strings.TrimSpace(c.AddressMode)
	c.Address = strings.TrimSpace(c.Address)
	c.SupernodeHost = strings.TrimSpace(c.SupernodeHost)

	if c.Community == "" {
		return errors.New("community is required")
	}

	if c.AddressMode != AddressModeDHCP && c.AddressMode != AddressModeStatic {
		return fmt.Errorf("addressMode must be %q or %q", AddressModeDHCP, AddressModeStatic)
	}

	if c.AddressMode == AddressModeStatic && c.Address == "" {
		return errors.New("address is required in static mode")
	}

	if c.SupernodeHost == "" {
		return errors.New("supernodeHost is required")
	}

	if c.SupernodePort < 1 || c.SupernodePort > 65535 {
		return errors.New("supernodePort must be between 1 and 65535")
	}

	if c.MTU < 500 || c.MTU > 9000 {
		return errors.New("mtu must be between 500 and 9000")
	}

	return nil
}

func (c Config) SupernodeAddress() string {
	return fmt.Sprintf("%s:%d", strings.TrimSpace(c.SupernodeHost), c.SupernodePort)
}

func (c Config) EdgeArgs() []string {
	args := []string{
		"-c", strings.TrimSpace(c.Community),
		"-l", c.SupernodeAddress(),
		"-M", strconv.Itoa(c.MTU),
	}

	if strings.TrimSpace(c.AddressMode) == AddressModeStatic {
		args = append(args, "-a", strings.TrimSpace(c.Address))
	} else {
		args = append(args, "-r", "-a", "dhcp:0.0.0.0")
	}

	if extra := strings.Fields(strings.TrimSpace(c.ExtraArgs)); len(extra) > 0 {
		args = append(args, extra...)
	}

	return args
}

func loadJSON(path string) (Config, error) {
	var cfg Config
	raw, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	if err := json.Unmarshal(raw, &cfg); err != nil {
		return Config{}, err
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func loadLegacyINI(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	cfg := Default()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])

		switch key {
		case "communityname":
			cfg.Community = value
		case "address":
			if value == "" {
				cfg.AddressMode = AddressModeDHCP
				cfg.Address = ""
			} else {
				cfg.AddressMode = AddressModeStatic
				cfg.Address = value
			}
		case "supernode":
			host, port, ok := strings.Cut(value, ":")
			if ok {
				cfg.SupernodeHost = strings.TrimSpace(host)
				parsed, err := strconv.Atoi(strings.TrimSpace(port))
				if err == nil {
					cfg.SupernodePort = parsed
				}
			} else {
				cfg.SupernodeHost = value
			}
		case "mtu":
			parsed, err := strconv.Atoi(value)
			if err == nil {
				cfg.MTU = parsed
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return Config{}, err
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
