package cli

import (
	"context"
	"fmt"
	"github.com/Sifchain/sifnode/x/whitelist/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func GetQueryCmd() *cobra.Command {
	// Group dispensation queries under a subcommand
	dispensationQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	dispensationQueryCmd.AddCommand(
		GetCmdQueryDenoms(),
	)
	return dispensationQueryCmd
}

func GetCmdQueryDenoms() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "denoms",
		Short: "query the complete whitelist",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Denoms(context.Background(), &types.QueryDenomsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res.List)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
