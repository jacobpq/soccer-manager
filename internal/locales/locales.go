package locales

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/text/language"
)

//go:embed *.json
var Content embed.FS

type LangKey string

const Key LangKey = "lang"

var (
	translations  = make(map[string]map[string]string)
	supportedTags []language.Tag
	mutex         sync.RWMutex
	loaded        bool
)

func Init() error {
	mutex.Lock()
	defer mutex.Unlock()

	if loaded {
		return nil
	}

	entries, err := Content.ReadDir(".")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		langCode := strings.TrimSuffix(entry.Name(), ".json")

		data, err := Content.ReadFile(entry.Name())
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", entry.Name(), err)
		}

		var dict map[string]string
		if err := json.Unmarshal(data, &dict); err != nil {
			return fmt.Errorf("failed to parse %s: %w", entry.Name(), err)
		}

		translations[langCode] = dict

		tag := language.Make(langCode)
		if langCode == "en" {
			supportedTags = append([]language.Tag{tag}, supportedTags...)
		} else {
			supportedTags = append(supportedTags, tag)
		}
	}

	loaded = true
	log.Printf("Locales loaded: %v", supportedTags)
	return nil
}

func GetMatcher() language.Matcher {
	mutex.RLock()
	defer mutex.RUnlock()
	return language.NewMatcher(supportedTags)
}

func T(ctx context.Context, key string, args ...interface{}) string {
	lang, ok := ctx.Value(Key).(string)
	if !ok || lang == "" {
		lang = "en"
	}

	mutex.RLock()
	dict, ok := translations[lang]
	mutex.RUnlock()

	if !ok {
		mutex.RLock()
		dict = translations["en"]
		mutex.RUnlock()
	}

	val, ok := dict[key]
	if !ok {
		mutex.RLock()
		if enVal, enOk := translations["en"][key]; enOk {
			val = enVal
		} else {
			return key
		}
		mutex.RUnlock()
	}

	if len(args) > 0 {
		return fmt.Sprintf(val, args...)
	}

	return val
}
