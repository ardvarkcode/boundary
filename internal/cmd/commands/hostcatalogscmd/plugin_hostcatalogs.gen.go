// Code generated by "make cli"; DO NOT EDIT.
package hostcatalogscmd

import (
	"errors"
	"fmt"

	"github.com/hashicorp/boundary/api"
	"github.com/hashicorp/boundary/api/hostcatalogs"
	"github.com/hashicorp/boundary/internal/cmd/base"
	"github.com/hashicorp/boundary/internal/cmd/common"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

func initPluginFlags() {
	flagsOnce.Do(func() {
		extraFlags := extraPluginActionsFlagsMapFunc()
		for k, v := range extraFlags {
			flagsPluginMap[k] = append(flagsPluginMap[k], v...)
		}
	})
}

var (
	_ cli.Command             = (*PluginCommand)(nil)
	_ cli.CommandAutocomplete = (*PluginCommand)(nil)
)

type PluginCommand struct {
	*base.Command

	Func string

	plural string
}

func (c *PluginCommand) AutocompleteArgs() complete.Predictor {
	initPluginFlags()
	return complete.PredictAnything
}

func (c *PluginCommand) AutocompleteFlags() complete.Flags {
	initPluginFlags()
	return c.Flags().Completions()
}

func (c *PluginCommand) Synopsis() string {
	if extra := extraPluginSynopsisFunc(c); extra != "" {
		return extra
	}

	synopsisStr := "host catalog"

	synopsisStr = fmt.Sprintf("%s %s", "plugin-type", synopsisStr)

	return common.SynopsisFunc(c.Func, synopsisStr)
}

func (c *PluginCommand) Help() string {
	initPluginFlags()

	var helpStr string
	helpMap := common.HelpMap("host catalog")

	switch c.Func {
	default:

		helpStr = c.extraPluginHelpFunc(helpMap)
	}

	// Keep linter from complaining if we don't actually generate code using it
	_ = helpMap
	return helpStr
}

var flagsPluginMap = map[string][]string{

	"create": {"scope-id", "name", "description", "plugin-id", "plugin-name", "attributes", "attr", "string-attr", "bool-attr", "num-attr", "secrets", "secret", "string-secret", "bool-secret", "num-secret"},

	"update": {"id", "name", "description", "version", "attributes", "attr", "string-attr", "bool-attr", "num-attr", "secrets", "secret", "string-secret", "bool-secret", "num-secret"},
}

func (c *PluginCommand) Flags() *base.FlagSets {
	if len(flagsPluginMap[c.Func]) == 0 {
		return c.FlagSet(base.FlagSetNone)
	}

	set := c.FlagSet(base.FlagSetHTTP | base.FlagSetClient | base.FlagSetOutputFormat)
	f := set.NewFlagSet("Command Options")
	common.PopulateCommonFlags(c.Command, f, "plugin-type host catalog", flagsPluginMap, c.Func)

	f = set.NewFlagSet("Attribute Options")
	common.PopulateAttributeFlags(c.Command, f, flagsPluginMap, c.Func)

	f = set.NewFlagSet("Secrets Options")
	common.PopulateSecretFlags(c.Command, f, flagsPluginMap, c.Func)

	extraPluginFlagsFunc(c, set, f)

	return set
}

func (c *PluginCommand) Run(args []string) int {
	initPluginFlags()

	switch c.Func {
	case "":
		return cli.RunResultHelp
	}

	c.plural = "plugin-type host catalog"
	switch c.Func {
	case "list":
		c.plural = "plugin-type host catalogs"
	}

	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.PrintCliError(err)
		return base.CommandUserError
	}

	if strutil.StrListContains(flagsPluginMap[c.Func], "id") && c.FlagId == "" {
		c.PrintCliError(errors.New("ID is required but not passed in via -id"))
		return base.CommandUserError
	}

	var opts []hostcatalogs.Option

	if strutil.StrListContains(flagsPluginMap[c.Func], "scope-id") {
		switch c.Func {
		case "create":
			if c.FlagScopeId == "" {
				c.PrintCliError(errors.New("Scope ID must be passed in via -scope-id or BOUNDARY_SCOPE_ID"))
				return base.CommandUserError
			}
		}
	}

	client, err := c.Client()
	if err != nil {
		c.PrintCliError(fmt.Errorf("Error creating API client: %s", err.Error()))
		return base.CommandCliError
	}
	hostcatalogsClient := hostcatalogs.NewClient(client)

	switch c.FlagName {
	case "":
	case "null":
		opts = append(opts, hostcatalogs.DefaultName())
	default:
		opts = append(opts, hostcatalogs.WithName(c.FlagName))
	}

	switch c.FlagDescription {
	case "":
	case "null":
		opts = append(opts, hostcatalogs.DefaultDescription())
	default:
		opts = append(opts, hostcatalogs.WithDescription(c.FlagDescription))
	}

	switch c.FlagRecursive {
	case true:
		opts = append(opts, hostcatalogs.WithRecursive(true))
	}

	if c.FlagFilter != "" {
		opts = append(opts, hostcatalogs.WithFilter(c.FlagFilter))
	}

	switch c.FlagPluginId {
	case "":
	default:
		opts = append(opts, hostcatalogs.WithPluginId(c.FlagPluginId))
	}
	switch c.FlagPluginName {
	case "":
	default:
		opts = append(opts, hostcatalogs.WithPluginName(c.FlagPluginName))
	}

	var version uint32

	switch c.Func {
	case "update":
		switch c.FlagVersion {
		case 0:
			opts = append(opts, hostcatalogs.WithAutomaticVersioning(true))
		default:
			version = uint32(c.FlagVersion)
		}
	}

	if err := common.HandleAttributeFlags(
		c.Command,
		"attr",
		c.FlagAttributes,
		c.FlagAttrs,
		func() {
			opts = append(opts, hostcatalogs.DefaultAttributes())
		},
		func(in map[string]interface{}) {
			opts = append(opts, hostcatalogs.WithAttributes(in))
		}); err != nil {
		c.PrintCliError(fmt.Errorf("Error evaluating attribute flags to: %s", err.Error()))
		return base.CommandCliError
	}

	if err := common.HandleAttributeFlags(
		c.Command,
		"secret",
		c.FlagSecrets,
		c.FlagScrts,
		func() {
			opts = append(opts, hostcatalogs.DefaultSecrets())
		},
		func(in map[string]interface{}) {
			opts = append(opts, hostcatalogs.WithSecrets(in))
		}); err != nil {
		c.PrintCliError(fmt.Errorf("Error evaluating secret flags to: %s", err.Error()))
		return base.CommandCliError
	}

	if ok := extraPluginFlagsHandlingFunc(c, f, &opts); !ok {
		return base.CommandUserError
	}

	var result api.GenericResult

	switch c.Func {

	case "create":
		result, err = hostcatalogsClient.Create(c.Context, "plugin", c.FlagScopeId, opts...)

	case "update":
		result, err = hostcatalogsClient.Update(c.Context, c.FlagId, version, opts...)

	}

	result, err = executeExtraPluginActions(c, result, err, hostcatalogsClient, version, opts)

	if err != nil {
		if apiErr := api.AsServerError(err); apiErr != nil {
			var opts []base.Option

			c.PrintApiError(apiErr, fmt.Sprintf("Error from controller when performing %s on %s", c.Func, c.plural), opts...)
			return base.CommandApiError
		}
		c.PrintCliError(fmt.Errorf("Error trying to %s %s: %s", c.Func, c.plural, err.Error()))
		return base.CommandCliError
	}

	output, err := printCustomPluginActionOutput(c)
	if err != nil {
		c.PrintCliError(err)
		return base.CommandUserError
	}
	if output {
		return base.CommandSuccess
	}

	switch c.Func {
	}

	switch base.Format(c.UI) {
	case "table":
		c.UI.Output(printItemTable(result))

	case "json":
		if ok := c.PrintJsonItem(result); !ok {
			return base.CommandCliError
		}
	}

	return base.CommandSuccess
}

var (
	extraPluginActionsFlagsMapFunc = func() map[string][]string { return nil }
	extraPluginSynopsisFunc        = func(*PluginCommand) string { return "" }
	extraPluginFlagsFunc           = func(*PluginCommand, *base.FlagSets, *base.FlagSet) {}
	extraPluginFlagsHandlingFunc   = func(*PluginCommand, *base.FlagSets, *[]hostcatalogs.Option) bool { return true }
	executeExtraPluginActions      = func(_ *PluginCommand, inResult api.GenericResult, inErr error, _ *hostcatalogs.Client, _ uint32, _ []hostcatalogs.Option) (api.GenericResult, error) {
		return inResult, inErr
	}
	printCustomPluginActionOutput = func(*PluginCommand) (bool, error) { return false, nil }
)