package geecache

type PeerPicker interface {
	/*
		通过传入的key选择相应节点的PeerGetter
	*/
	PickPeer(key string) (peer PeerGetter, ok bool)
}

//对应于HTTP客户端
type PeerGetter interface {
	//从对应的group中查找key对应的缓存值
	Get(group string, key string) ([]byte, error)
}

/*

在这里，抽象出 2 个接口，PeerPicker 的 PickPeer() 方法用于根据传入的 key 选择相应节点 PeerGetter。

接口 PeerGetter 的 Get() 方法用于从对应 group 查找缓存值。PeerGetter 就对应于上述流程中的 HTTP 客户端。


*/

/*

在 GeeCache 第三天 中我们为 HTTPPool 实现了服务端功能，通信不仅需要服务端还需要客户端，因此，我们接下来要为 HTTPPool 实现客户端的功能。

首先创建具体的 HTTP 客户端类 httpGetter，实现 PeerGetter 接口。

*/
