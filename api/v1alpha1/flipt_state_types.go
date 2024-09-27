package v1alpha1

import "encoding/json"

type Document struct {
	Version   string     `yaml:"version,omitempty" json:"version,omitempty"`
	Namespace *Namespace `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	Flags     []*Flag    `yaml:"flags,omitempty" json:"flags,omitempty"`
	Segments  []*Segment `yaml:"segments,omitempty" json:"segments,omitempty"`
}

type Flag struct {
	Key         string                     `yaml:"key,omitempty" json:"key,omitempty"`
	Name        string                     `yaml:"name,omitempty" json:"name,omitempty"`
	Type        string                     `yaml:"type,omitempty" json:"type,omitempty"`
	Description string                     `yaml:"description,omitempty" json:"description,omitempty"`
	Enabled     bool                       `yaml:"enabled" json:"enabled"`
	Metadata    map[string]json.RawMessage `yaml:"metadata,omitempty" json:"metadata,omitempty"`
	Variants    []*Variant                 `yaml:"variants,omitempty" json:"variants,omitempty"`
	Rules       []*Rule                    `yaml:"rules,omitempty" json:"rules,omitempty"`
	Rollouts    []*Rollout                 `yaml:"rollouts,omitempty" json:"rollouts,omitempty"`
}

type Variant struct {
	Default     bool            `yaml:"default,omitempty" json:"default,omitempty"`
	Key         string          `yaml:"key,omitempty" json:"key,omitempty"`
	Name        string          `yaml:"name,omitempty" json:"name,omitempty"`
	Description string          `yaml:"description,omitempty" json:"description,omitempty"`
	Attachment  json.RawMessage `yaml:"attachment,omitempty" json:"attachment,omitempty"`
}

type Rule struct {
	Segment       *Segments       `yaml:"segment,omitempty" json:"segment,omitempty"`
	Rank          uint            `yaml:"rank,omitempty" json:"rank,omitempty"`
	Distributions []*Distribution `yaml:"distributions,omitempty" json:"distributions,omitempty"`
}

type Distribution struct {
	VariantKey string  `yaml:"variant,omitempty" json:"variant,omitempty"`
	Rollout    float32 `yaml:"rollout,omitempty" json:"rollout,omitempty"`
}

type Rollout struct {
	Description string         `yaml:"description,omitempty" json:"description,omitempty"`
	Segment     *SegmentRule   `yaml:"segment,omitempty" json:"segment,omitempty"`
	Threshold   *ThresholdRule `yaml:"threshold,omitempty" json:"threshold,omitempty"`
}

type SegmentRule struct {
	Key      string   `yaml:"key,omitempty" json:"key,omitempty"`
	Keys     []string `yaml:"keys,omitempty" json:"keys,omitempty"`
	Operator string   `yaml:"operator,omitempty" json:"operator,omitempty"`
	Value    bool     `yaml:"value,omitempty" json:"value,omitempty"`
}

type ThresholdRule struct {
	Percentage float32 `yaml:"percentage,omitempty" json:"percentage,omitempty"`
	Value      bool    `yaml:"value,omitempty" json:"value,omitempty"`
}

type Segment struct {
	Key         string        `yaml:"key,omitempty" json:"key,omitempty"`
	Name        string        `yaml:"name,omitempty" json:"name,omitempty"`
	Description string        `yaml:"description,omitempty" json:"description,omitempty"`
	Constraints []*Constraint `yaml:"constraints,omitempty" json:"constraints,omitempty"`
	MatchType   string        `yaml:"match_type,omitempty" json:"match_type,omitempty"`
}

type Constraint struct {
	Type        string `yaml:"type,omitempty" json:"type,omitempty"`
	Property    string `yaml:"property,omitempty" json:"property,omitempty"`
	Operator    string `yaml:"operator,omitempty" json:"operator,omitempty"`
	Value       string `yaml:"value,omitempty" json:"value,omitempty"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

type Segments struct {
	Keys            []string `yaml:"keys,omitempty" json:"keys,omitempty"`
	SegmentOperator string   `yaml:"operator,omitempty" json:"operator,omitempty"`
}

type Namespace struct {
	Key         string `yaml:"key,omitempty" json:"key,omitempty"`
	Name        string `yaml:"name,omitempty" json:"name,omitempty"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

func (n *Namespace) GetKey() string {
	if n == nil {
		return ""
	}
	return n.Key
}
