package media

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/nuxtblog/nuxtblog/internal/consts"
	"github.com/nuxtblog/nuxtblog/internal/dao"
)

const (
	optKeyExtensionGroups = "media_extension_groups"
	optKeyFormatPolicies  = "media_format_policies"
)

// ── Extension Groups ──────────────────────────────────────────────────────────

func getExtensionGroups(ctx context.Context) []consts.ExtensionGroup {
	var groups []consts.ExtensionGroup
	if err := getOptionJSON(ctx, optKeyExtensionGroups, &groups); err != nil || len(groups) == 0 {
		return consts.DefaultExtensionGroups
	}
	return groups
}

func saveExtensionGroups(ctx context.Context, groups []consts.ExtensionGroup) error {
	data, err := json.Marshal(groups)
	if err != nil {
		return err
	}
	return upsertOption(ctx, optKeyExtensionGroups, string(data))
}

// ── Format Policies ───────────────────────────────────────────────────────────

func getFormatPolicies(ctx context.Context) []consts.FormatPolicy {
	var policies []consts.FormatPolicy
	if err := getOptionJSON(ctx, optKeyFormatPolicies, &policies); err != nil || len(policies) == 0 {
		return consts.DefaultFormatPolicies
	}
	return policies
}

func saveFormatPolicies(ctx context.Context, policies []consts.FormatPolicy) error {
	data, err := json.Marshal(policies)
	if err != nil {
		return err
	}
	return upsertOption(ctx, optKeyFormatPolicies, string(data))
}

func getFormatPolicyByName(ctx context.Context, name string) *consts.FormatPolicy {
	if name == "" {
		name = "default"
	}
	for _, p := range getFormatPolicies(ctx) {
		if p.Name == name {
			return &p
		}
	}
	return nil
}

// ── Upload validation ─────────────────────────────────────────────────────────

// resolveFormatPolicy gets the FormatPolicy for a given category.
func resolveFormatPolicy(ctx context.Context, category string) *consts.FormatPolicy {
	catDef := GetCategoryDef(ctx, category)
	policyName := ""
	if catDef != nil {
		policyName = catDef.FormatPolicy
	}
	return getFormatPolicyByName(ctx, policyName)
}

// validateFileExtension checks if a file extension is allowed by the given policy
// and whether the file size is within limits.
func validateFileExtension(ctx context.Context, policy *consts.FormatPolicy, ext string, fileSize int64) error {
	if policy == nil {
		// No policy = use default
		policy = getFormatPolicyByName(ctx, "default")
	}
	if policy == nil {
		return nil // no policy at all, allow everything
	}

	// Normalise extension (strip leading dot, lowercase)
	ext = strings.TrimPrefix(strings.ToLower(ext), ".")

	groups := getExtensionGroups(ctx)
	groupMap := make(map[string]*consts.ExtensionGroup, len(groups))
	for i := range groups {
		groupMap[groups[i].Name] = &groups[i]
	}

	for _, groupName := range policy.Groups {
		grp, ok := groupMap[groupName]
		if !ok {
			continue
		}
		for _, allowed := range grp.Extensions {
			if ext == allowed {
				// Extension matches — check size
				maxBytes := int64(grp.MaxSizeMB * 1024 * 1024)
				if maxBytes > 0 && fileSize > maxBytes {
					return fmt.Errorf("file size %.1f MB exceeds limit %.0f MB for %s",
						float64(fileSize)/(1024*1024), grp.MaxSizeMB, grp.LabelEn)
				}
				return nil // allowed
			}
		}
	}

	return fmt.Errorf("file extension .%s is not allowed by the current format policy", ext)
}

// ── IMedia method implementations for format policies ─────────────────────────

func (s *sMedia) GetExtensionGroups(ctx context.Context) ([]consts.ExtensionGroup, error) {
	return getExtensionGroups(ctx), nil
}

func (s *sMedia) SaveExtensionGroups(ctx context.Context, groups []consts.ExtensionGroup) error {
	return saveExtensionGroups(ctx, groups)
}

func (s *sMedia) GetFormatPolicies(ctx context.Context) ([]consts.FormatPolicy, error) {
	return getFormatPolicies(ctx), nil
}

func (s *sMedia) CreateFormatPolicy(ctx context.Context, policy consts.FormatPolicy) error {
	policies := getFormatPolicies(ctx)
	for _, p := range policies {
		if p.Name == policy.Name {
			return errors.New("format policy already exists")
		}
	}
	policies = append(policies, policy)
	return saveFormatPolicies(ctx, policies)
}

func (s *sMedia) UpdateFormatPolicy(ctx context.Context, name string, policy consts.FormatPolicy) error {
	policies := getFormatPolicies(ctx)
	for i := range policies {
		if policies[i].Name == name {
			if policies[i].IsSystem {
				// System policies: only allow updating groups
				policies[i].Groups = policy.Groups
			} else {
				policy.Name = name // ensure name stays the same
				policies[i] = policy
			}
			return saveFormatPolicies(ctx, policies)
		}
	}
	return errors.New("format policy not found")
}

func (s *sMedia) DeleteFormatPolicy(ctx context.Context, name string) error {
	policies := getFormatPolicies(ctx)
	for i, p := range policies {
		if p.Name == name {
			if p.IsSystem {
				return errors.New("cannot delete system format policy")
			}
			policies = append(policies[:i], policies[i+1:]...)
			return saveFormatPolicies(ctx, policies)
		}
	}
	return errors.New("format policy not found")
}

// ── Option helper ─────────────────────────────────────────────────────────────

func upsertOption(ctx context.Context, key, value string) error {
	cnt, _ := dao.Options.Ctx(ctx).Where("key", key).Count()
	var err error
	if cnt == 0 {
		_, err = dao.Options.Ctx(ctx).Data(g.Map{
			"key": key, "value": value, "autoload": 1,
		}).Insert()
	} else {
		_, err = dao.Options.Ctx(ctx).Where("key", key).Data(g.Map{"value": value}).Update()
	}
	return err
}
