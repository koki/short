package cmd

import (
	"bytes"
	"index/suffixarray"

	"github.com/spf13/cobra"

	"github.com/koki/short/generated"
	"github.com/koki/short/pager"
	"github.com/koki/short/util"
)

var (
	manCommand = &cobra.Command{
		Use:          "man",
		Short:        "Reference and Examples for resources and conversions",
		Long:         "Reference and Examples for koki <-> kubernetes conversions",
		RunE:         man,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
	}
)

func man(c *cobra.Command, args []string) error {
	resourceName := args[0]
	resourcePath := "../generated/" + resourceName + ".txt"
	data, err := generated.Asset(resourcePath)
	if err != nil {
		return util.UsageErrorf(c.CommandPath, "Unsupported resource name %s", resourceName)
	}

	index := suffixarray.New(data)

	buf := bytes.NewBuffer(data)
	p := pager.NewPager(buf, index)
	return p.Render()
}
