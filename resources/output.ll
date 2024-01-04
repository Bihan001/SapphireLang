declare i32 @printf(i8*, ...)
@format = constant [3 x i8] c"%d\00"
define i32 @main() {
entry:
%0 = mul i32 2, 3
%1 = sub i32 20, %0
%2 = sdiv i32 10, 5
%3 = add i32 %1, %2
%formatStr = getelementptr [3 x i8], [3 x i8]* @format, i32 0, i32 0 ; Get a pointer to the format string
call i32 (i8*, ...) @printf(i8* %formatStr, i32 %3)
ret i32 0
}