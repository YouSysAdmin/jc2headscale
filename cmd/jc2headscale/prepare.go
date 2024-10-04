package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yousysadmin/jc2headscale/pkg/jc"
	"github.com/yousysadmin/jc2headscale/pkg/policy"
)

func init() {
	cliCmd.AddCommand(preparePolicy)
}

var preparePolicy = &cobra.Command{
	Use:     "prepare",
	Short:   "Prepare policy",
	Aliases: []string{"p"},
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("jc-api-key").Value.String() == "" {
			errorInfo := map[string]any{
				"Error": "jc-api-key is required",
			}
			logger.Fatal("Prepare error:", logger.ArgsFromMap(errorInfo))
		}

		hsPolicy := policy.Policy{}
		client := jc.NewClient(jcApiKey)

		// Read policy template
		logger.Info(fmt.Sprintf("Read policy template from: %s", inputPolicyFile))
		err := hsPolicy.ReadPolicyFromFile(inputPolicyFile)
		if err != nil {
			errorInfo := map[string]any{
				"Step":        "Read policy file",
				"Policy file": inputPolicyFile,
				"Error":       err.Error(),
			}
			logger.Fatal("Prepare error:", logger.ArgsFromMap(errorInfo))
		}

		if len(hsPolicy.JCGroupList) <= 0 {
			errorInfo := map[string]any{
				"Step":        "Get groups from policy file",
				"Policy file": inputPolicyFile,
				"Error":       "The policy doesn't contain a group list or the list is empty",
			}
			logger.Fatal("Prepare error:", logger.ArgsFromMap(errorInfo))
		}

		// Get groups and group members
		var jcGroupsInfo []jc.Group
		for _, g := range hsPolicy.JCGroupList {
			logger.Info(fmt.Sprintf("Get group: %s", g))

			group, err := client.GetGroupByName(g)
			if err != nil {
				errorInfo := map[string]any{
					"Step":      "Get group",
					"GroupName": g,
					"Error":     err.Error(),
				}
				logger.Fatal("Prepare error:", logger.ArgsFromMap(errorInfo))
			}

			users, err := client.GetGroupMembers(group.ID, stripEmailDomain)
			if err != nil {
				errorInfo := map[string]any{
					"Step":      "Get user list for group",
					"GroupName": g,
					"Error":     err.Error(),
				}
				logger.Fatal("Prepare error:", logger.ArgsFromMap(errorInfo))
			}

			group.Users = users
			jcGroupsInfo = append(jcGroupsInfo, group)

			logger.Info(fmt.Sprintf("Collect %d members for group: %s", len(users), g))
		}

		hsGroups := map[string][]string{}
		for _, g := range jcGroupsInfo {
			var upg []string
			for _, u := range g.Users {
				upg = append(upg, u.Part)
			}
			groupName := fmt.Sprintf("group:%s", g.Name)
			hsGroups[groupName] = upg
		}

		hsPolicy.AppendGroups(hsGroups)

		logger.Info(fmt.Sprintf("Write policy to: %s", outputPolicyFile))
		err = hsPolicy.WritePolicyToFile(outputPolicyFile)
		if err != nil {
			errorInfo := map[string]any{
				"Step":  "Write prepared policy file",
				"Error": err.Error(),
			}
			logger.Fatal("Prepare error:", logger.ArgsFromMap(errorInfo))
		}
	},
}
