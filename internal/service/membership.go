package service

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/order/v1"
)

type IMembership interface {
	// ListTiers returns all active membership tiers (public).
	ListTiers(ctx context.Context) (*v1.MembershipTierListRes, error)
	// GetUserMembership returns the current user's active membership.
	GetUserMembership(ctx context.Context, userID int64) (*v1.UserMembershipRes, error)
	// Activate activates a membership for a user after payment.
	Activate(ctx context.Context, userID int64, tierID int64, orderID int64) error
	// CheckAccess checks if the user has membership access to paid content.
	CheckAccess(ctx context.Context, userID int64) (hasAccess bool, discountPct int, err error)
	// ExpireCheck checks and expires overdue memberships (cron).
	ExpireCheck(ctx context.Context)
	// AdminCreateTier creates a new membership tier.
	AdminCreateTier(ctx context.Context, req *v1.MembershipTierCreateReq) (*v1.MembershipTierItem, error)
	// AdminUpdateTier updates a membership tier.
	AdminUpdateTier(ctx context.Context, req *v1.MembershipTierUpdateReq) error
	// AdminDeleteTier deletes a membership tier.
	AdminDeleteTier(ctx context.Context, id int64) error
	// AdminListSubscribers lists users with active memberships.
	AdminListSubscribers(ctx context.Context, req *v1.MembershipSubscriberListReq) (*v1.MembershipSubscriberListRes, error)
}

var _membership IMembership

func Membership() IMembership {
	if _membership == nil {
		panic("implement not found for interface IMembership, forgot register?")
	}
	return _membership
}

func RegisterMembership(i IMembership) { _membership = i }
