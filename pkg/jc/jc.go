package jc

import (
	"context"
	"fmt"
	jcapiv1 "github.com/TheJumpCloud/jcapi-go/v1"
	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"strings"
)

// JCClient Jumpcloud client
type JCClient struct {
	V1          *jcapiv1.APIClient
	V1Auth      context.Context
	V2          *jcapiv2.APIClient
	V2Auth      context.Context
	ContentType string
}

// Group Jumpcloud group info
type Group struct {
	ID    string
	Name  string
	Users []User
}

// User Jumpcloud user info
type User struct {
	ID    string
	Email string
	Part  string
}

// NewClient create new Jumpcloud client
func NewClient(apiKey string) JCClient {
	c := JCClient{}
	c.V1 = jcapiv1.NewAPIClient(jcapiv1.NewConfiguration())
	c.V1Auth = context.WithValue(context.TODO(), jcapiv1.ContextAPIKey, jcapiv1.APIKey{
		Key: apiKey,
	})

	c.V2 = jcapiv2.NewAPIClient(jcapiv2.NewConfiguration())
	c.V2Auth = context.WithValue(context.TODO(), jcapiv2.ContextAPIKey, jcapiv2.APIKey{
		Key: apiKey,
	})

	c.ContentType = "application/json"
	return c
}

// GetGroupByName Get Jumpcloud group by name
func (c JCClient) GetGroupByName(grounName string) (Group, error) {

	filter := map[string]interface{}{
		"filter": []string{fmt.Sprintf("name:eq:%s", grounName)},
		"limit":  int32(100),
	}

	group, _, err := c.V2.UserGroupsApi.GroupsUserList(c.V2Auth, c.ContentType, c.ContentType, filter)
	if err != nil {
		return Group{}, err
	}

	if len(group) != 0 {
		return Group{
			ID:   group[0].Id,
			Name: group[0].Name,
		}, nil
	}

	return Group{}, fmt.Errorf("group '%s' not found", grounName)
}

// GetGroupMembers Get Jumpcloud group members
func (c JCClient) GetGroupMembers(groupId string, stripEmailDomain bool) ([]User, error) {

	var users []User

	options := map[string]interface{}{
		"limit": int32(100),
	}

	groupUsers, _, err := c.V2.UserGroupsApi.GraphUserGroupMembership(c.V2Auth, groupId, c.ContentType, c.ContentType, options)
	if err != nil {
		return nil, err
	}

	for _, u := range groupUsers {
		user, err := c.GetUserInfo(u.Id, stripEmailDomain)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUserInfo get Jumpcloud user info
func (c JCClient) GetUserInfo(userId string, stripEmailDomain bool) (User, error) {

	options := map[string]interface{}{
		"limit": int32(100),
	}

	user, _, err := c.V1.SystemusersApi.SystemusersGet(c.V1Auth, userId, c.ContentType, c.ContentType, options)
	if err != nil {
		return User{}, err
	}

	var userName string
	if stripEmailDomain {
		userName = strings.Split(user.Email, "@")[0]
	} else {
		userName = user.Email
	}

	userInfo := User{
		ID:    user.Id,
		Email: user.Email,
		Part:  userName,
	}

	return userInfo, nil
}
