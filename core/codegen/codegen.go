package codegen

var globalStr = "declare i32 @printf(i8*, ...)\n"

func GetGlobals() string {
	return globalStr
}

func AppendToGlobals(str string) {
	globalStr += str
}

func GenPrintInt(s string) string {
	globalFormatStr := "@format = constant [3 x i8] c\"%d\\00\"\n"
	AppendToGlobals(globalFormatStr)
	res := "%formatStr = getelementptr [3 x i8], [3 x i8]* @format, i32 0, i32 0 ; Get a pointer to the format string\n"
	res += "call i32 (i8*, ...) @printf(i8* %formatStr, i32 " + s + ")\n"
	return res
}
