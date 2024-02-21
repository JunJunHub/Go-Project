package errcode

//[125000,125999] 显控业务功能预留错误代码
const (
	codeMpusrvAlarmToneFileNotExist = iota + codeMpusrvAlarmManageStart //125000
	codeMpusrvAlarmToneFileReadFailed

	codeMpusrvAlarmRecordNotExist
	codeMpusrvAlarmRecordIsHandled
)

var (
	ErrMpusrvAlarmToneFileNotExist   = NewMpuapsErr(codeMpusrvAlarmToneFileNotExist, "{#ErrMpusrvAlarmToneFileNotExist}", nil)
	ErrMpusrvAlarmToneFileReadFailed = NewMpuapsErr(codeMpusrvAlarmToneFileReadFailed, "{#ErrMpusrvAlarmToneFileReadFailed}", nil)

	ErrMpusrvAlarmRecordNotExist  = NewMpuapsErr(codeMpusrvAlarmRecordNotExist, "{#ErrMpusrvAlarmRecordNotExist}", nil)
	ErrMpusrvAlarmRecordIsHandled = NewMpuapsErr(codeMpusrvAlarmRecordIsHandled, "{#ErrMpusrvAlarmRecordIsHandled}", nil)
)
