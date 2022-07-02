package lru

import (
	"container/list"
)

type Cache struct {
	//最大容量
	maxBytes int64
	//使用容量
	nbytes int64
	//双向链表
	ll *list.List
	//缓存map,value为结点指针类型，指向的结点存储各类型缓存数据
	cache map[string]*list.Element
	//可选,当entry节点被清除的时候执行
	OnEvicted func(key string, value Value)
}

//键值对，是双向链表结点的数据类型
type entry struct {
	key   string
	value Value
}

//Value使用Len方法去计算它所占用的字节数\
//为了通用性，我们允许值是实现了 Value 接口的任意类型，该接口只包含了一个方法 Len() int，用于返回值所占用的内存大小。

type Value interface {
	Len() int
}

/*
在这里我们直接使用 Go 语言标准库实现的双向链表list.List。
字典的定义是 map[string]*list.Element，键是字符串，值是双向链表中对应节点的指针。
maxBytes 是允许使用的最大内存，nbytes 是当前已使用的内存，OnEvicted 是某条记录被移除时的回调函数，可以为 nil。
键值对 entry 是双向链表节点的数据类型，在链表中仍保存每个值对应的 key 的好处在于，淘汰队首节点时，需要用 key 从字典中删除对应的映射。
为了通用性，我们允许值是实现了 Value 接口的任意类型，该接口只包含了一个方法 Len() int，用于返回值所占用的内存大小。
*/

/*
方便实例化Cache，实现New函数
*/

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {

	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

/*
查找功能有两步骤
- 第一次是从字典中找到对应的双向链表的结点
- 将该结点移动到队列尾部
*/

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		//ll集合的内置方法
		c.ll.MoveToFront(ele)
		//断言ele.Value是否是*entry类型，如果是的那就就转型并且返回，不是就会painc
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

/*
如果键对应的链表节点存在，则将对应节点移动到队尾，并返回查找到的值。
c.ll.MoveToFront(ele)，即将链表中的节点 ele 移动到队尾（双向链表作为队列，队首队尾是相对的，在这里约定 front 为队尾）
*/

/*
删除
这里的删除，实际上是缓存淘汰。即移除最近最少访问的节点（队首）
*/

/*
c.ll.Back() 取到队首节点，从链表中删除。
delete(c.cache, kv.key)，从字典中 c.cache 删除该节点的映射关系。
更新当前所用的内存 c.nbytes。
如果回调函数 OnEvicted 不为 nil，则调用回调函数。
*/

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())

		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}

}

/*
新增、修改，这一步也是比较关键的，因为上述方法中的断言之所以可以，是因为在这一步的正确操作

*/

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		//将一个元素push到front位置，同时返回这个节点元素指针
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	//为0的时候，表示内存的使用不受限
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

/*

最后为了测试方便,实现了Len方法来获取添加了多少条的数据
*/

func (c *Cache) Len() int {
	return c.ll.Len()
}
