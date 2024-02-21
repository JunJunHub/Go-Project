// ====================================
// [123000,123999]信号源管理相关错误码
// ====================================

package errcode

//[123000,123099] 信号源相关错误码
const (
	codeMpusrvChannelUnknown      = iota + codeMpusrvChannelStart //123000
	codeMpusrvChannelInvalidChnId                                 //123001
	codeMpusrvChannelNoPermission
	codeMpusrvChannelNoExist
)

//[123100,123299] 信号源分组管理相关错误码
const (
	codeMpusrvChannelGroupInvalidChnGroupId = iota + codeMpusrvChannelGroupStart //123100
	codeMpusrvChannelGroupInvalidParentGroupId
	codeMpusrvChannelGroupCreateSubGroupForbid

	codeMpusrvChannelGroupSysDefGroupNameUnavailable
	codeMpusrvChannelGroupNameDuplicate
	codeMpusrvChannelGroupNameInvalid
	codeMpusrvChannelGroupModifyForbid

	codeMpusrvChannelGroupMemMustBeVideoInput
	codeMpusrvChannelGroupMemMustBeVideoOutput
	codeMpusrvChannelGroupMemMustBeAudioInput
	codeMpusrvChannelGroupMemMustBeAudioOutput
)

//[123300,123399] 信号源收藏管理相关错误码
const (
	codeMpusrvChnFavoriteIdInvalid = iota + codeMpusrvChannelFavorite //123300
	codeMpusrvChnFavoriteMoveFailed
	codeMpusrvChnFavoriteNameInvalid
	codeMpusrvChnFavoriteParentNodeInvalid
	codeMpusrvChnFavoriteNameDuplicate
	codeMpusrvChnFavoriteNumLimit
	codeMpusrvChnFavoriteDelRootNodeForbid
)

//[123400,123499] 模拟信号绑定网络信号源关联配置
const (
	codeMpuSrvAssociationChnOnlySupportLocalChn = iota + codeMpusrvAssociationChn
	codeMpuSrvAssociationChnBoundOther
)

var (
	ErrMpusrvChannelUnknown      = NewMpuapsErr(codeMpusrvChannelUnknown, "{#ErrMpusrvChannelUnknown}", nil)
	ErrMpusrvChannelInvalidChnId = NewMpuapsErr(codeMpusrvChannelInvalidChnId, "{#ErrMpusrvChannelInvalidChnId}", nil)
	ErrMpusrvChannelNoPermission = NewMpuapsErr(codeMpusrvChannelNoPermission, "{#ErrMpusrvChannelNoPermission}", nil)
	ErrMpusrvChannelNoExist      = NewMpuapsErr(codeMpusrvChannelNoExist, "{#ErrMpusrvChannelNoExist}", nil)

	ErrMpusrvChannelGroupInvalidChnGroupId          = NewMpuapsErr(codeMpusrvChannelGroupInvalidChnGroupId, "{#ErrMpusrvChannelGroupInvalidChnGroupId}", nil)
	ErrMpusrvChannelGroupInvalidParentGroupId       = NewMpuapsErr(codeMpusrvChannelGroupInvalidParentGroupId, "{#ErrMpusrvChannelGroupInvalidParentGroupId}", nil)
	ErrMpusrvChannelGroupCreateSubGroupForbid       = NewMpuapsErr(codeMpusrvChannelGroupCreateSubGroupForbid, "{#ErrMpusrvChannelGroupCreateSubGroupForbid}", nil)
	ErrMpusrvChannelGroupSysDefGroupNameUnavailable = NewMpuapsErr(codeMpusrvChannelGroupSysDefGroupNameUnavailable, "{#ErrMpusrvChannelGroupSysDefGroupNameUnavailable}", nil)
	ErrMpusrvChannelGroupNameDuplicate              = NewMpuapsErr(codeMpusrvChannelGroupNameDuplicate, "{#ErrMpusrvChannelGroupNameDuplicate}", nil)
	ErrMpusrvChannelGroupNameInvalid                = NewMpuapsErr(codeMpusrvChannelGroupNameInvalid, "{#ErrMpusrvChannelGroupNameInvalid}", nil)
	ErrMpusrvChannelGroupModifyForbid               = NewMpuapsErr(codeMpusrvChannelGroupModifyForbid, "{#ErrMpusrvChannelGroupModifyForbid}", nil)

	ErrMpusrvChannelGroupMemMustBeVideoInput  = NewMpuapsErr(codeMpusrvChannelGroupMemMustBeVideoInput, "{#ErrMpusrvChannelGroupMemMustBeVideoInput}", nil)
	ErrMpusrvChannelGroupMemMustBeVideoOutput = NewMpuapsErr(codeMpusrvChannelGroupMemMustBeVideoOutput, "{#ErrMpusrvChannelGroupMemMustBeVideoOutput}", nil)
	ErrMpusrvChannelGroupMemMustBeAudioInput  = NewMpuapsErr(codeMpusrvChannelGroupMemMustBeAudioInput, "{#ErrMpusrvChannelGroupMemMustBeAudioInput}", nil)
	ErrMpusrvChannelGroupMemMustBeAudioOutput = NewMpuapsErr(codeMpusrvChannelGroupMemMustBeAudioOutput, "{#ErrMpusrvChannelGroupMemMustBeAudioOutput}", nil)

	ErrMpuSrvChnFavoriteIdInvalid         = NewMpuapsErr(codeMpusrvChnFavoriteIdInvalid, "{#ErrMpuSrvChnFavoriteIdInvalid}", nil)
	ErrMpusrvChnFavoriteMoveFailed        = NewMpuapsErr(codeMpusrvChnFavoriteMoveFailed, "{#ErrMpusrvChnFavoriteMoveFailed}", nil)
	ErrMpusrvChnFavoriteNameInvalid       = NewMpuapsErr(codeMpusrvChnFavoriteNameInvalid, "{#ErrMpusrvChnFavoriteNameInvalid}", nil)
	ErrMpusrvChnFavoriteParentNodeInvalid = NewMpuapsErr(codeMpusrvChnFavoriteParentNodeInvalid, "{#ErrMpusrvChnFavoriteParentNodeInvalid}", nil)
	ErrMpusrvChnFavoriteNameDuplicate     = NewMpuapsErr(codeMpusrvChnFavoriteNameDuplicate, "{#ErrMpusrvChnFavoriteNameDuplicate}", nil)
	ErrMpusrvChnFavoriteNumLimit          = NewMpuapsErr(codeMpusrvChnFavoriteNumLimit, "{#ErrMpusrvChnFavoriteNumLimit}", nil)
	ErrMpusrvChnFavoriteDelRootNodeForbid = NewMpuapsErr(codeMpusrvChnFavoriteDelRootNodeForbid, "{#ErrMpusrvChnFavoriteDelRootNodeForbid}", nil)

	ErrMpuSrvAssociationChnOnlySupportLocalChn = NewMpuapsErr(codeMpuSrvAssociationChnOnlySupportLocalChn, "{#ErrMpuSrvAssociationChnOnlySupportLocalChn}", nil)
	ErrMpuSrvAssociationChnBoundOther          = NewMpuapsErr(codeMpuSrvAssociationChnBoundOther, "{#ErrMpuSrvAssociationChnBoundOther}", nil)
)
