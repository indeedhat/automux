package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var splitMergeCases = []struct {
	name     string
	target   []Split
	override []Split
	final    []Split
}{
	{
		"no-override",
		[]Split{{t_ptr(true), t_ptr("nvim"), t_ptr(10), t_ptr(false)}},
		[]Split{{}},
		[]Split{{t_ptr(true), t_ptr("nvim"), t_ptr(10), t_ptr(false)}},
	},
	{
		"full-override",
		[]Split{{t_ptr(true), t_ptr("nvim"), t_ptr(10), t_ptr(false)}},
		[]Split{{t_ptr(false), t_ptr("vim"), t_ptr(20), t_ptr(true)}},
		[]Split{{t_ptr(false), t_ptr("vim"), t_ptr(20), t_ptr(true)}},
	},
	{
		"multi-splits",
		[]Split{
			{t_ptr(true), t_ptr("nvim"), t_ptr(10), t_ptr(false)},
			{t_ptr(true), t_ptr("nvim"), t_ptr(10), t_ptr(false)},
		},
		[]Split{
			{t_ptr(false), t_ptr("vim"), t_ptr(20), t_ptr(true)},
			{},
		},
		[]Split{
			{t_ptr(false), t_ptr("vim"), t_ptr(20), t_ptr(true)},
			{t_ptr(true), t_ptr("nvim"), t_ptr(10), t_ptr(false)},
		},
	},
	{
		"extra-splits",
		[]Split{
			{t_ptr(true), t_ptr("nvim"), t_ptr(10), t_ptr(false)},
		},
		[]Split{
			{},
			{t_ptr(true), t_ptr("vim"), t_ptr(15), t_ptr(false)},
		},
		[]Split{
			{t_ptr(true), t_ptr("nvim"), t_ptr(10), t_ptr(false)},
			{t_ptr(true), t_ptr("vim"), t_ptr(15), t_ptr(false)},
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
