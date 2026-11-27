package casbin

import (
	"context"
	"fmt"
	"sort"

	"base/ent"
	"base/ent/syscasbinrule"
	"base/ent/sysmenu"
	"base/ent/sysrolemenu"
	"base/ent/sysuser"
	"base/ent/sysuserrole"
)

func SyncRolePolicies(ctx context.Context, client *ent.Client, roleID int) error {
	roleSubject := fmt.Sprintf("role:%d", roleID)

	if _, err := client.SysCasbinRule.Delete().
		Where(
			syscasbinrule.PtypeEQ("p"),
			syscasbinrule.V0EQ(roleSubject),
		).
		Exec(ctx); err != nil {
		return err
	}

	menuIDs, err := client.SysRoleMenu.Query().
		Where(sysrolemenu.RoleID(roleID)).
		Select(sysrolemenu.FieldMenuID).
		Ints(ctx)
	if err != nil {
		return err
	}
	if len(menuIDs) == 0 {
		return nil
	}

	menuList, err := client.SysMenu.Query().
		Where(sysmenu.IDIn(menuIDs...)).
		All(ctx)
	if err != nil {
		return err
	}

	codeMap := map[string]struct{}{}
	for _, menu := range menuList {
		if menu.PermissionCode != "" {
			codeMap[menu.PermissionCode] = struct{}{}
		}
	}

	codes := make([]string, 0, len(codeMap))
	for code := range codeMap {
		codes = append(codes, code)
	}
	sort.Strings(codes)

	for _, code := range codes {
		if err := client.SysCasbinRule.Create().
			SetPtype("p").
			SetV0(roleSubject).
			SetV1(code).
			SetV2("access").
			Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}

func SyncUserRoles(ctx context.Context, client *ent.Client, userID int) error {
	user, err := client.SysUser.Query().
		Where(sysuser.IDEQ(userID)).
		Only(ctx)
	if err != nil {
		return err
	}

	userSubject := "user:" + user.Username
	if _, err := client.SysCasbinRule.Delete().
		Where(
			syscasbinrule.PtypeEQ("g"),
			syscasbinrule.V0EQ(userSubject),
		).
		Exec(ctx); err != nil {
		return err
	}

	roleIDs, err := client.SysUserRole.Query().
		Where(sysuserrole.UserID(userID)).
		Select(sysuserrole.FieldRoleID).
		Ints(ctx)
	if err != nil {
		return err
	}

	for _, roleID := range roleIDs {
		if err := client.SysCasbinRule.Create().
			SetPtype("g").
			SetV0(userSubject).
			SetV1(fmt.Sprintf("role:%d", roleID)).
			Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}
