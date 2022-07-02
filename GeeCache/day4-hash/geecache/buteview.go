package geecache

/*

抽象出一个只读数据结构ByteView用来表示缓存值

*/
//这个结构体保持住了一个不可变的一系列字节的视图
type ByteView struct {
	//可以将string转换为字节数组
	b []byte
}

//这个方法返回视图的长度
func (v ByteView) Len() int {
	return len(v.b)
}

//字节切片返回，数据的副本作为一个字节切片
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

//这个方法返回作为一个字符串的数据，copy副本如果必要的话
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

/*

- ByteView 只有一个数据成员，b []byte，b 将会存储真实的缓存值。选择 byte 类型是为了能够支持任意的数据类型的存储，例如字符串、图片等。
- 实现 Len() int 方法，我们在 lru.Cache 的实现中，要求被缓存对象必须实现 Value 接口，即 Len() int 方法，返回其所占的内存大小。
- b 是只读的，使用 ByteSlice() 方法返回一个拷贝，防止缓存值被外部程序修改。

*/
