package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var sessionMergeCases = []struct {
	name     string
	target   Session
	override Session
	final    Session
}{
	// NB: the last two fields are not merged (Debug, Logger)
	{
		"no-override",
		Session{"./", "test-session", t_ptr(true), t_ptr("./.automux.hcl"), nil, false, nil},
		Session{},
		Session{"./", "test-session", t_ptr(true), t_ptr("./.automux.hcl"), nil, false, nil},
	},
	{
		"full-override",
		Session{"./", "test-session", t_ptr(true), t_ptr("./.automux.hcl"), nil, false, nil},
		Session{"../", "better-session", t_ptr(false), t_ptr("../.automux.hcl"), []Window{{}}, false, nil},
		Session{"../", "better-session", t_ptr(false), t_ptr("../.automux.hcl"), []Window{{}}, false, nil},
	},
}

func TestMergeSessions(t *testing.T) {
	for _, c := range sessionMergeCases {
		t.Run(c.name, func(t *testing.T) {
			merged := mergeSessions(c.target, c.override)
			require.Equal(t, c.final, merged)
		})
	}
}

var windowMergeCases = []struct {
	name     string
	target   []Window
	override []Window
	final    []Window
}{
	{
		"no-override",
		[]Window{{"win-1", t_ptr("nvim"), t_ptr(true), nil}},
		[]Window{{Title: "win-1"}},
		[]Window{{"win-1", t_ptr("nvim"), t_ptr(true), nil}},
	},
	{
		"no-override",
		[]Window{{"win-2", t_ptr("nvim"), t_ptr(true), nil}},
		[]Window{{"win-2", t_ptr("vim"), t_ptr(false), []Split{{}}}},
		[]Window{{"win-2", t_ptr("vim"), t_ptr(false), []Split{{}}}},
	},
	{
		"multi-windows",
		[]Window{{"win-1", t_ptr("nvim"), t_ptr(true), nil}},
		[]Window{{Title: "win-1"}, {"win-2", t_ptr("vim"), t_ptr(false), nil}},
		[]Window{{"win-1", t_ptr("nvim"), t_ptr(true), nil}, {"win-2", t_ptr("vim"), t_ptr(false), nil}},
	},
}

func TestMergeWindows(t *testing.T) {
	for _, c := range windowMergeCases {
		t.Run(c.name, func(t *testing.T) {
			merged := mergeWindows(c.target, c.override)
			require.Equal(t, c.final, merged)
		})
	}
}

var splitMergeCases = []struct {
	name     string
	target   []Split
	override []Split
	final    []Split
}{
	{
		"no-override",
		[]Split{{t_ptr(true), t_ptr("nvim"), t_ptr(10), t_ptr(false), nil}},
		[]Split{{}},
		[]Split{{t_ptr(true), t_ptr("nvim"), t_ptr(10), t_ptr(false), nil}},
	},
	{
		"full-override",
		[]Split{{t_ptr(true), t_ptr("nvim"), t_ptr(10), t_ptr(false), nil}},
		[]Split{{t_ptr(false), t_ptr("vim"), t_ptr(20), t_ptr(true), nil}},
		[]Split{{t_ptr(false), t_ptr("vim"), t_ptr(20), t_ptr(true), nil}},
	},
	{
		"multi-splits",
		[]Split{
			{t_ptr(true), t_ptr("nvim"), t_ptr(10), t_ptr(false), nil},
			{t_ptr(true), t_ptr("nvim"), t_ptr(10), t_ptr(false), nil},
		},
		[]Split{
			{t_ptr(false), t_ptr("vim"), t_ptr(20), t_ptr(true), nil},
			{},
		},
		[]Split{
			{t_ptr(false), t_ptr("vim"), t_ptr(20), t_ptr(true), nil},
			{t_ptr(true), t_ptr("nvim"), t_ptr(10), t_ptr(false), nil},
		},
	},
	{
		"extra-splits",
		[]Split{
			{t_ptr(true), t_ptr("nvim"), t_ptr(10), t_ptr(false), t_ptr("sub/")},
		},
		[]Split{
			{},
			{t_ptr(true), t_ptr("vim"), t_ptr(15), t_ptr(false), nil},
		},
		[]Split{
			{t_ptr(true), t_ptr("nvim"), t_ptr(10), t_ptr(false), t_ptr("sub/")},
			{t_ptr(true), t_ptr("vim"), t_ptr(15), t_ptr(false), nil},
		},
	},
}

func TestMergeSplits(t *testing.T) {
	for _, c := range splitMergeCases {
		t.Run(c.name, func(t *testing.T) {
			merged := mergeSplits(c.target, c.override)
			require.Equal(t, c.final, merged)
		})
	}
}

func t_ptr[T any](v T) *T {
	return &v
}
