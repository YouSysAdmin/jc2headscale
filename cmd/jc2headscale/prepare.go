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

		// Get group names from policy file
		groups := hsPolicy.GetGroupNames()
		if len(groups) <= 0 {
			errorInfo := map[string]any{
				"Step":        "Get groups from policy file",
				"Policy file": inputPolicyFile,
				"Error":       "The policy doesn't contain a group list or the list is empty",
			}
			logger.Fatal("Prepare error:", logger.ArgsFromMap(errorInfo))
		}

		// Get groups and group members
		var jcGroupsInfo []*jc.Group
		for _, g := range groups {
			logger.Info(fmt.Sprintf("Get group: %s", g))

			// Get group info from Jumpcloud
			// If group doesn't find, returns nil
			group, err := client.GetGroupByName(g)
			if err != nil {
				errorInfo := map[string]any{
					"Step":      "Get group",
					"GroupName": g,
					"Error":     err.Error(),
				}
				logger.Fatal("Prepare error:", logger.ArgsFromMap(errorInfo))
			}

			// If a group is found in Jumpcloud, try to get a members
			if group != nil {
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
			} else {
				logger.Info(fmt.Sprintf("Group '%s' not foud in the Jumpcloud", g))
			}
		}

		hsGroups := map[string][]string{}
		for _, g := range jcGroupsInfo {
			var upg []string
			for _, u := range g.Users {
				upg = append(upg, u.Part)
			}

			// Add the prefix 'group' to a group name
			groupName := fmt.Sprintf("group:%s", g.Name)

			// If, in the policy, there are static users for a group,
			// then we add them, too
			upg = append(upg, hsPolicy.Groups[groupName]...)

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
