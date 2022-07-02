package geecache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

/*

在这个测试用例中，我们借助 GetterFunc 的类型转换，将一个匿名回调函数转换成了接口 f Getter。
调用该接口的方法 f.Get(key string)，实际上就是在调用匿名回调函数。

定义一个函数类型 F，并且实现接口 A 的方法，然后在这个方法中调用自己。这是 Go 语言中将其他函数（参数返回值定义与 F 一致）转换为接口 A 的常用技巧。

*/
func TestGetterFunc_Get(t *testing.T) {

	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}

}

func TestGroup_Get(t *testing.T) {

	//首先使用一个map模拟耗时的数据库
	var db = map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}

	//创建group实例，并且测试get方法
	loadCounts := make(map[string]int, len(db))
	gee := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}

				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {

		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatalf("failed to get value of Tom")
		}
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}
}

/*

在这个测试用例中，我们主要测试了 2 种情况
1）在缓存为空的情况下，能够通过回调函数获取到源数据。
2）在缓存已经存在的情况下，是否直接从缓存中获取，为了实现这一点，使用 loadCounts 统计某个键调用回调函数的次数，如果次数大于1，则表示调用了多次回调函数，
没有缓存。

*/
