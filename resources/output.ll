define i32 @main() {
entry:
%0 = add i32 0, 20
%1 = add i32 0, 2
%2 = add i32 0, 3
%3 = mul i32 %1, %2
%4 = sub i32 %0, %3
%5 = add i32 0, 10
%6 = add i32 0, 5
%7 = sdiv i32 %5, %6
%8 = add i32 %4, %7
ret i32 %8
}