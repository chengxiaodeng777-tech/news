package auth

import (
	"context"
	"sort"
	"strings"
	"time"

	"base/ent"
	"base/ent/sysmenu"
	"base/ent/sysrolemenu"
	"base/ent/sysuser"
	"base/ent/sysuserrole"
	authmodel "base/internal/model/auth"
	"base/internal/pkg/apperror"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	client    *ent.Client
	jwtSecret string
}

type Claims struct {
	UserID   int      `json:"userId"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

func New(client *ent.Client, jwtSecret string) *Service {
	return &Service{
		client:    client,
		jwtSecret: jwtSecret,
	}
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.client.SysUser.Query().
		Where(sysuser.UsernameEQ(username)).
		Only(ctx)
	if err != nil {
		return "", err
	}
	if user.Status == 0 {
		return "", apperror.New(apperror.CodeForbidden, "用户已停用")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", apperror.New(apperror.CodeInvalidCredential, "用户名或密码错误")
	}

	roleCodes, err := s.GetRoleCodes(ctx, user.ID)
	if err != nil {
		return "", err
	}

	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Roles:    roleCodes,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   username,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *Service) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, apperror.New(apperror.CodeUnauthorized, "认证令牌无效")
	}
	return claims, nil
}

func (s *Service) GetUserInfo(ctx context.Context, userID int) (*authmodel.UserInfoResponse, error) {
	user, err := s.client.SysUser.Query().
		Where(sysuser.IDEQ(userID)).
		WithOrganization().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	roleCodes, err := s.GetRoleCodes(ctx, userID)
	if err != nil {
		return nil, err
	}

	orgName := ""
	if user.Edges.Organization != nil {
		orgName = user.Edges.Organization.Name
	}

	return &authmodel.UserInfoResponse{
		UserID:       user.ID,
		Username:     user.Username,
		RealName:     user.RealName,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		HomePath:     "/dashboard/monitor",
		Roles:        roleCodes,
		Organization: orgName,
	}, nil
}

func (s *Service) GetRoleCodes(ctx context.Context, userID int) ([]string, error) {
	userRoleList, err := s.client.SysUserRole.Query().
		Where(sysuserrole.UserID(userID)).
		WithRole().
		All(ctx)
	if err != nil {
		return nil, err
	}

	roleCodes := make([]string, 0, len(userRoleList))
	for _, item := range userRoleList {
		if item.Edges.Role != nil {
			roleCodes = append(roleCodes, item.Edges.Role.Code)
		}
	}
	sort.Strings(roleCodes)
	return roleCodes, nil
}

func (s *Service) GetPermissionCodes(ctx context.Context, userID int) ([]string, error) {
	user, err := s.client.SysUser.Query().Where(sysuser.IDEQ(userID)).Only(ctx)
	if err != nil {
		return nil, err
	}

	menuList, err := s.queryUserMenus(ctx, user)
	if err != nil {
		return nil, err
	}

	uniqueCodes := map[string]struct{}{}
	for _, menu := range menuList {
		if menu.PermissionCode != "" {
			uniqueCodes[menu.PermissionCode] = struct{}{}
		}
	}

	codes := make([]string, 0, len(uniqueCodes))
	for code := range uniqueCodes {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	return codes, nil
}

func (s *Service) GetMenuTree(ctx context.Context, userID int) ([]authmodel.MenuResponse, error) {
	user, err := s.client.SysUser.Query().Where(sysuser.IDEQ(userID)).Only(ctx)
	if err != nil {
		return nil, err
	}

	menuList, err := s.queryUserMenus(ctx, user)
	if err != nil {
		return nil, err
	}

	menuMap := make(map[int]authmodel.MenuResponse, len(menuList))
	childrenMap := make(map[int][]authmodel.MenuResponse)
	var roots []authmodel.MenuResponse

	for _, menu := range menuList {
		if menu.Type != "catalog" && menu.Type != "menu" {
			continue
		}
		menuMap[menu.ID] = authmodel.MenuResponse{
			Name:      menu.Name,
			Path:      menu.Path,
			Component: menu.Component,
			Redirect:  menu.Redirect,
			Meta: authmodel.MenuMeta{
				Title:  menu.Title,
				Icon:   menu.Icon,
				Hidden: menu.Hidden,
			},
		}
	}

	for _, menu := range menuList {
		if menu.Type != "catalog" && menu.Type != "menu" {
			continue
		}
		node := menuMap[menu.ID]
		if menu.ParentID == nil || *menu.ParentID == 0 {
			roots = append(roots, node)
			continue
		}
		childrenMap[*menu.ParentID] = append(childrenMap[*menu.ParentID], node)
	}

	var build func(nodes []authmodel.MenuResponse) []authmodel.MenuResponse
	build = func(nodes []authmodel.MenuResponse) []authmodel.MenuResponse {
		result := make([]authmodel.MenuResponse, 0, len(nodes))
		for _, node := range nodes {
			if children, ok := childrenMap[findMenuIDByRouteName(menuMap, node.Name)]; ok {
				node.Children = build(children)
			}
			result = append(result, node)
		}
		return result
	}

	return build(roots), nil
}

func (s *Service) queryUserMenus(ctx context.Context, user *ent.SysUser) ([]*ent.SysMenu, error) {
	if user.IsSuperAdmin {
		return s.client.SysMenu.Query().
			Order(ent.Asc(sysmenu.FieldSort), ent.Asc(sysmenu.FieldID)).
			All(ctx)
	}

	roleIDs, err := s.client.SysUserRole.Query().
		Where(sysuserrole.UserID(user.ID)).
		Select(sysuserrole.FieldRoleID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	if len(roleIDs) == 0 {
		return []*ent.SysMenu{}, nil
	}

	menuIDs, err := s.client.SysRoleMenu.Query().
		Where(sysrolemenu.RoleIDIn(roleIDs...)).
		Select(sysrolemenu.FieldMenuID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	if len(menuIDs) == 0 {
		return []*ent.SysMenu{}, nil
	}

	return s.client.SysMenu.Query().
		Where(sysmenu.IDIn(menuIDs...)).
		Order(ent.Asc(sysmenu.FieldSort), ent.Asc(sysmenu.FieldID)).
		All(ctx)
}

func findMenuIDByRouteName(menuMap map[int]authmodel.MenuResponse, name string) int {
	for id, menu := range menuMap {
		if menu.Name == name {
			return id
		}
	}
	return 0
}

func NormalizeToken(token string) string {
	return strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))
}
