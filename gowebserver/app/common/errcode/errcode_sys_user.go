// ===================================
// 系统用户相关错误码定义[111000,111999]
// ===================================

package errcode

//[111000,111099] 系统用户预留错误码
const (
	codeSysUserUnknown = iota + codeSysUserStart //111000
)

//[111100,111199] 系统用户管理错误码
const (
	codeSysUserLock                 = iota + codeSysUserManageStart //111100
	codeSysUserInvalidUid                                           //111101
	codeSysUserLoginNameDuplicate                                   //111102
	codeSysUserEmailDuplicate                                       //111103
	codeSysUserPhonenumberDuplicate                                 //111104
	codeSysUserDelDefSysUserForbid                                  //111105
	codeSysUserDelOnlineUserForbid                                  //111106
	codeSysUserSysAdminPasswdErr                                    //111107
	codeSysUserBindRoleUnknow                                       //111108
	codeSysUserUnBindRole                                           //111109
	codeSysUserPasswdErr                                            //111110
	codeSysUserNotLogin                                             //111111
	codeSysUserModifyDefSysUserRoleForbid
	codeSysUserPassInvalid
	codeSysUserLockDefSysUserForbid
	codeSysUserLoginNameInvalid
	codeSysUserUpdateNewPassSameAsOldPass
	codeSysUserUpdatePassNotSameAsConfirm
	codeSysUserUpdatePassInputOldPassErr
)

//[1111200,1111299] 系统用户登录错误码
const (
	codeSysUserLoginAuthFailed                              = iota + codeSysUserLoginStart //111200
	codeSysUserLoginCaptchaVerifyFailed                                                    //111201
	codeSysUserLoginLock                                                                   //111202
	codeSysUserCascadeLoginForbidRelationsUnknown                                          //111203
	codeSysUserCascadeLoginForbidSameCascadeId                                             //111204
	codeSysUserCascadeLoginForbidSubordinates                                              //111205
	codeSysUserCascadeLoginNetErr                                                          //111206
	codeSysUserCascadeLoginForbidHaveOtherParentCascadePlat                                //111207
	codeSysUserLoginAlready                                                                //111208
	codeSysUserLoginPermissionDenied                                                       //111209
	codeSysUserLoginBlackForbid                                                            //111210
	codeSysUserLoginAccessConnectionsNumExceedMax                                          //111211
	codeSysUserLoginFailed                                                                 //111212
)

//用户相关业务错误码对象
var (
	ErrSysUserUnknown                    = NewMpuapsErr(codeSysUserUnknown, "{#ErrSysUserUnknown}", nil)
	ErrSysUserLock                       = NewMpuapsErr(codeSysUserLock, "{#ErrSysUserLock}", nil)
	ErrSysUserInvalidUid                 = NewMpuapsErr(codeSysUserInvalidUid, "{#ErrSysUserInvalidUid}", nil)
	ErrSysUserLoginNameDuplicate         = NewMpuapsErr(codeSysUserLoginNameDuplicate, "{#ErrSysUserLoginNameDuplicate}", nil)
	ErrSysUserEmailDuplicate             = NewMpuapsErr(codeSysUserEmailDuplicate, "{#ErrSysUserEmailDuplicate}", nil)
	ErrSysUserPhonenumberDuplicate       = NewMpuapsErr(codeSysUserPhonenumberDuplicate, "{#ErrSysUserPhonenumberDuplicate}", nil)
	ErrSysUserDelDefSysUserForbid        = NewMpuapsErr(codeSysUserDelDefSysUserForbid, "{#ErrSysUserDelDefSysUserForbid}", nil)
	ErrSysUserDelOnlineUserForbid        = NewMpuapsErr(codeSysUserDelOnlineUserForbid, "{#ErrSysUserDelOnlineUserForbid}", nil)
	ErrSysUserSysAdminPasswdErr          = NewMpuapsErr(codeSysUserSysAdminPasswdErr, "{#ErrSysUserSysAdminPasswdErr}", nil)
	ErrSysUserBindRoleUnknow             = NewMpuapsErr(codeSysUserBindRoleUnknow, "{#ErrSysUserBindRoleUnknow}", nil)
	ErrSysUserUnBindRole                 = NewMpuapsErr(codeSysUserUnBindRole, "{#ErrSysUserUnBindRole}", nil)
	ErrSysUserPasswdErr                  = NewMpuapsErr(codeSysUserPasswdErr, "{#ErrSysUserPasswdErr}", nil)
	ErrSysUserNotLogin                   = NewMpuapsErr(codeSysUserNotLogin, "{#ErrSysUserNotLogin}", nil)
	ErrSysUserModifyDefSysUserRoleForbid = NewMpuapsErr(codeSysUserModifyDefSysUserRoleForbid, "{#ErrSysUserModifyDefSysUserRoleForbid}", nil)
	ErrSysUserPassInvalid                = NewMpuapsErr(codeSysUserPassInvalid, "{#ErrSysUserPassInvalid}", nil)
	ErrSysUserLockDefSysUserForbid       = NewMpuapsErr(codeSysUserLockDefSysUserForbid, "{#ErrSysUserLockDefSysUserForbid}", nil)
	ErrSysUserLoginNameInvalid           = NewMpuapsErr(codeSysUserLoginNameInvalid, "{#ErrSysUserLoginNameInvalid}", nil)
	ErrSysUserUpdateNewPassSameAsOldPass = NewMpuapsErr(codeSysUserUpdateNewPassSameAsOldPass, "{#ErrSysUserUpdateNewPassSameAsOldPass}", nil)
	ErrSysUserUpdatePassNotSameAsConfirm = NewMpuapsErr(codeSysUserUpdatePassNotSameAsConfirm, "{#ErrSysUserUpdatePassNotSameAsConfirm}", nil)
	ErrSysUserUpdatePassInputOldPassErr  = NewMpuapsErr(codeSysUserUpdatePassInputOldPassErr, "{#ErrSysUserUpdatePassInputOldPassErr}", nil)

	ErrSysUserLoginAuthFailed                              = NewMpuapsErr(codeSysUserLoginAuthFailed, "{#ErrSysUserLoginAuthFailed}", nil)
	ErrSysUserLoginCaptchaVerifyFailed                     = NewMpuapsErr(codeSysUserLoginCaptchaVerifyFailed, "{#ErrSysUserLoginCaptchaVerifyFailed}", nil)
	ErrSysUserLoginLock                                    = NewMpuapsErr(codeSysUserLoginLock, "{#ErrSysUserLoginLock}", nil)
	ErrSysUserCascadeLoginForbidRelationsUnknown           = NewMpuapsErr(codeSysUserCascadeLoginForbidRelationsUnknown, "{#ErrSysUserCascadeLoginForbidRelationsUnknown}", nil)
	ErrSysUserCascadeLoginForbidSameCascadeId              = NewMpuapsErr(codeSysUserCascadeLoginForbidSameCascadeId, "{#ErrSysUserCascadeLoginForbidSameCascadeId}", nil)
	ErrSysUserCascadeLoginForbidSubordinates               = NewMpuapsErr(codeSysUserCascadeLoginForbidSubordinates, "{#ErrSysUserCascadeLoginForbidSubordinates}", nil)
	ErrSysUserCascadeLoginNetErr                           = NewMpuapsErr(codeSysUserCascadeLoginNetErr, "{#ErrSysUserCascadeLoginNetErr}", nil)
	ErrSysUserCascadeLoginForbidHaveOtherParentCascadePlat = NewMpuapsErr(codeSysUserCascadeLoginForbidHaveOtherParentCascadePlat, "{#ErrSysUserCascadeLoginForbidHaveOtherParentCascadePlat}", nil)
	ErrSysUserLoginAlready                                 = NewMpuapsErr(codeSysUserLoginAlready, "{#ErrSysUserLoginAlready}", nil)
	ErrSysUserLoginPermissionDenied                        = NewMpuapsErr(codeSysUserLoginPermissionDenied, "{#ErrSysUserLoginPermissionDenied}", nil)
	ErrSysUserLoginBlackForbid                             = NewMpuapsErr(codeSysUserLoginBlackForbid, "{#ErrSysUserLoginBlackForbid}", nil)
	ErrSysUserLoginAccessConnectionsNumExceedMax           = NewMpuapsErr(codeSysUserLoginAccessConnectionsNumExceedMax, "{#ErrSysUserLoginAccessConnectionsNumExceedMax}", nil)
	ErrSysUserLoginFailed                                  = NewMpuapsErr(codeSysUserLoginFailed, "{#ErrSysUserLoginFailed}", nil)
)
