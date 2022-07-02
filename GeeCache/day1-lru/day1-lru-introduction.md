1 FIFO/LFU/LRU 算法简介

GeeCache 的缓存全部存储在内存中，内存是有限的，因此不可能无限制地添加数据。假定我们设置缓存能够使用的内存大小为 N，那么在某一个时间点，
添加了某一条缓存记录之后，占用内存超过了 N，这个时候就需要从缓存中移除一条或多条数据了。那移除谁呢？我们肯定希望尽可能移除“没用”的数据，
那如何判定数据“有用”还是“没用”呢？

1.1 FIFO(First In First Out)
先进先出，也就是淘汰缓存中最老(最早添加)的记录。FIFO 认为，最早添加的记录，其不再被使用的可能性比刚添加的可能性大。这种算法的实现也非常
简单，创建一个队列，新增记录添加到队尾，每次内存不够时，淘汰队首。但是很多场景下，部分记录虽然是最早添加但也最常被访问，而不得不因为呆的
时间太长而被淘汰。这类数据会被频繁地添加进缓存，又被淘汰出去，导致缓存命中率降低。

1.2 LFU(Least Frequently Used)
最少使用，也就是淘汰缓存中访问频率最低的记录。LFU 认为，如果数据过去被访问多次，那么将来被访问的频率也更高。LFU 的实现需要维护一个按照
访问次数排序的队列，每次访问，访问次数加1，队列重新排序，淘汰时选择访问次数最少的即可。LFU 算法的命中率是比较高的，但缺点也非常明显，维
护每个记录的访问次数，对内存的消耗是很高的；另外，如果数据的访问模式发生变化，LFU 需要较长的时间去适应，也就是说 LFU 算法受历史数据的
影响比较大。例如某个数据历史上访问次数奇高，但在某个时间点之后几乎不再被访问，但因为历史访问次数过高，而迟迟不能被淘汰。

1.3 LRU(Least Recently Used)
最近最少使用，相对于仅考虑时间因素的 FIFO 和仅考虑访问频率的 LFU，LRU 算法可以认为是相对平衡的一种淘汰算法。LRU 认为，如果数据最近
被访问过，那么将来被访问的概率也会更高。LRU 算法的实现非常简单，维护一个队列，如果某条记录被访问了，则移动到队尾，那么队首则是最近最少
访问的数据，淘汰该条记录即可。

2 LRU 算法实现
2.1 核心数据结构

https://geektutu.com/post/geecache-day1/lru.jpg

这张图很好地表示了 LRU 算法最核心的 2 个数据结构

绿色的是字典(map)，存储键和值的映射关系。这样根据某个键(key)查找对应的值(value)的复杂是O(1)，在字典中插入一条记录的复杂度也是O(1)。
红色的是双向链表(double linked list)实现的队列。将所有的值放到双向链表中，这样，当访问到某个值时，将其移动到队尾的复杂度是O(1)，在队尾新增一条记录以及删除一条记录的复杂度均为O(1)。
接下来我们创建一个包含字典和双向链表的结构体类型 Cache，方便实现后续的增删查改操作。


学习点
- list官方文档 https://golang.org/pkg/container/list/
- go语言中文网 https://studygolang.com/pkgdoc
- go语言标准库 https://studygolang.com/pkgdoc
- reflect反射中的DeepEqual方法  可参考：https://blog.csdn.net/a765717/article/details/112508290
  Go语言提供了运行时反射的内置支持实现，并允许程序借助反射包来操纵任意类型的对象。在Golang的reflect.DeepEqual()函数用于检查是否x和y是“deeply equal”与否。要访问此函数，需要在程序中导入反射包。

用法：
func DeepEqual(x, y interface{}) bool
参数：此函数采用两个参数，其值可以是任意类型，即x，y。

返回值：此函数返回布尔值。

以下示例说明了以上方法在Golang中的用法：

对于array、slice、map、struct等类型，当比较两个值是否相等时，是不能使用==运算符的。

func DeepEqual(x, y interface{}) bool
深度比较，反馈两个对象是否深等价。
用来判断两个值是否深度一致
使用reflect.DeepEqual来比较两个slice、struct、map是否相等
基本类型值等比较会使用==，当比较array、slice的成员、map映射的键值对、struct结构体的字段时，需要进行深入比对，比如map的键值对，对键只使用==，但值会继续往深层比对。

深等价

x和y同nil或同non-nil
x和y具有相同的长度
x和y指向同一个底层数组所初始化的实体对象

&x[0] == &y[0]

注意

一个non-nil的空切片和一个nil的切片并不是深等价的，比如[]byte{}和[]byte{nil}是非等价的。
numbers、bools、strings、channels使用==相等则是深等价的

func TestSlice(t *testing.T) {
m1 := map[string]int{"id": 1, "pid": 0}
m2 := map[string]int{"pid": 0, "id": 1}
//t.Log(m1 == m2)//invalid operation: m1 == m2 (map can only be compared to nil)

    //map变量只能和空（nil）比较
    //t.Log(m1 == nil) //false
    //t.Log(m2 != nil) //true

    t.Log(reflect.DeepEqual(m1, m2)) //true
}

