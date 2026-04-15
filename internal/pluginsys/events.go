package pluginsys

import (
	"context"

	"github.com/nuxtblog/nuxtblog/internal/event"
	"github.com/nuxtblog/nuxtblog/internal/event/payload"
)

// RegisterEventListeners subscribes to the global event bus and fans out events
// to all currently loaded plugin runtimes. Call once at startup.
func RegisterEventListeners() {
	type mapping struct {
		name  string
		toMap func(any) map[string]any
	}
	for _, m := range []mapping{
		// ── Post ────────────────────────────────────────────────────────────
		{
			event.PostCreated,
			func(p any) map[string]any {
				pp := p.(payload.PostCreated)
				return map[string]any{
					"id": pp.PostID, "title": pp.Title, "slug": pp.Slug,
					"excerpt": pp.Excerpt, "post_type": pp.PostType,
					"author_id": pp.AuthorID, "status": pp.Status,
				}
			},
		},
		{
			event.PostUpdated,
			func(p any) map[string]any {
				pp := p.(payload.PostUpdated)
				return map[string]any{
					"id": pp.PostID, "title": pp.Title, "slug": pp.Slug,
					"excerpt": pp.Excerpt, "post_type": pp.PostType,
					"author_id": pp.AuthorID, "status": pp.Status,
				}
			},
		},
		{
			event.PostPublished,
			func(p any) map[string]any {
				pp := p.(payload.PostPublished)
				return map[string]any{
					"id": pp.PostID, "title": pp.Title, "slug": pp.Slug,
					"excerpt": pp.Excerpt, "post_type": pp.PostType,
					"author_id": pp.AuthorID,
				}
			},
		},
		{
			event.PostDeleted,
			func(p any) map[string]any {
				pp := p.(payload.PostDeleted)
				return map[string]any{
					"id": pp.PostID, "title": pp.Title, "slug": pp.Slug,
					"post_type": pp.PostType, "author_id": pp.AuthorID,
				}
			},
		},
		// ── Comment ──────────────────────────────────────────────────────────
		{
			event.CommentCreated,
			func(p any) map[string]any {
				pp := p.(payload.CommentCreated)
				m := map[string]any{
					"id": pp.CommentID, "status": pp.Status,
					"object_type": pp.ObjectType, "object_id": pp.ObjectID,
					"object_title": pp.ObjectTitle, "object_slug": pp.ObjectSlug,
					"post_author_id":   pp.PostAuthorID,
					"author_id":        pp.AuthorID,
					"author_name":      pp.AuthorName,
					"author_email":     pp.AuthorEmail,
					"content":          pp.Content,
					"parent_author_id": pp.ParentAuthorID,
				}
				if pp.ParentID != nil {
					m["parent_id"] = *pp.ParentID
				}
				return m
			},
		},
		{
			event.CommentDeleted,
			func(p any) map[string]any {
				pp := p.(payload.CommentDeleted)
				return map[string]any{
					"id": pp.CommentID, "object_type": pp.ObjectType, "object_id": pp.ObjectID,
				}
			},
		},
		{
			event.CommentStatusChanged,
			func(p any) map[string]any {
				pp := p.(payload.CommentStatusChanged)
				return map[string]any{
					"id": pp.CommentID, "object_type": pp.ObjectType, "object_id": pp.ObjectID,
					"old_status": pp.OldStatus, "new_status": pp.NewStatus,
					"moderator_id": pp.ModeratorID,
				}
			},
		},
		// ── User ─────────────────────────────────────────────────────────────
		{
			event.UserRegistered,
			func(p any) map[string]any {
				pp := p.(payload.UserRegistered)
				return map[string]any{
					"id": pp.UserID, "username": pp.Username, "email": pp.Email,
					"display_name": pp.DisplayName, "locale": pp.Locale, "role": pp.Role,
				}
			},
		},
		{
			event.UserUpdated,
			func(p any) map[string]any {
				pp := p.(payload.UserUpdated)
				return map[string]any{
					"id": pp.UserID, "username": pp.Username, "email": pp.Email,
					"display_name": pp.DisplayName, "locale": pp.Locale,
					"role": pp.Role, "status": pp.Status,
				}
			},
		},
		{
			event.UserDeleted,
			func(p any) map[string]any {
				pp := p.(payload.UserDeleted)
				return map[string]any{
					"id": pp.UserID, "username": pp.Username, "email": pp.Email,
				}
			},
		},
		{
			event.UserFollowed,
			func(p any) map[string]any {
				pp := p.(payload.UserFollowed)
				return map[string]any{
					"follower_id":     pp.FollowerID,
					"follower_name":   pp.FollowerName,
					"follower_avatar": pp.FollowerAvatar,
					"following_id":    pp.FollowingID,
				}
			},
		},
		// ── Media ────────────────────────────────────────────────────────────
		{
			event.MediaUploaded,
			func(p any) map[string]any {
				pp := p.(payload.MediaUploaded)
				return map[string]any{
					"id": pp.MediaID, "uploader_id": pp.UploaderID,
					"filename": pp.Filename, "mime_type": pp.MimeType,
					"file_size": pp.FileSize, "url": pp.URL,
					"category": pp.Category, "width": pp.Width, "height": pp.Height,
				}
			},
		},
		{
			event.MediaDeleted,
			func(p any) map[string]any {
				pp := p.(payload.MediaDeleted)
				return map[string]any{
					"id": pp.MediaID, "uploader_id": pp.UploaderID,
					"filename": pp.Filename, "mime_type": pp.MimeType,
					"category": pp.Category,
				}
			},
		},
		// ── Taxonomy / Term ──────────────────────────────────────────────────
		{
			event.TaxonomyCreated,
			func(p any) map[string]any {
				pp := p.(payload.TaxonomyCreated)
				return map[string]any{
					"id": pp.TaxID, "term_id": pp.TermID,
					"term_name": pp.TermName, "term_slug": pp.TermSlug,
					"taxonomy": pp.Taxonomy,
				}
			},
		},
		{
			event.TaxonomyDeleted,
			func(p any) map[string]any {
				pp := p.(payload.TaxonomyDeleted)
				return map[string]any{
					"id": pp.TaxID, "term_name": pp.TermName,
					"term_slug": pp.TermSlug, "taxonomy": pp.Taxonomy,
				}
			},
		},
		{
			event.TermCreated,
			func(p any) map[string]any {
				pp := p.(payload.TermCreated)
				return map[string]any{
					"id": pp.TermID, "name": pp.Name, "slug": pp.Slug,
				}
			},
		},
		{
			event.TermDeleted,
			func(p any) map[string]any {
				pp := p.(payload.TermDeleted)
				return map[string]any{
					"id": pp.TermID, "name": pp.Name, "slug": pp.Slug,
				}
			},
		},
		// ── Reaction / Checkin ───────────────────────────────────────────────
		{
			event.ReactionAdded,
			func(p any) map[string]any {
				pp := p.(payload.ReactionAdded)
				return map[string]any{
					"user_id": pp.UserID, "object_type": pp.ObjectType,
					"object_id": pp.ObjectID, "type": pp.Type,
				}
			},
		},
		{
			event.ReactionRemoved,
			func(p any) map[string]any {
				pp := p.(payload.ReactionRemoved)
				return map[string]any{
					"user_id": pp.UserID, "object_type": pp.ObjectType,
					"object_id": pp.ObjectID, "type": pp.Type,
				}
			},
		},
		{
			"checkin.done",
			func(p any) map[string]any {
				if m, ok := p.(map[string]any); ok {
					return m
				}
				return nil
			},
		},
		// ── Post viewed ──────────────────────────────────────────────────────
		{
			event.PostViewed,
			func(p any) map[string]any {
				pp := p.(payload.PostViewed)
				return map[string]any{"id": pp.PostID, "user_id": pp.UserID}
			},
		},
		// ── Comment approved ─────────────────────────────────────────────────
		{
			event.CommentApproved,
			func(p any) map[string]any {
				pp := p.(payload.CommentApproved)
				return map[string]any{
					"id": pp.CommentID, "object_type": pp.ObjectType,
					"object_id": pp.ObjectID, "moderator_id": pp.ModeratorID,
				}
			},
		},
		// ── User login / logout ───────────────────────────────────────────────
		{
			event.UserLoggedIn,
			func(p any) map[string]any {
				pp := p.(payload.UserLoggedIn)
				return map[string]any{
					"id": pp.UserID, "username": pp.Username,
					"email": pp.Email, "role": pp.Role,
				}
			},
		},
		{
			event.UserLoggedOut,
			func(p any) map[string]any {
				pp := p.(payload.UserLoggedOut)
				return map[string]any{"id": pp.UserID}
			},
		},
		// ── Option updated ───────────────────────────────────────────────────
		{
			event.OptionUpdated,
			func(p any) map[string]any {
				pp := p.(payload.OptionUpdated)
				return map[string]any{"key": pp.Key, "value": pp.Value}
			},
		},
		// ── Plugin lifecycle ─────────────────────────────────────────────────
		{
			event.PluginInstalled,
			func(p any) map[string]any {
				pp := p.(payload.PluginInstalled)
				return map[string]any{
					"id": pp.PluginID, "title": pp.Title,
					"version": pp.Version, "author": pp.Author,
				}
			},
		},
		{
			event.PluginUninstalled,
			func(p any) map[string]any {
				pp := p.(payload.PluginUninstalled)
				return map[string]any{"id": pp.PluginID}
			},
		},
	} {
		m := m
		event.OnAsync(m.name, func(_ context.Context, e event.Event) error {
			fanOut(e.Name, m.toMap(e.Payload))
			return nil
		})
	}
}
