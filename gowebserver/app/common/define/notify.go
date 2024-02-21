// ========================================
// 通知消息topic枚举
// ========================================

package define

import "fmt"

//ENotifyTopic 通知类型
type ENotifyTopic int32

const (
	ENotifyAlive                ENotifyTopic = 0  //保活
	ENotifyDisconnect           ENotifyTopic = 1  //断链
	ENotifyResourcesReady       ENotifyTopic = 2  //后端资源准备就绪(与本级显控数据同步完成)
	ENotifyCascadeInfoUpdate    ENotifyTopic = 3  //级联平台信息更新通知
	ENotifyCascadeChnCfgUpdate  ENotifyTopic = 4  //级联通道配置更新通知
	ENotifyAccessPlatformUpdate ENotifyTopic = 5  //接入平台信息更新通知
	ENotifyChnUpdate            ENotifyTopic = 6  //通道更新通知
	ENotifyChnGroupUpdate       ENotifyTopic = 7  //通道分组通知
	ENotifyChnGroupMemUpdate    ENotifyTopic = 8  //通道分组成员
	ENotifyTVWallUpdate         ENotifyTopic = 9  //大屏配置通知
	ENotifyTVWallSceneUpdate    ENotifyTopic = 10 //场景更新通知
	ENotifyTVWallSceneLoad      ENotifyTopic = 11 //场景切换通知
	ENotifyTVWallWindowUpdate   ENotifyTopic = 12 //大屏窗口变化通知

	ENotifyUmtChnSyncWarnState             ENotifyTopic = 15 //Umt点位同步异常通知
	ENotifyUmtPlatStateUpdate              ENotifyTopic = 16 //Umt平台状态更新通知
	ENotifyUmtChnGroupUpdate               ENotifyTopic = 17 //Umt平台设备分组更新通知
	ENotifyUmtChnUpdate                    ENotifyTopic = 18 //Umt点位更新通知
	ENotifySuperiorPushChn                 ENotifyTopic = 19 //上级推送点位通知
	ENotifyChnFavoritesUpdate              ENotifyTopic = 20 //通道收藏夹更新通知
	ENotifyChnFavoritesMemUpdate           ENotifyTopic = 21 //通道收藏夹收藏成员更新通知
	ENotifyChnCollectStateUpdate           ENotifyTopic = 22 //通道收藏状态变更通知
	ETvWallSceneLoopNotify                 ENotifyTopic = 23 //大屏场景轮循通知
	ETvWallWindowChnLoopNotify             ENotifyTopic = 24 //窗口信号源轮循通知
	ENotifyForceLogout                     ENotifyTopic = 25 //强制用户登出通知
	ENotifyMatrixDispatchSchemeGroupUpdate ENotifyTopic = 26 //矩阵预案分组更新通知
	ENotifyMatrixDispatchSchemeUpdate      ENotifyTopic = 27 //矩阵预案更新通知
	ENotifyMatrixDispatchSchemeEditUpdate  ENotifyTopic = 28 //矩阵预案编辑更新通知
	ENotifyMatrixDispatchUpdate            ENotifyTopic = 29 //矩阵调度更新通知
	ENotifyMatrixDispatchDefaultSrcUpdate  ENotifyTopic = 30 //矩阵调度默认源配置更新通知
	ENotifyAudioChnFavoritesMemUpdate      ENotifyTopic = 31 //音频收藏夹收藏成员更新通知
	ENotifyUserPermissionUpdate            ENotifyTopic = 32 //用户数据权限更新通知
	ENotifyAlarmEventUpdate                ENotifyTopic = 33 //告警更新通知
	ENotifyAlarmCfgUpdate                  ENotifyTopic = 34 //告警配置通知

	ENotifyApplianceSchemeGroupUpdate ENotifyTopic = 40 //环境设备分组更新通知
	ENotifyApplianceDevUpdate         ENotifyTopic = 41 //环境设备更新通知
	ENotifyMpuRelayUpdate             ENotifyTopic = 42 //继电器更新通知

	ENotifyMeetGroupUpdate    ENotifyTopic = 50 //会议分组更新通知
	ENotifyMeetDispatchUpdate ENotifyTopic = 51 //会议调度更新通知
	ENotifyMeetingUpdate      ENotifyTopic = 52 //会议信息更新通知

	ETvWallBatchChnLoopNotify ENotifyTopic = 70 //大屏批量通道轮循通知

	ENotifyMatrixAVInputChannelBindCfgUpdate ENotifyTopic = 80 //大屏批量通道轮循通知
)

func (x ENotifyTopic) String() string {
	switch x {
	case ENotifyAlive:
		return "emNotifyAlive"
	case ENotifyDisconnect:
		return "emNotifyDisconnect"
	case ENotifyResourcesReady:
		return "emNotifyResourcesReady"
	case ENotifyCascadeInfoUpdate:
		return "emNotifyCascadeInfoUpdate"
	case ENotifyCascadeChnCfgUpdate:
		return "emNotifyCascadeChnCfgUpdate"
	case ENotifyAccessPlatformUpdate:
		return "emNotifyAccessPlatformUpdate"
	case ENotifyChnUpdate:
		return "emNotifyChnUpdate"
	case ENotifyChnGroupUpdate:
		return "emNotifyChnGroupUpdate"
	case ENotifyChnGroupMemUpdate:
		return "emNotifyChnGroupMemUpdate"
	case ENotifyTVWallUpdate:
		return "emNotifyTVWallUpdate"
	case ENotifyTVWallSceneUpdate:
		return "emNotifyTVWallSceneUpdate"
	case ENotifyTVWallSceneLoad:
		return "emNotifyTVWallSceneLoad"
	case ENotifyTVWallWindowUpdate:
		return "emNotifyTVWallWindowUpdate"

	case ENotifyUmtChnSyncWarnState:
		return "emNotifyUmtChnSyncWarnState"
	case ENotifyUmtPlatStateUpdate:
		return "emNotifyUmtPlatStateUpdate"
	case ENotifyUmtChnGroupUpdate:
		return "emNotifyUmtChnGroupUpdate"
	case ENotifyUmtChnUpdate:
		return "emNotifyUmtChnUpdate"
	case ENotifySuperiorPushChn:
		return "emNotifySuperiorPushChn"
	case ENotifyChnCollectStateUpdate:
		return "emNotifyChnCollectStateUpdate"
	case ENotifyChnFavoritesUpdate:
		return "emNotifyChnFavoritesUpdate"
	case ENotifyChnFavoritesMemUpdate:
		return "emNotifyChnFavoritesMemUpdate"
	case ETvWallSceneLoopNotify:
		return "emTvWallSceneLoopNotify"
	case ETvWallWindowChnLoopNotify:
		return "emTvWallWindowChnLoopNotify"
	case ENotifyForceLogout:
		return "emNotifyForceLogout"
	case ENotifyMatrixDispatchSchemeGroupUpdate:
		return "emNotifyMatrixDispatchSchemeGroupUpdate"
	case ENotifyMatrixDispatchSchemeUpdate:
		return "emNotifyMatrixDispatchSchemeUpdate"
	case ENotifyMatrixDispatchSchemeEditUpdate:
		return "emNotifyMatrixDispatchSchemeEditUpdate"
	case ENotifyMatrixDispatchUpdate:
		return "emNotifyMatrixDispatchUpdate"
	case ENotifyMatrixDispatchDefaultSrcUpdate:
		return "emNotifyMatrixDispatchDefaultSrcUpdate"
	case ENotifyMatrixAVInputChannelBindCfgUpdate:
		return "emNotifyMatrixAVInputChannelBindCfgUpdate"

	case ENotifyAudioChnFavoritesMemUpdate:
		return "emNotifyAudioChnFavoritesMemUpdate"
	case ENotifyUserPermissionUpdate:
		return "emNotifyUserPermissionUpdate"
	case ENotifyAlarmEventUpdate:
		return "emNotifyAlarmEventUpdate"
	case ENotifyAlarmCfgUpdate:
		return "emNotifyAlarmCfgUpdate"
	case ENotifyApplianceSchemeGroupUpdate:
		return "emNotifyApplianceSchemeGroupUpdate"
	case ENotifyApplianceDevUpdate:
		return "emNotifyApplianceDevUpdate"
	case ENotifyMpuRelayUpdate:
		return "emNotifyMpuRelayUpdate"

	case ENotifyMeetGroupUpdate:
		return "emNotifyMeetGroupUpdate"
	case ENotifyMeetDispatchUpdate:
		return "emNotifyMeetDispatchUpdate"
	case ENotifyMeetingUpdate:
		return "emNotifyMeetingUpdate"

	case ETvWallBatchChnLoopNotify:
		return "emTvWallWindowBatchChnLoopNotify"
	default:
		return fmt.Sprintf("ENotifyTopic(%d) undefined, Please check file \"app/common/define/notify.go\"", x)
	}
}
