package report

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/report/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sReport struct{}

func New() service.IReport  { return &sReport{} }
func init()                 { service.RegisterReport(New()) }

func (s *sReport) Create(ctx context.Context, req *v1.ReportCreateReq) (int64, error) {
	userID, _ := middleware.GetCurrentUserID(ctx)
	if userID == 0 {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, g.I18n().T(ctx, "error.unauthorized"))
	}
	// Prevent duplicate pending report from same user
	cnt, _ := dao.Reports.Ctx(ctx).
		Where("reporter_id", userID).
		Where("target_type", req.TargetType).
		Where("target_id", req.TargetId).
		Where("status", "pending").
		Count()
	if cnt > 0 {
		return 0, gerror.NewCode(gcode.CodeBusinessValidationFailed, g.I18n().T(ctx, "report.duplicate"))
	}
	newId := idgen.New()
	reportId, err := dao.Reports.Ctx(ctx).Data(g.Map{
		"id":          newId,
		"reporter_id": userID,
		"target_type": req.TargetType,
		"target_id":   req.TargetId,
		"reason":      req.Reason,
		"detail":      req.Detail,
		"status":      "pending",
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}

	// Notify reporter: report received
	go func() {
		_ = service.Notification().Create(
			context.Background(),
			"system", "report_received",
			nil, "", "",
			userID,
			"report", &reportId,
			req.Reason, "",
			g.I18n().T(ctx, "report.received_content"),
		)
	}()

	return reportId, nil
}

// reportListRow is used for the JOIN projection in List.
type reportListRow struct {
	Id           int64   `orm:"id"`
	ReporterId   int64   `orm:"reporter_id"`
	ReporterName string  `orm:"reporter_name"`
	TargetType   string  `orm:"target_type"`
	TargetId     int64   `orm:"target_id"`
	TargetName   string  `orm:"target_name"`
	Reason       string  `orm:"reason"`
	Detail       string  `orm:"detail"`
	Status       string  `orm:"status"`
	Notes        string  `orm:"notes"`
	CreatedAt    string  `orm:"created_at"`
	ResolvedAt   *string `orm:"resolved_at"`
}

func (s *sReport) List(ctx context.Context, req *v1.ReportListReq) (*v1.ReportListRes, error) {
	callerRole := middleware.GetCurrentUserRole(ctx)
	if callerRole < middleware.RoleEditor {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, g.I18n().T(ctx, "error.forbidden"))
	}

	offset := (req.Page - 1) * req.Size

	// Build parameterized count query using ORM builder
	countM := dao.Reports.Ctx(ctx)
	if req.Status != "" && req.Status != "all" {
		countM = countM.Where("status", req.Status)
	}
	total, _ := countM.Count()

	// For the JOIN query, use Raw with parameterized placeholders (no string concat).
	// LEFT JOIN users tu to resolve user targets, LEFT JOIN posts tp to resolve post targets.
	const selectSQL = `
		SELECT r.id, r.reporter_id, COALESCE(u.display_name, u.username) AS reporter_name,
		       r.target_type, r.target_id, r.reason, r.detail, r.status, r.notes,
		       r.created_at, r.resolved_at,
		       CASE
		           WHEN r.target_type = 'user' THEN COALESCE(tu.display_name, tu.username, '')
		           WHEN r.target_type = 'post' THEN COALESCE(tp.title, '')
		           ELSE ''
		       END AS target_name
		FROM reports r
		LEFT JOIN users u  ON u.id  = r.reporter_id
		LEFT JOIN users tu ON tu.id = r.target_id AND r.target_type = 'user'
		LEFT JOIN posts tp ON tp.id = r.target_id AND r.target_type = 'post'
	`
	var rows []reportListRow
	if req.Status != "" && req.Status != "all" {
		_ = dao.Reports.DB().Ctx(ctx).Raw(
			selectSQL+` WHERE r.status = ? ORDER BY r.created_at DESC LIMIT ? OFFSET ?`,
			req.Status, req.Size, offset,
		).Scan(&rows)
	} else {
		_ = dao.Reports.DB().Ctx(ctx).Raw(
			selectSQL+` ORDER BY r.created_at DESC LIMIT ? OFFSET ?`,
			req.Size, offset,
		).Scan(&rows)
	}

	items := make([]v1.ReportItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, v1.ReportItem{
			Id: r.Id, ReporterId: r.ReporterId, ReporterName: r.ReporterName,
			TargetType: r.TargetType, TargetId: r.TargetId, TargetName: r.TargetName,
			Reason: r.Reason, Detail: r.Detail, Status: r.Status,
			Notes: r.Notes, CreatedAt: r.CreatedAt, ResolvedAt: r.ResolvedAt,
		})
	}
	return &v1.ReportListRes{Items: items, Total: total}, nil
}

func (s *sReport) Handle(ctx context.Context, req *v1.ReportHandleReq) error {
	callerRole := middleware.GetCurrentUserRole(ctx)
	if callerRole < middleware.RoleEditor {
		return gerror.NewCode(gcode.CodeNotAuthorized, g.I18n().T(ctx, "error.forbidden"))
	}

	// Fetch reporter_id + reason before updating
	type reportRow struct {
		ReporterId int64  `orm:"reporter_id"`
		Reason     string `orm:"reason"`
	}
	var r reportRow
	_ = dao.Reports.Ctx(ctx).Where("id", req.Id).Fields("reporter_id, reason").Scan(&r)

	_, err := dao.Reports.Ctx(ctx).Where("id", req.Id).Data(g.Map{
		"status":      req.Status,
		"notes":       req.Notes,
		"resolved_at": gtime.Now(),
	}).Update()
	if err != nil {
		return err
	}

	// Notify reporter with result
	if r.ReporterId > 0 {
		subType := "report_resolved"
		if req.Status == "dismissed" {
			subType = "report_dismissed"
		}
		content := req.Notes
		if content == "" {
			content = g.I18n().T(ctx, "report.handled_content_"+req.Status)
		}
		reportId := req.Id
		go func() {
			_ = service.Notification().Create(
				context.Background(),
				"system", subType,
				nil, "", "",
				r.ReporterId,
				"report", &reportId,
				r.Reason, "",
				content,
			)
		}()
	}

	return nil
}
