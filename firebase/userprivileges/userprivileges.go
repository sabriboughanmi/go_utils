package userprivileges

import (
	"../../utils"
	"context"
	"firebase.google.com/go/auth"
)

//UserPrivilegeClient is a client extending the firebase auth package to manage users Privileges
type UserPrivilegeClient struct {
	authClient *auth.Client
	ctx        context.Context
}

func GetUserPrivilegeClient (authClient *auth.Client, ctx context.Context) UserPrivilegeClient {
	return  UserPrivilegeClient{
		authClient,
		ctx,
	}
}


//CustomClaimsPrivileges is the key to access Auth.UserRecord Privileges.
//Note! its forbidden to add "privileges" key to a firebase.auth User when extending userprivileges Package
var CustomClaimsPrivileges = "privileges"

//Privilege is type (uint16) used to define users Privileges
type Privilege uint16


//HasPrivileges checks if a user has specified privileges
func (client *UserPrivilegeClient)HasPrivilegesByID(userID string, privileges ...Privilege) (bool,error) {

	userRecord, err := client.authClient.GetUser(client.ctx, userID)
	if err != nil {
		return false,err
	}

	var userPrivileges []Privilege

	if value, exist := userRecord.CustomClaims[CustomClaimsPrivileges]; exist {
		if err =  utils.InterfaceToType(value, &userPrivileges); err != nil {
			return false,err
		}

	}else{ //userPrivileges are empty
		return false,nil
	}

	if len(difference(privileges, userPrivileges)) > 0{
		return false,nil
	}

	return true,nil
}

//HasPrivileges checks if a user has specified privileges
func (client *UserPrivilegeClient)HasPrivileges(userRecord *auth.UserRecord, privileges ...Privilege) bool {

	var userPrivileges []Privilege

	if value, exist := userRecord.CustomClaims[CustomClaimsPrivileges]; exist {
		if err := utils.InterfaceToType(value, &userPrivileges); err != nil {
			return false
		}

	}else{ //userPrivileges are empty
		return false
	}

	if len(difference(privileges, userPrivileges)) > 0{
		return false
	}

	return true
}




//SetUserPrivilegesByID gives a user some privileges
func (client *UserPrivilegeClient)SetUserPrivilegesByID(userID string, privileges ...Privilege) error {

	userRecord, err := client.authClient.GetUser(client.ctx, userID)
	if err != nil {
		return err
	}

	var userPrivileges []Privilege
	if value, exist := userRecord.CustomClaims[CustomClaimsPrivileges]; exist {
		if err = utils.InterfaceToType(value, &userPrivileges); err != nil {
			return err
		}
	}

	privilegesToAdd :=  difference(privileges, userPrivileges)
	userPrivileges = append(userPrivileges, privilegesToAdd...)

	userRecord.CustomClaims[CustomClaimsPrivileges] = userPrivileges

	if err = client.authClient.SetCustomUserClaims(client.ctx, userID, userRecord.CustomClaims); err != nil {
		return err
	}
	return nil
}

//SetUserPrivileges gives a user some privileges
func (client *UserPrivilegeClient)SetUserPrivileges(userRecord *auth.UserRecord, privileges ...Privilege) error {

	var userPrivileges []Privilege
	if value, exist := userRecord.CustomClaims[CustomClaimsPrivileges]; exist {
		if err := utils.InterfaceToType(value, &userPrivileges); err != nil {
			return err
		}
	}

	privilegesToAdd :=  difference(privileges, userPrivileges)
	userPrivileges = append(userPrivileges, privilegesToAdd...)

	userRecord.CustomClaims[CustomClaimsPrivileges] = userPrivileges

	if err := client.authClient.SetCustomUserClaims(client.ctx, userRecord.UID, userRecord.CustomClaims); err != nil {
		return err
	}
	return nil
}




//GetUserPrivilegesByID returns user Granted privileges
func (client *UserPrivilegeClient)GetUserPrivilegesByID(userID string) ( []Privilege,error) {

	userRecord, err := client.authClient.GetUser(client.ctx, userID)
	if err != nil {
		return nil,err
	}

	var userPrivileges []Privilege
	if value, exist := userRecord.CustomClaims[CustomClaimsPrivileges]; exist {
		if err = utils.InterfaceToType(value, &userPrivileges); err != nil {
			return nil,err
		}
	}
	return userPrivileges,nil
}

//GetUserPrivileges returns user Granted privileges
func (client *UserPrivilegeClient)GetUserPrivileges(userRecord *auth.UserRecord) []Privilege {
	var userPrivileges []Privilege
	if value, exist := userRecord.CustomClaims[CustomClaimsPrivileges]; exist {
		if err := utils.InterfaceToType(value, &userPrivileges); err != nil {
			return nil
		}
	}
	return userPrivileges
}


