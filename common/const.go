package common

var (
	StoryCreate       = "story::create"
	StoryStatusChange = "story::status_change"
	StoryComment      = "story::comment"
	BugCreate         = "bug::create"
	BugStatusChange   = "bug::status_change"
	BugComment        = "bug::comment"
	TaskCreate        = "task::create"
	TaskStatusChange  = "task::status_change"
	TaskComment       = "task::comment"
)

var (
	CommentAdd    = "add"
	CommentUpdate = "update"
	CommentDelete = "delete"
)

var StoryStatusMap = map[string]string{
	"planning":   "规划中",
	"status_2":   "已评审",
	"rejected":   "已拒绝",
	"status_3":   "已转测",
	"status_4":   "测试中",
	"status_5":   "测试完成",
	"resolved":   "已实现",
	"developing": "实现中",
}

var BugStatusMap = map[string]string{
	"new":         "新",
	"in_progress": "接受/处理",
	"resolved":    "已解决",
	"verified":    "已验证",
	"reopened":    "重新打开",
	"rejected":    "已拒绝",
	"closed":      "已关闭",
	"postponed":   "延期",
	"suspended":   "挂起",
}

var TaskStatusMap = map[string]string{
	"open":        "未开始",
	"progressing": "进行中",
	"done":        "已完成",
}

var EventKeyMap = map[string]string{
	StoryCreate:       "需求创建",
	StoryStatusChange: "需求状态变更",
	StoryComment:      "需求评论",
	BugCreate:         "缺陷创建",
	BugStatusChange:   "缺陷状态变更",
	BugComment:        "缺陷评论",
	TaskCreate:        "任务创建",
	TaskStatusChange:  "任务状态变更",
	TaskComment:       "任务评论",
}

var CommentActionMap = map[string]string{
	CommentAdd:    "新增",
	CommentUpdate: "更新",
	CommentDelete: "删除",
}

var (
	TAPDProjectTableName  = "tapd_project"
	DingdingUserTableName = "dingding_user"
	DingdingDeptTableName = "dingding_dept"
)

const StrTimeFormat = "2006-01-02 15:04:05"
