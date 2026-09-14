package cmd

import (
	"errors"
	"fmt"
	"github.com/previder/vault-cli/pkg/model"
	"github.com/spf13/cobra"
)

func init() {
	var secretCmd = &cobra.Command{
		Use:   "token",
		Short: "Token commands",
	}
	rootCmd.AddCommand(secretCmd)

	var cmdList = &cobra.Command{
		Use:    "list",
		Short:  "Get a list of tokens",
		Args:   cobra.NoArgs,
		RunE:   listToken,
		PreRun: setupClient,
	}
	secretCmd.AddCommand(cmdList)

	var cmdGet = &cobra.Command{
		Use:        "get [id]",
		Short:      "Get a token",
		Args:       cobra.ExactArgs(1),
		ArgAliases: []string{"id"},
		RunE:       getToken,
		PreRun:     setupClient,
	}
	secretCmd.AddCommand(cmdGet)

	var cmdCreate = &cobra.Command{
		Use:    "create",
		Short:  "Create a new token",
		RunE:   createToken,
		PreRun: setupClient,
	}
	secretCmd.AddCommand(cmdCreate)

	cmdCreate.Flags().StringP("description", "d", "", "Description of the token")
	cmdCreate.Flags().StringP("expires", "e", "", "Expire date of the token, date in following formats are allowed [2006-01-02 /  2006-01-02 15:04:05 / 2006-01-02T15:04:05 ]")
	cmdCreate.Flags().StringP("type", "r", model.TokenTypeReadWrite,
		fmt.Sprintf("Type of token [%s / %s / %s / %s]",
			model.TokenTypeReadOnly, model.TokenTypeReadWrite,
			model.TokenTypeEnvironmentAdmin, model.TokenTypeMasterAdmin))

	var cmdDelete = &cobra.Command{
		Use:        "delete [id]",
		Short:      "Delete a token",
		Args:       cobra.ExactArgs(1),
		ArgAliases: []string{"id"},
		RunE:       deleteToken,
		PreRun:     setupClient,
	}
	secretCmd.AddCommand(cmdDelete)
}

func listToken(cmd *cobra.Command, args []string) error {
	content, err := vaultClient.GetTokens()
	if err != nil {
		return err
	}
	if outputType == "pretty" {
		intContent := make([]interface{}, len(content))
		for idx, row := range content {
			intContent[idx] = row
		}
		if verbose {
			fmt.Printf("%+v\n", intContent)
		}
		printTable([]string{"Id", "Description", "TokenType"}, intContent)
	} else {
		printJson(content)
	}
	return nil
}

func getToken(cmd *cobra.Command, args []string) error {
	content, err := vaultClient.GetToken(args[0])
	if err != nil {
		return err
	}
	if outputType == "pretty" {
		fmt.Printf("%+v\n", content)
	} else {
		printJson(content)
	}
	return nil
}

func createToken(cmd *cobra.Command, args []string) error {
	create := model.TokenCreate{}
	description, err := cmd.Flags().GetString("description")
	if err != nil {
		return err
	}
	if description == "" {
		return errors.New("description cannot be empty")
	}
	create.Description = description

	expires, err := cmd.Flags().GetString("expire")
	if err == nil {
		create.ExpiresAt = expires
	}

	tokenType, err := cmd.Flags().GetString("type")
	if err != nil {
		return err
	}
	if tokenType == "" {
		return errors.New("type cannot be empty")
	}
	create.TokenType = tokenType

	createResponse, err := vaultClient.CreateToken(create)
	if err != nil {
		return err
	}

	if outputType == "pretty" {
		fmt.Printf("%+v\n", createResponse)
	} else {
		printJson(createResponse)
	}

	return nil
}

func deleteToken(cmd *cobra.Command, args []string) error {
	err := vaultClient.DeleteToken(args[0])
	if err != nil {
		return err
	}
	return nil
}
