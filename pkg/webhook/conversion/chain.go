package conversion

import (
	"strings"

	"github.com/flant/shell-operator/pkg/utils/string_helper"
)

type ChainStorage map[string]*Chain

func NewChainStorage() ChainStorage {
	return make(map[string]*Chain)
}

type Chain struct {
	// Index: ruleID ("srcVer->desiredVer") to a sequence of "primitive" ruleIDs.
	PathsCache map[string][]string
	// from -> to cache
	FromToCache map[string]map[string]bool
}

// Access

func (cs ChainStorage) Get(crdName string) *Chain {
	if _, ok := cs[crdName]; !ok {
		cs[crdName] = &Chain{
			PathsCache:  make(map[string][]string),
			FromToCache: make(map[string]map[string]bool),
		}
	}
	return cs[crdName]
}

func (c *Chain) Put(rule ConversionRule) {
	// paths "from->to"
	ruleID := rule.String()
	c.PathsCache[ruleID] = []string{ruleID}
	// from -> to
	if _, ok := c.FromToCache[rule.FromVersion]; !ok {
		c.FromToCache[rule.FromVersion] = make(map[string]bool)
	}
	c.FromToCache[rule.FromVersion][rule.ToVersion] = true
}

// Calculations

func (cs ChainStorage) FindConversionChain(crdName string, rule ConversionRule) []string {
	chain, ok := cs[crdName]
	if !ok {
		return nil
	}

	// Return if there is no path to a desired version, trimmed or full.
	if !chain.HasTargetVersion(rule.ToVersion) {
		return nil
	}

	for {
		p := chain.SearchPathForRule(rule)
		if len(p) > 0 {
			return p
		}

		// Fill cache with more paths.
		newPaths := map[string][]string{}

		// Try only ids that starts from a source version.
		for _, ruleIDToCheck := range chain.IDsByFromVersion(rule) {
			ruleToCheck := RuleFromString(ruleIDToCheck)

			if ruleToCheck.ShortToVersion() == rule.ShortFromVersion() {
				// Ignore loops.
				continue
			}

			// toVersion in ruleIDToCheck is a new start. Get toVersions available starting from it.
			for _, nextRule := range chain.NextRules(ruleToCheck.ToVersion) {
				newRule := ConversionRule{
					FromVersion: rule.FromVersion,
					ToVersion:   nextRule.ToVersion,
				}

				if newRule.ShortToVersion() == rule.ShortFromVersion() {
					// Ignore loops.
					continue
				}

				newPath := append(chain.PathsCache[ruleIDToCheck], nextRule.String())

				// This path is already discovered.
				p := chain.SearchPathForRule(newRule)
				if len(p) != 0 {
					continue
				}

				newPaths[newRule.String()] = newPath
			}
		}

		// break if no new paths are discovered.
		if len(newPaths) == 0 {
			break
		}

		// Put new paths in cache.
		for id, path := range newPaths {
			chain.PathsCache[id] = path
		}
	}

	return nil
}

func (c Chain) SearchPathForRule(rule ConversionRule) []string {
	IDs := []string{}
	for k := range c.PathsCache {
		r := RuleFromString(k)
		// Return is full equal is found.
		if k == rule.String() {
			return c.PathsCache[k]
		}
		if VersionsMatched(rule.ToVersion, r.ToVersion) && VersionsMatched(rule.FromVersion, r.FromVersion) {
			IDs = append(IDs, k)
		}
	}

	// Oops. No similar paths.
	if len(IDs) == 0 {
		return nil
	}
	// Return if only one ID is found.
	if len(IDs) == 1 {
		return c.PathsCache[IDs[0]]
	}

	// Try to find a more stricter match. Prefer match of full toVersions.
	// IDs len should not be more than 3 items from these variants:
	// 1. group1/v1->v2
	// 2. group1/v1->group2/v2
	// 3. v1->group2/v2
	// 4. v1->v2
	// There should be a full equal variant and it is already returned earlier.
	//
	fromMatches := []string{}
	toMatches := []string{}
	idxFrom := strings.IndexRune(rule.FromVersion, '/')
	idxTo := strings.IndexRune(rule.ToVersion, '/')
	for _, k := range IDs {
		r := RuleFromString(k)
		cc := 0
		if idxFrom >= 0 && r.FromVersion == rule.FromVersion {
			fromMatches = append(fromMatches, k)
			cc++
		}
		if idxTo >= 0 && r.ToVersion == rule.ToVersion {
			toMatches = append(toMatches, k)
			cc++
		}
		if cc == 2 {
			// Full equal is found.
			return c.PathsCache[k]
		}
	}
	if len(toMatches) > 0 {
		return c.PathsCache[toMatches[0]]
	}
	if len(fromMatches) > 0 {
		return c.PathsCache[fromMatches[0]]
	}
	return c.PathsCache[IDs[0]]
}

func (c Chain) IDsByFromVersion(rule ConversionRule) []string {
	IDs := []string{}
	for k := range c.PathsCache {
		idxFrom := strings.Index(k, rule.ShortFromVersion())
		if idxFrom == -1 {
			continue
		}
		idxSep := strings.Index(k, "->")
		if idxSep == -1 {
			continue
		}
		if idxFrom < idxSep {
			IDs = append(IDs, k)
		}
	}
	return IDs
}

// AvailableToVersions finds all toVersions by an input fromVer in FromToCache maps.
func (c Chain) NextRules(fromVer string) []ConversionRule {
	rules := []ConversionRule{}
	shortVer := string_helper.TrimGroup(fromVer)
	for k := range c.FromToCache {
		//
		if k == fromVer {
			for toVer := range c.FromToCache[k] {
				rules = append(rules, ConversionRule{
					FromVersion: k,
					ToVersion:   toVer,
				})
			}
			continue
		}

		idxFrom := strings.Index(k, shortVer)
		if idxFrom == -1 {
			continue
		}
		for toVer := range c.FromToCache[k] {
			rules = append(rules, ConversionRule{
				FromVersion: k,
				ToVersion:   toVer,
			})
		}
	}

	return rules
}

// HasTargetVersion returns true if there is a final version that matches input version.
func (c Chain) HasTargetVersion(finalVer string) bool {
	for fromVer := range c.FromToCache {
		for toVer := range c.FromToCache[fromVer] {
			if VersionsMatched(finalVer, toVer) {
				return true
			}
		}
	}
	return false
}

// VersionsMatched when:
// - v0 equals to v1
// - v0 is short and v1 is full and short v1 is equals to v0
// - v0 is full and v1 is short and short v0 is equals to v
func VersionsMatched(v0, v1 string) bool {
	if v0 == v1 {
		return true
	}
	idx0 := strings.IndexRune(v0, '/')
	idx1 := strings.IndexRune(v1, '/')
	if idx0 == -1 && idx1 >= 0 {
		return v0 == v1[idx1+1:]
	}
	if idx0 >= 0 && idx1 == -1 {
		return v0[idx0+1:] == v1
	}
	return false
}
