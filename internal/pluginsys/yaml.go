// Package plugin — yaml.go implements Layer 0: declarative YAML plugins.
//
// A YAML plugin is a single .yaml file that defines webhooks, simple filter rules,
// content injection, and settings — with zero code. Example:
//
//	id: my-webhook-forwarder
//	title: Webhook 转发
//	version: 1.0.0
//	icon: i-tabler-webhook
//	author: admin
//	description: 将文章事件转发到第三方 Webhook
//	settings:
//	  - key: webhook_url
//	    label: Webhook URL
//	    type: string
//	    required: true
//	webhooks:
//	  - url: "{{settings.webhook_url}}"
//	    events: ["post.published", "comment.created"]
//	filters:
//	  - event: comment.create
//	    rules:
//	      - field: content
//	        min_length: 5
//	        message: 评论太短
package pluginsys

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"gopkg.in/yaml.v3"
)

// ─── YAML plugin definition ─────────────────────────────────────────────────

// YAMLPlugin is the in-memory representation of a declarative YAML plugin.
type YAMLPlugin struct {
	ID          string         `yaml:"id"`
	Title       string         `yaml:"title"`
	Version     string         `yaml:"version"`
	Icon        string         `yaml:"icon"`
	Author      string         `yaml:"author"`
	Description string         `yaml:"description"`
	TrustLevel  string         `yaml:"trust_level"`
	Priority    int            `yaml:"priority"`
	Settings    []SettingField `yaml:"settings"`
	Webhooks    []YAMLWebhook  `yaml:"webhooks"`
	Filters     []YAMLFilter   `yaml:"filters"`
	CSS         string         `yaml:"css"`
	PublicJS    string         `yaml:"public_js"`
	AdminJS     string         `yaml:"admin_js"`
	Navigation  []NavigationDef `yaml:"navigation"`
}

// YAMLWebhook is a declarative webhook forwarding rule.
type YAMLWebhook struct {
	URL     string            `yaml:"url"`
	Events  []string          `yaml:"events"`
	Headers map[string]string `yaml:"headers"`
}

// YAMLFilter is a declarative filter with simple matching rules.
type YAMLFilter struct {
	Event string     `yaml:"event"` // e.g. "comment.create"
	Rules []YAMLRule `yaml:"rules"`
}

// YAMLRule is a single validation rule applied to a field in filter data.
type YAMLRule struct {
	Field        string   `yaml:"field"`                   // data field to check
	MinLength    int      `yaml:"min_length,omitempty"`    // minimum string length
	MaxLength    int      `yaml:"max_length,omitempty"`    // maximum string length
	BlockedWords []string `yaml:"blocked_words,omitempty"` // reject if any word found
	Regex        string   `yaml:"regex,omitempty"`         // reject if matches
	NotRegex     string   `yaml:"not_regex,omitempty"`     // reject if does NOT match
	Message      string   `yaml:"message"`                 // abort message
}

// ─── YAML plugin registry ────────────────────────────────────────────────────

var (
	yamlMu      sync.RWMutex
	yamlPlugins = make(map[string]*YAMLPlugin)
)

// LoadYAMLPlugins loads all .yaml/.yml files from dir as declarative plugins.
func LoadYAMLPlugins(ctx context.Context, dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			g.Log().Warningf(ctx, "[yaml-plugin] read dir %s: %v", dir, err)
		}
		return
	}

	count := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext != ".yaml" && ext != ".yml" {
			continue
		}
		path := filepath.Join(dir, e.Name())
		if err := loadOneYAMLPlugin(ctx, path); err != nil {
			g.Log().Warningf(ctx, "[yaml-plugin] load %s: %v", path, err)
		} else {
			count++
		}
	}
	g.Log().Infof(ctx, "[yaml-plugin] %d declarative plugin(s) loaded from %s", count, dir)
}

func loadOneYAMLPlugin(ctx context.Context, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var yp YAMLPlugin
	if err := yaml.Unmarshal(data, &yp); err != nil {
		return fmt.Errorf("parse: %w", err)
	}
	if yp.ID == "" {
		return fmt.Errorf("missing 'id' field")
	}
	if yp.Priority == 0 {
		yp.Priority = 10
	}

	// Pre-compile regex rules
	for i, f := range yp.Filters {
		for j, r := range f.Rules {
			if r.Regex != "" {
				if _, err := regexp.Compile(r.Regex); err != nil {
					return fmt.Errorf("filter[%d].rule[%d] bad regex: %w", i, j, err)
				}
			}
			if r.NotRegex != "" {
				if _, err := regexp.Compile(r.NotRegex); err != nil {
					return fmt.Errorf("filter[%d].rule[%d] bad not_regex: %w", i, j, err)
				}
			}
		}
	}

	yamlMu.Lock()
	yamlPlugins[yp.ID] = &yp
	yamlMu.Unlock()

	g.Log().Infof(ctx, "[yaml-plugin] loaded %s (%s v%s)", yp.ID, yp.Title, yp.Version)
	return nil
}

// ─── YAML filter execution ──────────────────────────────────────────────────

// RunYAMLFilters runs all declarative YAML plugin filters for the given event.
// It modifies data in place and returns an error if any rule aborts.
func RunYAMLFilters(ctx context.Context, event string, data map[string]any) error {
	yamlMu.RLock()
	defer yamlMu.RUnlock()

	for _, yp := range yamlPlugins {
		for _, f := range yp.Filters {
			if f.Event != event {
				continue
			}
			if err := runYAMLRules(yp, f.Rules, data); err != nil {
				return err
			}
		}
	}
	return nil
}

func runYAMLRules(yp *YAMLPlugin, rules []YAMLRule, data map[string]any) error {
	for _, r := range rules {
		val, _ := data[r.Field].(string)

		// min_length
		if r.MinLength > 0 && len(strings.TrimSpace(val)) < r.MinLength {
			return fmt.Errorf("[%s] %s", yp.ID, r.Message)
		}

		// max_length
		if r.MaxLength > 0 && len([]rune(val)) > r.MaxLength {
			return fmt.Errorf("[%s] %s", yp.ID, r.Message)
		}

		// blocked_words
		if len(r.BlockedWords) > 0 {
			lower := strings.ToLower(val)
			for _, w := range r.BlockedWords {
				if strings.Contains(lower, strings.ToLower(w)) {
					return fmt.Errorf("[%s] %s", yp.ID, r.Message)
				}
			}
		}

		// regex — reject if matches
		if r.Regex != "" {
			if matched, _ := regexp.MatchString(r.Regex, val); matched {
				return fmt.Errorf("[%s] %s", yp.ID, r.Message)
			}
		}

		// not_regex — reject if NOT matches
		if r.NotRegex != "" {
			if matched, _ := regexp.MatchString(r.NotRegex, val); !matched {
				return fmt.Errorf("[%s] %s", yp.ID, r.Message)
			}
		}
	}
	return nil
}

// ─── YAML webhook dispatch ──────────────────────────────────────────────────

// FanOutYAMLWebhooks triggers all YAML-declared webhooks matching the event.
// Settings interpolation is done at dispatch time via the plugin's DB settings.
func FanOutYAMLWebhooks(eventName string, payload map[string]any) {
	yamlMu.RLock()
	defer yamlMu.RUnlock()

	for _, yp := range yamlPlugins {
		for _, wh := range yp.Webhooks {
			for _, pattern := range wh.Events {
				if isEventMatch(pattern, eventName) {
					settings := loadYAMLPluginSettings(yp.ID)
					url := interpolateSettings(wh.URL, settings)
					headers := make(map[string]string, len(wh.Headers))
					for k, v := range wh.Headers {
						headers[k] = interpolateSettings(v, settings)
					}
					go fireWebhookHTTP(url, headers, eventName, payload)
					break // one match per webhook def is enough
				}
			}
		}
	}
}

// loadYAMLPluginSettings reads settings from the plugins table for a YAML plugin.
func loadYAMLPluginSettings(pluginID string) map[string]any {
	val, _ := g.DB().Ctx(context.Background()).
		Model("plugins").Where("id", pluginID).Value("settings")
	if val == nil || val.IsNil() {
		return nil
	}
	var m map[string]any
	_ = val.Scan(&m)
	return m
}

// interpolateSettings replaces {{settings.key}} placeholders.
func interpolateSettings(s string, settings map[string]any) string {
	if settings == nil || !strings.Contains(s, "{{") {
		return s
	}
	for k, v := range settings {
		placeholder := fmt.Sprintf("{{settings.%s}}", k)
		s = strings.ReplaceAll(s, placeholder, fmt.Sprintf("%v", v))
	}
	return s
}

// ─── YAML plugin queries ────────────────────────────────────────────────────

// GetYAMLPlugin returns a loaded YAML plugin by ID, or nil.
func GetYAMLPlugin(id string) *YAMLPlugin {
	yamlMu.RLock()
	defer yamlMu.RUnlock()
	return yamlPlugins[id]
}

// GetAllYAMLPlugins returns all loaded YAML plugins.
func GetAllYAMLPlugins() []*YAMLPlugin {
	yamlMu.RLock()
	defer yamlMu.RUnlock()
	result := make([]*YAMLPlugin, 0, len(yamlPlugins))
	for _, yp := range yamlPlugins {
		result = append(result, yp)
	}
	return result
}

// IsYAMLPlugin reports whether a plugin ID belongs to a loaded YAML plugin.
func IsYAMLPlugin(id string) bool {
	yamlMu.RLock()
	defer yamlMu.RUnlock()
	_, ok := yamlPlugins[id]
	return ok
}

// YAMLPluginToManifest converts a YAML plugin to the standard Manifest type
// so the admin UI can display it uniformly.
func YAMLPluginToManifest(yp *YAMLPlugin) Manifest {
	webhooks := make([]WebhookDef, len(yp.Webhooks))
	for i, w := range yp.Webhooks {
		webhooks[i] = WebhookDef{
			URL:     w.URL,
			Events:  w.Events,
			Headers: w.Headers,
		}
	}
	tl := TrustLevel(yp.TrustLevel)
	if tl == "" {
		tl = TrustLevelCommunity
	}
	var contributes *Contributes
	if len(yp.Navigation) > 0 {
		contributes = &Contributes{
			Navigation: yp.Navigation,
		}
	}
	return Manifest{
		Name:        yp.ID,
		Title:       yp.Title,
		Description: yp.Description,
		Version:     yp.Version,
		Author:      yp.Author,
		Icon:        yp.Icon,
		Priority:    yp.Priority,
		Settings:    yp.Settings,
		CSS:         yp.CSS,
		AdminJS:     yp.AdminJS,
		PublicJS:    yp.PublicJS,
		Webhooks:    webhooks,
		TrustLevel:  tl,
		Contributes: contributes,
	}
}
