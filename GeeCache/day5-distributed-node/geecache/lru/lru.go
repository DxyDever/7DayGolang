package lru

import (
	"container/list"
)

type Cache struct {
	maxBytes int64
	nbytes   int64
	ll       *list.List
	cache    map[string]*list.Element
	//可选,当entry节点被清楚的时候执行
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

//Value使用Len方法去计算它所占用的字节数
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
		c.ll.MoveToFront(ele)
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
新增、修改
*/

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

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
