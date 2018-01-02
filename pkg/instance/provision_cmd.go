package instance

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/command"
	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/Azure/service-catalog-cli/pkg/parameters"
	"github.com/spf13/cobra"
)

type provisonCmd struct {
	*command.Context
	ns        string
	className string
	planName  string
	params    []string
	secrets   []string
}

// NewProvisionCmd builds a "svcat provision" command
func NewProvisionCmd(cxt *command.Context) *cobra.Command {
	provisionCmd := &provisonCmd{Context: cxt}
	cmd := &cobra.Command{
		Use:   "provision NAME --plan PLAN --class CLASS",
		Short: "Create a new instance of a service",
		Example: `
  svcat provision wordpress-mysql-instance --class azure-mysqldb --plan standard800 -p location=eastus -p sslEnforcement=disabled
  svcat provision wordpress-mysql-instance --class azure-mysqldb --plan standard800 -s mysecret[dbparams]
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return provisionCmd.run(args)
		},
	}
	cmd.Flags().StringVarP(&provisionCmd.ns, "namespace", "n", "default",
		"The namespace in which to create the instance")
	cmd.Flags().StringVar(&provisionCmd.className, "class", "",
		"The class name (Required)")
	cmd.MarkFlagRequired("class")
	cmd.Flags().StringVar(&provisionCmd.planName, "plan", "",
		"The plan name (Required)")
	cmd.MarkFlagRequired("plan")
	cmd.Flags().StringArrayVarP(&provisionCmd.params, "param", "p", nil,
		"Additional parameter to use when provisioning the service, format: NAME=VALUE")
	cmd.Flags().StringArrayVarP(&provisionCmd.secrets, "secret", "s", nil,
		"Additional parameter, whose value is stored in a secret, to use when provisioning the service, format: SECRET[KEY]")
	return cmd
}

func (c *provisonCmd) run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("name is required")
	}
	name := args[0]

	params, err := parameters.ParseVariableAssignments(c.params)
	if err != nil {
		return fmt.Errorf("invalid --param value (%s)", err)
	}

	secrets, err := parameters.ParseKeyMaps(c.secrets)
	if err != nil {
		return fmt.Errorf("invalid --secret value (%s)", err)
	}

	return c.provision(name, params, secrets)
}

func (c *provisonCmd) provision(name string, params map[string]string, secrets map[string]string) error {
	instance, err := provision(c.Client, c.ns, name, c.className, c.planName, params, secrets)
	if err != nil {
		return err
	}

	output.WriteInstanceDetails(c.Output, instance)

	return nil
}