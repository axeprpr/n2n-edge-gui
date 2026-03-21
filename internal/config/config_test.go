package config

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestLoadLegacyINI(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	content := "communityname = alpha\naddress = 10.0.0.8\nsupernode = edge.example.com:7777\nmtu = 1400\n"
	if err := os.WriteFile(filepath.Join(dir, "conf.ini"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Community != "alpha" {
		t.Fatalf("unexpected community: %q", cfg.Community)
	}
	if cfg.AddressMode != AddressModeStatic {
		t.Fatalf("unexpected address mode: %q", cfg.AddressMode)
	}
	if cfg.SupernodeHost != "edge.example.com" || cfg.SupernodePort != 7777 {
		t.Fatalf("unexpected supernode: %s:%d", cfg.SupernodeHost, cfg.SupernodePort)
	}
	if cfg.MTU != 1400 {
		t.Fatalf("unexpected mtu: %d", cfg.MTU)
	}
}

func TestEdgeArgs(t *testing.T) {
	t.Parallel()

	cfg := Config{
		Community:     "alpha",
		AddressMode:   AddressModeDHCP,
		SupernodeHost: "supernode.local",
		SupernodePort: 1234,
		MTU:           1300,
		ExtraArgs:     "-f -v",
	}

	got := cfg.EdgeArgs()
	want := []string{"-c", "alpha", "-l", "supernode.local:1234", "-M", "1300", "-r", "-a", "dhcp:0.0.0.0", "-f", "-v"}
	if !slices.Equal(got, want) {
		t.Fatalf("unexpected args:\n got: %#v\nwant: %#v", got, want)
	}
}
