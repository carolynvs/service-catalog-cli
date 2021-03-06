package broker

import (
	"fmt"

	"github.com/Azure/service-catalog-cli/pkg/output"
	"github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type getCmd struct {
	cl *clientset.Clientset
}

// NewGetCmd builds a "svc-cat get brokers" command
func NewGetCmd(cl *clientset.Clientset) *cobra.Command {
	getCmd := getCmd{cl: cl}
	cmd := &cobra.Command{
		Use:     "brokers [name]",
		Aliases: []string{"broker", "brk"},
		Short:   "List brokers, optionally filtered by name",
		Example: `
  svc-cat get brokers
  svc-cat get broker asb
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmd.run(args)
		},
	}

	return cmd
}

func (c *getCmd) run(args []string) error {
	if len(args) == 0 {
		return c.getAll()
	} else {
		name := args[0]
		return c.get(name)
	}
}

func (c *getCmd) getAll() error {
	brokers, err := c.cl.ServicecatalogV1beta1().ClusterServiceBrokers().List(v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("Error listing brokers (%s)", err)
	}

	output.WriteBrokerList(brokers.Items...)
	return nil
}

func (c *getCmd) get(name string) error {
	broker, err := retrieveByName(c.cl, name)
	if err != nil {
		return err
	}

	output.WriteBrokerList(*broker)
	return nil
}
