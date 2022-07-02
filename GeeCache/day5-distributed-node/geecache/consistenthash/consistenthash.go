package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//函数类型Hash，采取依赖注入的方式，允许用于替换成自定义的Hash函数，也方便测试时候替换

type Hash func(data []byte) uint32

type Map struct {
	//Hash函数
	hash Hash
	//虚拟结点倍数
	replicas int
	//哈希环
	keys []int //sorted
	//虚拟结点与真实结点的映射表
	hashMap map[int]string
}

//构造函数 New() 允许自定义虚拟节点倍数和 Hash 函数。

func New(replicas int, fn Hash) *Map {

	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}

	if m.hash == nil {
		/*
			// ChecksumIEEE返回CRC-32校验和

			//使用IEEE多项式。

		*/
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

/*
下面来实现添加真实结点/机器的Add方法
*/

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

/*
Add 函数允许传入 0 或 多个真实节点的名称。
对每一个真实节点 key，对应创建 m.replicas 个虚拟节点，虚拟节点的名称是：strconv.Itoa(i) + key，即通过添加编号的方式区分不同虚拟节点。
使用 m.hash() 计算虚拟节点的哈希值，使用 append(m.keys, hash) 添加到环上。
在 hashMap 中增加虚拟节点和真实节点的映射关系。
最后一步，环上的哈希值排序。

*/

/*
最后一步，实现选择结点的get方法
*/

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))

	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	/*
		如果idx等于len(m.keys)表示经过search之后，没有找到第一次下标m.keys[i] >= hash的时候，
		表示hash值大于所有的key里面的元素，由于这是一个环，所以取m.keys[0]处的值，也就是哈希值。
	*/
	return m.hashMap[m.keys[idx%len(m.keys)]]
}

/*

选择节点就非常简单了，第一步，计算 key 的哈希值。
第二步，顺时针找到第一个匹配的虚拟节点的下标 idx，从 m.keys 中获取到对应的哈希值。如果 idx == len(m.keys)，说明应选择 m.keys[0]，因为 m.keys 是一个环状结构，所以用取余数的方式来处理这种情况。
第三步，通过 hashMap 映射得到真实的节点。
至此，整个一致性哈希算法就实现完成了。


*/
