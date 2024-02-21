// ================================================
// 级联资源管理相关错误码
// ================================================

package errcode

//[122000,122999] 级联相关错误码
const (
	_ = iota + codeMpusrvCascadeStart //122000
	//无效的级联ID
	codeMpusrvCascadeInvalidCascadeId
	//禁止删除本级级联信息
	codeMpusrvCascadeDelDefaultCascadePlatForbid
	//级联平台资源使用中,禁止删除
	codeMpusrvCascadeDelInuseCascadePlatForbid
	//查询本级级联信息失败
	codeMpusrvCascadeGetDefaultCascadePlatFailed
	//本级级联ID更新,资源同步中,请稍后
	codeMpusrvCascadeDefaultCascadeIdResetting
	//本级级联信息错误
	codeMpusrvCascadeDefaultCascadeInfoErr
	//级联ID与其他平台重复
	codeMpusrvCascadeCheckCascadeIdDuplicateWithOtherCascadePlat
	//级联平台不在线
	codeMpusrvCascadeOffline
	//添加级联平台失败
	codeMpusrvCascadeAddLowerCascadePlatFailed
	//获取级联台级联ID失败
	codeMpusrvCascadeGetLowerCascadePlatIdFailed
	//校验级联ID是否重复失败
	codeMpusrvCascadeCheckCascadeIdDuplicateErr
	//已与其他平台建立级联关系,禁止修改级联ID
	codeMpusrvCascadeRestCascadeIdFailedCascadeRelationEstablished

	//级联推送信号源成功
	codeMpusrvCascadePushChnSuccess
	//级联推送信号源失败
	codeMpusrvCascadePushChnFailed
	//接收上级推送信号源失败
	codeMpusrvCascadeRecvPushChnFailed
	//重复推送
	codeMpusrvCascadePushChnDuplicate
	//不能将下级的信号源推给下级本身
	codeMpusrvCascadePushSubordinateChnToSubordinate
	//不可推送未关联网络信号源的模拟信号
	codeMpusrvCascadePushAnalogChnNotAssociationChn
	//包含不支持推送的通道
	codeMpusrvCascadePushChnContainsUnsupportedChn
	//级联平台IP重复
	codeMpusrvCascadeCheckCascadePlatIpDuplicate
	//只能向直属下级推送信号源
	codeMpusrvCascadePushToDirectSubordinatePlatOnly
	//向下级推送点位个数达到上限
	codeMpusrvCascadePushChannelMaxLimit
	//点位未推送不可撤销
	codeMpusrvCascadePushChnNotFind

	//级联获取关联网络信号源失败
	codeMpusrvCascadeGetAssociationChnFailed
	//模拟信号未关联网络信号源
	codeMpusrvCascadeAnalogChnNotAssociationChn
)

var (
	ErrMpusrvCascadeInvalidCascadeId                              = NewMpuapsErr(codeMpusrvCascadeInvalidCascadeId, "{#ErrMpusrvCascadeInvalidCascadeId}", nil)
	ErrMpusrvCascadeDelDefaultCascadePlatForbid                   = NewMpuapsErr(codeMpusrvCascadeDelDefaultCascadePlatForbid, "{#ErrMpusrvCascadeDelDefaultCascadePlatForbid}", nil)
	ErrMpusrvCascadeDelInuseCascadePlatForbid                     = NewMpuapsErr(codeMpusrvCascadeDelInuseCascadePlatForbid, "{#ErrMpusrvCascadeDelInuseCascadePlatForbid}", nil)
	ErrMpusrvCascadeGetDefaultCascadePlatFailed                   = NewMpuapsErr(codeMpusrvCascadeGetDefaultCascadePlatFailed, "{#ErrMpusrvCascadeGetDefaultCascadePlatFailed}", nil)
	ErrMpusrvCascadeDefaultCascadeIdResetting                     = NewMpuapsErr(codeMpusrvCascadeDefaultCascadeIdResetting, "{#ErrMpusrvCascadeDefaultCascadeIdResetting}", nil)
	ErrMpusrvCascadeGenCascadeIdErr                               = NewMpuapsErr(codeMpusrvCascadeDefaultCascadeInfoErr, "{#ErrMpusrvCascadeGenCascadeIdErr}", nil)
	ErrMpusrvCascadeCheckCascadeIdDuplicateWithOtherCascadePlat   = NewMpuapsErr(codeMpusrvCascadeCheckCascadeIdDuplicateWithOtherCascadePlat, "{#ErrMpusrvCascadeCheckCascadeIdDuplicateWithOtherCascadePlat}", nil)
	ErrMpusrvCascadeOffline                                       = NewMpuapsErr(codeMpusrvCascadeOffline, "{#ErrMpusrvCascadeOffline}", nil)
	ErrMpusrvCascadeAddLowerCascadePlatFailed                     = NewMpuapsErr(codeMpusrvCascadeAddLowerCascadePlatFailed, "{#ErrMpusrvCascadeAddLowerCascadePlatFailed}", nil)
	ErrMpusrvCascadeGetLowerCascadePlatIdFailed                   = NewMpuapsErr(codeMpusrvCascadeGetLowerCascadePlatIdFailed, "{#ErrMpusrvCascadeGetLowerCascadePlatIdFailed}", nil)
	ErrMpusrvCascadeCheckCascadeIdDuplicateErr                    = NewMpuapsErr(codeMpusrvCascadeCheckCascadeIdDuplicateErr, "{#ErrMpusrvCascadeCheckCascadeIdDuplicateErr}", nil)
	ErrMpusrvCascadeRestCascadeIdFailedCascadeRelationEstablished = NewMpuapsErr(codeMpusrvCascadeRestCascadeIdFailedCascadeRelationEstablished, "{#ErrMpusrvCascadeRestCascadeIdFailedCascadeRelationEstablished}", nil)

	ErrMpusrvCascadePushChnSuccess                  = NewMpuapsErr(codeMpusrvCascadePushChnSuccess, "{#ErrMpusrvCascadePushChnSuccess}", nil)
	ErrMpusrvCascadePushChnFailed                   = NewMpuapsErr(codeMpusrvCascadePushChnFailed, "{#ErrMpusrvCascadePushChnFailed}", nil)
	ErrMpusrvCascadeRecvPushChnFailed               = NewMpuapsErr(codeMpusrvCascadeRecvPushChnFailed, "{#ErrMpusrvCascadeRecvPushChnFailed}", nil)
	ErrMpusrvCascadePushChnDuplicate                = NewMpuapsErr(codeMpusrvCascadePushChnDuplicate, "{#ErrMpusrvCascadePushChnDuplicate}", nil)
	ErrMpusrvCascadePushSubordinateChnToSubordinate = NewMpuapsErr(codeMpusrvCascadePushSubordinateChnToSubordinate, "{#ErrMpusrvCascadePushSubordinateChnToSubordinate}", nil)
	ErrMpusrvCascadePushAnalogChnNotAssociationChn  = NewMpuapsErr(codeMpusrvCascadePushAnalogChnNotAssociationChn, "{#ErrMpusrvCascadePushAnalogChnNotAssociationChn}", nil)
	ErrMpusrvCascadePushChnContainsUnsupportedChn   = NewMpuapsErr(codeMpusrvCascadePushChnContainsUnsupportedChn, "{#ErrMpusrvCascadePushChnContainsUnsupportedChn}", nil)
	ErrMpusrvCascadeCheckCascadePlatIpDuplicate     = NewMpuapsErr(codeMpusrvCascadeCheckCascadePlatIpDuplicate, "{#ErrMpusrvCascadeCheckCascadePlatIpDuplicate}", nil)
	ErrMpusrvCascadePushToDirectSubordinatePlatOnly = NewMpuapsErr(codeMpusrvCascadePushToDirectSubordinatePlatOnly, "{#ErrMpusrvCascadePushToDirectSubordinatePlatOnly}", nil)
	ErrMpusrvCascadePushChannelMaxLimit             = NewMpuapsErr(codeMpusrvCascadePushChannelMaxLimit, "{#ErrMpusrvCascadePushChannelMaxLimit}", nil)
	ErrMpusrvCascadePushChnNotFind                  = NewMpuapsErr(codeMpusrvCascadePushChnNotFind, "{#ErrMpusrvCascadePushChnNotFind}", nil)

	ErrMpusrvCascadeGetAssociationChnFailed    = NewMpuapsErr(codeMpusrvCascadeGetAssociationChnFailed, "{#ErrMpusrvCascadeGetAssociationChnFailed}", nil)
	ErrMpusrvCascadeAnalogChnNotAssociationChn = NewMpuapsErr(codeMpusrvCascadeAnalogChnNotAssociationChn, "{#ErrMpusrvCascadeAnalogChnNotAssociationChn}", nil)
)
