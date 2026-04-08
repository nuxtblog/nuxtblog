package user

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	apiv1 "github.com/nuxtblog/nuxtblog/api/user/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/event"
	"github.com/nuxtblog/nuxtblog/internal/event/payload"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	eng "github.com/nuxtblog/nuxtblog/internal/pluginsys"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/nuxtblog/nuxtblog/internal/util/password"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sUser struct{}

func New() service.IUser { return &sUser{} }

func init() {
	service.RegisterUser(New())
}

func hashPassword(plain string) (string, error) {
	return password.Hash(plain)
}

func checkPassword(plain, hash string) bool {
	return password.Verify(plain, hash)
}

func (s *sUser) GetList(ctx context.Context, req *apiv1.UserGetListReq) (*apiv1.UserGetListRes, error) {
	m := dao.Users.Ctx(ctx).WhereNull("deleted_at")
	if req.Role != nil {
		m = m.Where("role", int(*req.Role))
	}
	if req.Status != nil {
		m = m.Where("status", int(*req.Status))
	}
	if req.Keyword != nil && *req.Keyword != "" {
		kw := fmt.Sprintf("%%%s%%", *req.Keyword)
		m = m.WhereOrLike("username", kw).WhereOrLike("email", kw).WhereOrLike("display_name", kw)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	res := &apiv1.UserGetListRes{
		List:  []*apiv1.UserItem{},
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}
	if total == 0 {
		return res, nil
	}

	type UserRow struct {
		Id          int64       `orm:"id"`
		Username    string      `orm:"username"`
		Email       string      `orm:"email"`
		DisplayName string      `orm:"display_name"`
		AvatarId    *int64      `orm:"avatar_id"`
		Bio         string      `orm:"bio"`
		Role        int         `orm:"role"`
		Status      int         `orm:"status"`
		Locale      string      `orm:"locale"`
		CreatedAt   *gtime.Time `orm:"created_at"`
	}
	var rows []UserRow
	err = m.Page(req.Page, req.Size).OrderDesc("created_at").Scan(&rows)
	if err != nil {
		return nil, err
	}

	// Collect avatar IDs to batch-resolve URLs
	avatarIdSet := map[int64]bool{}
	for _, row := range rows {
		if row.AvatarId != nil && *row.AvatarId > 0 {
			avatarIdSet[*row.AvatarId] = true
		}
	}
	avatarUrlMap := map[int64]string{}
	if len(avatarIdSet) > 0 {
		ids := make([]int64, 0, len(avatarIdSet))
		for id := range avatarIdSet {
			ids = append(ids, id)
		}
		type MediaRow struct {
			Id     int64  `orm:"id"`
			CdnUrl string `orm:"cdn_url"`
		}
		var mediaRows []MediaRow
		_ = dao.Medias.Ctx(ctx).Fields("id, cdn_url").WhereIn("id", ids).Scan(&mediaRows)
		for _, mr := range mediaRows {
			if mr.CdnUrl != "" {
				avatarUrlMap[mr.Id] = mr.CdnUrl
			}
		}
	}

	for _, row := range rows {
		item := &apiv1.UserItem{
			Id:          row.Id,
			Username:    row.Username,
			Email:       row.Email,
			DisplayName: row.DisplayName,
			Bio:         row.Bio,
			Role:        apiv1.UserRole(row.Role),
			Status:      apiv1.UserStatus(row.Status),
			Locale:      row.Locale,
			AvatarId:    row.AvatarId,
			CreatedAt:   row.CreatedAt,
		}
		if row.AvatarId != nil {
			if url, ok := avatarUrlMap[*row.AvatarId]; ok {
				item.Avatar = &url
			}
		}
		res.List = append(res.List, item)
	}
	return res, nil
}

func (s *sUser) GetOne(ctx context.Context, id int64) (*apiv1.UserItem, error) {
	u, err := dao.Users.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}
	item := &apiv1.UserItem{
		Id:          int64(u.Id),
		Username:    u.Username,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		Bio:         u.Bio,
		Role:        apiv1.UserRole(u.Role),
		Status:      apiv1.UserStatus(u.Status),
		Locale:      u.Locale,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
	if u.AvatarId != 0 {
		aid := int64(u.AvatarId)
		item.AvatarId = &aid
		// Resolve avatar URL from media table
		type MediaRow struct {
			CdnUrl string `orm:"cdn_url"`
		}
		var media MediaRow
		if e := dao.Medias.Ctx(ctx).Fields("cdn_url").Where("id", u.AvatarId).Scan(&media); e == nil && media.CdnUrl != "" {
			item.Avatar = &media.CdnUrl
		}
	}

	// Fetch user_profiles
	type ProfileRow struct {
		Location    string `orm:"location"`
		Website     string `orm:"website"`
		Github      string `orm:"github"`
		Twitter     string `orm:"twitter"`
		SocialLinks string `orm:"social_links"`
	}
	var profile ProfileRow
	if e := dao.UserProfiles.Ctx(ctx).
		Fields("location", "website", "github", "twitter", "social_links").
		Where("user_id", id).Scan(&profile); e == nil {
		metas := map[string]string{}
		if profile.Location != "" {
			metas["location"] = profile.Location
		}
		if profile.Website != "" {
			metas["website"] = profile.Website
		}
		if profile.Github != "" {
			metas["github"] = profile.Github
		}
		if profile.Twitter != "" {
			metas["twitter"] = profile.Twitter
		}
		if profile.SocialLinks != "" {
			var extra map[string]string
			if json.Unmarshal([]byte(profile.SocialLinks), &extra) == nil {
				for k, v := range extra {
					if v != "" {
						metas[k] = v
					}
				}
			}
		}
		if len(metas) > 0 {
			item.Metas = metas
		}
		// Expose cover from metas
		if cover, ok := metas["cover"]; ok {
			item.Cover = &cover
		}
	}

	return item, nil
}

func (s *sUser) Create(ctx context.Context, req *apiv1.UserCreateReq) (int64, error) {
	count, err := dao.Users.Ctx(ctx).Where("username", req.Username).WhereNull("deleted_at").Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, gerror.NewCode(gcode.CodeBusinessValidationFailed, g.I18n().T(ctx, "user.username_exists"))
	}
	count, err = dao.Users.Ctx(ctx).Where("email", req.Email).WhereNull("deleted_at").Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, gerror.NewCode(gcode.CodeBusinessValidationFailed, g.I18n().T(ctx, "user.email_registered"))
	}
	hash, err := hashPassword(req.Password)
	if err != nil {
		return 0, err
	}
	result, err := dao.Users.Ctx(ctx).Insert(g.Map{
		"username":      req.Username,
		"email":         req.Email,
		"password_hash": hash,
		"display_name":  req.DisplayName,
		"role":          int(req.Role),
		"status":        int(apiv1.UserStatusActive),
		"locale":        req.Locale,
	})
	if err != nil {
		return 0, uniqueConstraintError(err)
	}
	id, _ := result.LastInsertId()

	_ = event.Emit(ctx, event.UserRegistered, payload.UserRegistered{
		UserID:      id,
		Username:    req.Username,
		Email:       req.Email,
		DisplayName: req.DisplayName,
		Locale:      req.Locale,
		Role:        int(req.Role),
	})

	return id, nil
}

func (s *sUser) Update(ctx context.Context, req *apiv1.UserUpdateReq) error {
	// Run plugin filter:user.update — allows plugins to modify or reject user updates
	{
		filterIn := map[string]any{}
		if req.DisplayName != nil {
			filterIn["display_name"] = *req.DisplayName
		}
		if req.Bio != nil {
			filterIn["bio"] = *req.Bio
		}
		if req.Locale != nil {
			filterIn["locale"] = *req.Locale
		}
		if req.Status != nil {
			filterIn["status"] = int(*req.Status)
		}
		if len(filterIn) > 0 {
			if filtered, ferr := eng.Filter(ctx, eng.FilterUserUpdate, filterIn); ferr != nil {
				return ferr
			} else {
				if v, ok := filtered["display_name"].(string); ok && req.DisplayName != nil {
					req.DisplayName = &v
				}
				if v, ok := filtered["bio"].(string); ok && req.Bio != nil {
					req.Bio = &v
				}
				if v, ok := filtered["locale"].(string); ok && req.Locale != nil {
					req.Locale = &v
				}
			}
		}
	}

	// Update users table
	data := g.Map{}
	if req.DisplayName != nil {
		data["display_name"] = *req.DisplayName
	}
	if req.Bio != nil {
		data["bio"] = *req.Bio
	}
	if req.AvatarId != nil {
		data["avatar_id"] = *req.AvatarId
	}
	if req.Locale != nil {
		data["locale"] = *req.Locale
	}
	if req.Status != nil {
		data["status"] = int(*req.Status)
	}
	if req.Role != nil {
		callerRole := middleware.GetCurrentUserRole(ctx)
		if !service.Permission().Can(ctx, callerRole, "promote_users") {
			return gerror.NewCode(gcode.CodeNotAuthorized, "permission denied: promote_users")
		}
		data["role"] = int(*req.Role)
	}
	if len(data) > 0 {
		_, err := dao.Users.Ctx(ctx).Where("id", req.Id).WhereNull("deleted_at").Update(data)
		if err != nil {
			return err
		}
	}

	// Update user_profiles table if any profile field provided
	hasProfile := req.Location != nil || req.Website != nil || req.Github != nil ||
		req.Twitter != nil || req.Instagram != nil || req.Linkedin != nil ||
		req.Youtube != nil || req.Cover != nil || req.CoverId != nil
	if !hasProfile {
		return nil
	}

	// Read existing profile to merge social_links JSON
	type ProfileRow struct {
		SocialLinks string `orm:"social_links"`
	}
	var existing ProfileRow
	dao.UserProfiles.Ctx(ctx).Fields("social_links").Where("user_id", req.Id).Scan(&existing)

	extra := map[string]string{}
	if existing.SocialLinks != "" {
		json.Unmarshal([]byte(existing.SocialLinks), &extra)
	}

	profileData := g.Map{}
	if req.Location != nil {
		profileData["location"] = *req.Location
	}
	if req.Website != nil {
		profileData["website"] = *req.Website
	}
	if req.Github != nil {
		profileData["github"] = *req.Github
	}
	if req.Twitter != nil {
		profileData["twitter"] = *req.Twitter
	}

	// Extra social links stored in JSON
	setExtra := func(key string, val *string) {
		if val == nil {
			return
		}
		if *val == "" {
			delete(extra, key)
		} else {
			extra[key] = *val
		}
	}
	setExtra("instagram", req.Instagram)
	setExtra("linkedin", req.Linkedin)
	setExtra("youtube", req.Youtube)
	setExtra("cover", req.Cover)
	if req.CoverId != nil {
		if *req.CoverId == 0 {
			delete(extra, "cover_id")
		} else {
			extra["cover_id"] = fmt.Sprintf("%d", *req.CoverId)
		}
	}

	if len(extra) > 0 {
		b, _ := json.Marshal(extra)
		profileData["social_links"] = string(b)
	} else {
		profileData["social_links"] = ""
	}

	if len(profileData) == 0 {
		goto emitUpdate
	}

	// Upsert user_profiles
	{
		result, err := dao.UserProfiles.Ctx(ctx).Where("user_id", req.Id).Update(profileData)
		if err != nil {
			return err
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			profileData["user_id"] = req.Id
			if _, err = dao.UserProfiles.Ctx(ctx).Insert(profileData); err != nil {
				return err
			}
		}
	}

emitUpdate:
	{
		type snapRow struct {
			Username    string `orm:"username"`
			Email       string `orm:"email"`
			DisplayName string `orm:"display_name"`
			Locale      string `orm:"locale"`
			Role        int    `orm:"role"`
			Status      int    `orm:"status"`
		}
		var snap snapRow
		if e := dao.Users.Ctx(ctx).Where("id", req.Id).Scan(&snap); e == nil {
			_ = event.Emit(ctx, event.UserUpdated, payload.UserUpdated{
				UserID:      req.Id,
				Username:    snap.Username,
				Email:       snap.Email,
				DisplayName: snap.DisplayName,
				Locale:      snap.Locale,
				Role:        snap.Role,
				Status:      snap.Status,
			})
		}
	}
	return nil
}

func (s *sUser) Delete(ctx context.Context, id int64) error {
	type snapRow struct {
		Username string `orm:"username"`
		Email    string `orm:"email"`
	}
	var snap snapRow
	_ = dao.Users.Ctx(ctx).Where("id", id).WhereNull("deleted_at").Scan(&snap)

	_, err := dao.Users.Ctx(ctx).Where("id", id).WhereNull("deleted_at").Update(g.Map{
		"deleted_at": gtime.Now(),
	})
	if err != nil {
		return err
	}
	_ = event.Emit(ctx, event.UserDeleted, payload.UserDeleted{
		UserID:   id,
		Username: snap.Username,
		Email:    snap.Email,
	})
	return nil
}

func (s *sUser) ChangePassword(ctx context.Context, req *apiv1.UserChangePasswordReq) error {
	u, err := dao.Users.GetById(ctx, req.Id)
	if err != nil {
		return err
	}
	if u == nil {
		return gerror.NewCode(gcode.CodeNotFound, g.I18n().T(ctx, "error.user_not_found"))
	}

	// Editors and above resetting another user's password skip the old-password check.
	// Regular users changing their own password must provide the old password,
	// UNLESS they are OAuth-only users (empty password_hash) setting a password for the first time.
	callerID, _ := middleware.GetCurrentUserID(ctx)
	callerRole := middleware.GetCurrentUserRole(ctx)
	isAdminReset := callerRole >= middleware.RoleEditor && callerID != req.Id
	isOAuthOnly := u.PasswordHash == ""
	if !isAdminReset && !isOAuthOnly {
		if req.OldPassword == "" {
			return gerror.NewCode(gcode.CodeValidationFailed, g.I18n().T(ctx, "user.old_password_required"))
		}
		if !checkPassword(req.OldPassword, u.PasswordHash) {
			return gerror.NewCode(gcode.CodeBusinessValidationFailed, g.I18n().T(ctx, "user.old_password_wrong"))
		}
	}

	newHash, err := hashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	_, err = dao.Users.Ctx(ctx).Where("id", req.Id).Update(g.Map{
		"password_hash": newHash,
	})
	return err
}

// uniqueConstraintError converts a raw DB unique-constraint error into a
// user-friendly business error. Works for both SQLite and MySQL/PostgreSQL.
func uniqueConstraintError(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	switch {
	case strings.Contains(msg, "users.email") || strings.Contains(msg, "email"):
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, g.I18n().T(context.Background(), "user.email_registered"))
	case strings.Contains(msg, "users.username") || strings.Contains(msg, "username"):
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, g.I18n().T(context.Background(), "user.username_exists"))
	case strings.Contains(msg, "UNIQUE") || strings.Contains(msg, "Duplicate entry") || strings.Contains(msg, "unique"):
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, g.I18n().T(context.Background(), "user.data_exists"))
	}
	return err
}
