package api

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/2qif49lt/agent/pkg/system"
	"github.com/docker/libtrust"
)

func formGroup(key string, start, last int) string {
	parts := strings.Split(key, "/")
	groupType := parts[0]
	var ip string
	if len(parts) > 1 {
		ip = parts[0]
		groupType = parts[1]
	}
	group := strconv.Itoa(start)
	if start != last {
		group = fmt.Sprintf("%s-%d", group, last)
	}
	if ip != "" {
		group = fmt.Sprintf("%s:%s->%s", ip, group, group)
	}
	return fmt.Sprintf("%s/%s", group, groupType)
}

// LoadOrCreateTrustKey attempts to load the libtrust key at the given path,
// otherwise generates a new one
func LoadOrCreateTrustKey(trustKeyPath string) (libtrust.PrivateKey, error) {
	err := system.MkdirAll(filepath.Dir(trustKeyPath), 0700)
	if err != nil {
		return nil, err
	}
	trustKey, err := libtrust.LoadKeyFile(trustKeyPath)
	if err == libtrust.ErrKeyFileDoesNotExist {
		trustKey, err = libtrust.GenerateECP256PrivateKey()
		if err != nil {
			return nil, fmt.Errorf("Error generating key: %s", err)
		}
		if err := libtrust.SaveKey(trustKeyPath, trustKey); err != nil {
			return nil, fmt.Errorf("Error saving key file: %s", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("Error loading key file %s: %s", trustKeyPath, err)
	}
	return trustKey, nil
}

// LoadOrCreateTrustKey attempts to load the libtrust key at the given path,
// otherwise generates a new one
func LoadSignPubKey(pubkey string) (libtrust.PublicKey, error) {
	return libtrust.LoadPublicKeyFile(pubkey)
}
