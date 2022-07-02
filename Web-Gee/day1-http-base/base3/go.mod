module example

go 1.17

// 在 go.mod 中使用 replace 将 gee 指向 ./gee

require gee v0.0.0

replace gee => ./gee