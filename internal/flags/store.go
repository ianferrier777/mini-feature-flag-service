package flags

import (
	"sync"
	"encoding/json"
	"io/ioutil"
	"os"
)

const flagFile = "flags.json"

var (
	flagStore = make(map[string]Flag)
	mu        sync.RWMutex
)

// InitFlags initializes the store with some sample flags
func InitFlags() {
	mu.Lock()
	defer mu.Unlock()

	flagStore["new-homepage"] = Flag{
		Name:          "new-homepage",
		Enabled:       false,
		TargetUsers:   []string{"123", "456"},
		TargetRegions: []string{"us", "ca"},
	}

	flagStore["beta-dashboard"] = Flag{
		Name:    "beta-dashboard",
		Enabled: true, // fallback enabled for everyone
	}
}


// GetFlag returns a flag by name, and a bool indicating if it was found
func GetFlag(name string) (Flag, bool) {
	mu.RLock()
	defer mu.RUnlock()

	flag, found := flagStore[name]
	return flag, found
}

// SetFlag creates or updates a flag (used by POST)
func SetFlag(name string, enabled bool, users []string, regions []string) {
	mu.Lock()
	defer mu.Unlock()
	flagStore[name] = Flag{
		Name:          name,
		Enabled:       enabled,
		TargetUsers:   users,
		TargetRegions: regions,
	}
}

// UpdateFlag updates an existing flag; returns false if not found
func UpdateFlag(name string, enabled bool, users []string, regions []string) bool {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := flagStore[name]; !exists {
		return false
	}

	flagStore[name] = Flag{
		Name:          name,
		Enabled:       enabled,
		TargetUsers:   users,
		TargetRegions: regions,
	}
	return true
}

func SaveToFile() error {
	mu.RLock()
	defer mu.RUnlock()

	data, err := json.MarshalIndent(flagStore, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(flagFile, data, 0644)
}

func LoadFromFile() error {
	mu.Lock()
	defer mu.Unlock()

	if _, err := os.Stat(flagFile); os.IsNotExist(err) {
		return nil // no file to load
	}

	data, err := ioutil.ReadFile(flagFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &flagStore)
}

func GetAllFlags() []Flag {
	mu.RLock()
	defer mu.RUnlock()

	flagsList := make([]Flag, 0, len(flagStore))
	for _, flag := range flagStore {
		flagsList = append(flagsList, flag)
	}
	return flagsList
}

func DeleteFlag(name string) bool {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := flagStore[name]; !exists {
		return false
	}
	delete(flagStore, name)
	return true
}