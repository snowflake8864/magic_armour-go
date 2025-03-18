package enum_backup

const (
	BackupStatusDefault int = 1 //备份状态: 未开始
	BackupStatusRunning int = 2 //备份状态: 进行中
	BackupStatusSuccess int = 3 //备份状态: 成功
	BackupStatusFail    int = 4 //备份状态: 失败

	RollbackStatusDefault int = 1 //回滚状态: 未开始
	RollbackStatusRunning int = 2 //回滚状态: 进行中
	RollbackStatusSuccess int = 3 //回滚状态: 成功
	RollbackStatusFail    int = 4 //回滚状态: 失败
)

var EnumIntMap = map[string]map[int]string{
	"backup_status": {
		BackupStatusDefault: "未开始",
		BackupStatusRunning: "进行中",
		BackupStatusSuccess: "成功",
		BackupStatusFail:    "失败",
	},
	"rollback_status": {
		RollbackStatusDefault: "未开始",
		RollbackStatusRunning: "进行中",
		RollbackStatusSuccess: "成功",
		RollbackStatusFail:    "失败",
	},
}

var WebMapInt = map[string][]int{
	"backup_status": {
		BackupStatusDefault,
		BackupStatusRunning,
		BackupStatusSuccess,
		BackupStatusFail,
	},
	"rollback_status": {
		RollbackStatusDefault,
		RollbackStatusRunning,
		RollbackStatusSuccess,
		RollbackStatusFail,
	},
}
