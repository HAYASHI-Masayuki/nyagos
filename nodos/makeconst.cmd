pushd "%~dp0"
go run ..\dos\importconst_run.go -p nodos ^
	COINIT_APARTMENTTHREADED ^
	COINIT_MULTITHREADED ^
	COINIT_DISABLE_OLE1DDE ^
	COINIT_SPEED_OVER_MEMORY ^
	FSCTL_SET_REPARSE_POINT
popd
