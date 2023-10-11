package serve

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func Command(parents ...string) *cobra.Command {
	var command options
	cmd := &cobra.Command{
		Use:          "serve",
		Short:        "serve",
		Long:         "Serve API",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE:         command.Run,
	}
	// server flags
	cmd.Flags().StringVar(&command.serverFlags.host, "server-host", "0.0.0.0", "server host")
	cmd.Flags().IntVar(&command.serverFlags.port, "server-port", 8080, "server port")
	// gin flags
	cmd.Flags().StringVar(&command.ginFlags.mode, "gin-mode", gin.ReleaseMode, "gin run mode")
	cmd.Flags().BoolVar(&command.ginFlags.log, "gin-log", true, "enable gin logger")
	cmd.Flags().BoolVar(&command.ginFlags.cors, "gin-cors", true, "enable gin cors")
	cmd.Flags().IntVar(&command.ginFlags.maxBodySize, "gin-max-body-size", 2*1024*1024, "gin max body size")
	// // ui flags
	// res.Flags().StringVar(&command.uiFlags.sponsor, "ui-sponsor", "", "sponsor text")
	// // engine flags
	// res.Flags().StringSliceVar(&command.engineFlags.builtInCrds, "engine-builtin-crds", nil, "list of enabled builtin custom resource definitions")
	// res.Flags().StringSliceVar(&command.engineFlags.localCrds, "engine-local-crds", nil, "list of folders containing custom resource definitions")
	// // cluster flags
	// res.Flags().BoolVar(&command.clusterFlags.cluster, "cluster", false, "enable cluster connected mode")
	// clientcmd.BindOverrideFlags(&command.clusterFlags.kubeConfigOverrides, res.Flags(), clientcmd.RecommendedConfigOverrideFlags("kube-"))
	return cmd
}
