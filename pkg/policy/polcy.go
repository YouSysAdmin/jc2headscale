package policy

import (
	"encoding/json"
	HsPolicy "github.com/juanfont/headscale/hscontrol/policy"
	"github.com/tailscale/hujson"
	"os"
)

// Policy extend Headscale policy
type Policy struct {
	HsPolicy.ACLPolicy
	JCGroupList []string `json:"jc_group_list"`
}

// ReadPolicyFromFile read Headscale policy from file
func (p *Policy) ReadPolicyFromFile(path string) error {
	policyData, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	ast, err := hujson.Parse(policyData)
	if err != nil {
		return err
	}
	ast.Standardize()
	data := ast.Pack()
	err = json.Unmarshal(data, &p)

	return err
}

// WritePolicyToFile write Headscale policy from file
func (p *Policy) WritePolicyToFile(path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(p)

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// AppendGroups append group to policy
func (p *Policy) AppendGroups(groups map[string][]string) {
	for g, u := range groups {
		p.Groups[g] = u
	}
}
