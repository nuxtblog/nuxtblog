package listener

import (
	"context"
	"fmt"

	"github.com/nuxtblog/nuxtblog/internal/event"
	"github.com/nuxtblog/nuxtblog/internal/event/payload"
	"github.com/nuxtblog/nuxtblog/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func registerNotificationListeners() {
	event.OnAsync(event.CommentCreated, onCommentCreated)
	event.OnAsync(event.UserFollowed, onUserFollowed)
	event.OnAsync(event.PostPublished, onPostPublished)
}

// onCommentCreated notifies the post author on new comments,
// and the parent comment author on replies (when they differ).
func onCommentCreated(ctx context.Context, e event.Event) error {
	p, ok := e.Payload.(payload.CommentCreated)
	if !ok {
		return nil
	}
	objectID := p.ObjectID

	if p.ParentID != nil && p.ParentAuthorID != 0 && p.ParentAuthorID != p.PostAuthorID {
		_ = service.Notification().Create(ctx,
			"reply", "",
			nil, p.AuthorName, "",
			p.ParentAuthorID,
			p.ObjectType, &objectID,
			p.ObjectTitle, "/posts/"+p.ObjectSlug,
			p.Content,
		)
	}

	notifType := "comment"
	if p.ParentID != nil {
		notifType = "reply"
	}
	if p.PostAuthorID > 0 {
		_ = service.Notification().Create(ctx,
			notifType, "",
			nil, p.AuthorName, "",
			p.PostAuthorID,
			p.ObjectType, &objectID,
			p.ObjectTitle, "/posts/"+p.ObjectSlug,
			p.Content,
		)
	}
	return nil
}

// onUserFollowed notifies the followed user.
func onUserFollowed(ctx context.Context, e event.Event) error {
	p, ok := e.Payload.(payload.UserFollowed)
	if !ok {
		return nil
	}
	followerID := p.FollowerID
	link := fmt.Sprintf("/user/%d", followerID)
	_ = service.Notification().Create(
		ctx, "follow", "follow",
		&followerID, p.FollowerName, p.FollowerAvatar,
		p.FollowingID, "user", &followerID,
		p.FollowerName, link, "",
	)
	return nil
}

// onPostPublished notifies the post author when their post is approved/published.
func onPostPublished(ctx context.Context, e event.Event) error {
	p, ok := e.Payload.(payload.PostPublished)
	if !ok {
		return nil
	}
	postID := p.PostID
	_ = service.Notification().Create(ctx,
		"system", "approved",
		nil, "", "",
		p.AuthorID,
		"post", &postID,
		p.Title, "/posts/"+p.Slug,
		g.I18n().Tf(ctx, "notification.post_published", p.Title),
	)
	return nil
}
