package checkin

import (
	"context"
	"time"

	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/event"
	"github.com/nuxtblog/nuxtblog/internal/event/payload"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

type sCheckin struct{}

func New() service.ICheckin { return &sCheckin{} }

func init() {
	service.RegisterCheckin(New())
}

const actionCheckin = "checkin"

func (s *sCheckin) DoCheckin(ctx context.Context, userID int64) (alreadyCheckedIn bool, streak int, err error) {
	today := time.Now().Format("2006-01-02")

	alreadyCheckedIn, err = dao.UserActions.CheckedInToday(ctx, userID, today)
	if err != nil {
		return false, 0, err
	}
	if alreadyCheckedIn {
		streak, _, err = dao.UserProfiles.GetCheckinCache(ctx, userID)
		if err != nil {
			return true, 0, err
		}
		if streak == 0 {
			streak, err = s.calcStreak(ctx, userID, today)
			if err != nil {
				return true, 0, err
			}
			_ = dao.UserProfiles.UpdateCheckinStreak(ctx, userID, streak, today)
		}
		return true, streak, nil
	}

	if _, err = dao.UserActions.InsertAction(ctx, userID, actionCheckin, "", 0, "{}"); err != nil {
		return false, 0, err
	}

	streak, err = s.calcStreak(ctx, userID, today)
	if err != nil {
		return false, 0, err
	}

	_ = dao.UserProfiles.UpdateCheckinStreak(ctx, userID, streak, today)

	_ = event.Emit(ctx, event.CheckinDone, payload.CheckinDone{
		UserID: userID,
		Streak: streak,
	})

	return false, streak, nil
}

func (s *sCheckin) GetStatus(ctx context.Context, userID int64) (checkedInToday bool, streak int, err error) {
	today := time.Now().Format("2006-01-02")

	checkedInToday, err = dao.UserActions.CheckedInToday(ctx, userID, today)
	if err != nil {
		return false, 0, err
	}

	streak, _, err = dao.UserProfiles.GetCheckinCache(ctx, userID)
	if err != nil {
		return checkedInToday, 0, err
	}

	// Cache miss: calculate from user_actions and refresh the cache.
	if streak == 0 && checkedInToday {
		streak, err = s.calcStreak(ctx, userID, today)
		if err != nil {
			return checkedInToday, 0, err
		}
		_ = dao.UserProfiles.UpdateCheckinStreak(ctx, userID, streak, today)
	}

	return checkedInToday, streak, nil
}

// calcStreak counts consecutive check-in days ending on today.
func (s *sCheckin) calcStreak(ctx context.Context, userID int64, today string) (int, error) {
	dates, err := dao.UserActions.CheckinDates(ctx, userID, 365)
	if err != nil || len(dates) == 0 {
		return 1, err
	}

	streak := 0
	cursor, _ := time.Parse("2006-01-02", today)
	for _, d := range dates {
		t, err := time.Parse("2006-01-02", d)
		if err != nil {
			break
		}
		if !t.Equal(cursor) {
			break
		}
		streak++
		cursor = cursor.AddDate(0, 0, -1)
	}
	if streak == 0 {
		streak = 1
	}
	return streak, nil
}
