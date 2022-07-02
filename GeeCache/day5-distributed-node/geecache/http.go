package geecache

import (
	"fmt"
	"geecache/consistenthash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

//const defaultBasePath = "/_geecache/"

//type HTTPPool struct {
//	self     string
//	basePath string
//}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

/*
HTTPPool 只有 2 个参数，一个是 self，用来记录自己的地址，包括主机名/IP 和端口。
另一个是 basePath，作为节点间通讯地址的前缀，默认是 /_geecache/，那么 http://example.com/_geecache/ 开头的请求，就用于节点间的访问。因为一个主机上还可能承载其他的服务，加一段 Path 是一个好习惯。比如，大部分网站的 API 接口，一般以 /api 作为前缀。

*/

/*
下面来实现最为核心的方法，serverhttp方法
*/

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

//ServerHttp 处理所有的http请求
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}

	p.Log("%s %s", r.Method, r.URL.Path)

	// /<basepath>/<groupname>/<key> required
	/*
		将切片拆分为以sep分隔的子字符串，并返回
		分隔符之间的子字符串。
		计数决定返回的子字符串的数量:
		n > 0:最多n个子字符串;最后一个子串将是未拆分的余数。
		n == 0:结果为nil(零子字符串)
		n < 0:所有子字符串
	*/
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)

	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)

	if group == nil {
		http.Error(w, "no such group"+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())

}

/*

ServeHTTP 的实现逻辑是比较简单的，首先判断访问路径的前缀是否是 basePath，不是返回错误。
我们约定访问路径格式为 /<basepath>/<groupname>/<key>，通过 groupname 得到 group 实例，再使用 group.Get(key) 获取缓存数据。
最终使用 w.Write() 将缓存值作为 httpResponse 的 body 返回。
到这里，HTTP 服务端已经完整地实现了。接下来，我们将在单机上启动 HTTP 服务，使用 curl 进行测试。

*/

/*
创建具体的 HTTP 客户端类 httpGetter，实现 PeerGetter 接口。
*/
//HTTP客户端
type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		/*
			QueryEscape 对字符串进行转义，以便可以将其安全地放置在 URL 查询中。

			func QueryEscape(s string) string
			QueryEscape函数对s进行转码使之可以安全的用在URL查询里。

			func QueryUnescape(s string) (string, error)
			QueryUnescape函数用于将QueryEscape转码的字符串还原。它会把%AB改为字节0xAB，将’+’改为’ ‘。如果有某个%后面未跟两个十六进制数字，本函数会返回错误。

		*/
		url.QueryEscape(group),
		url.QueryEscape(key),
	)

	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned : %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}

/*
这个是什么意思呢？
确保这个类型实现了这个接口 如果没有实现会报错的
*/
var _ PeerGetter = (*httpGetter)(nil)

/*
baseURL 表示将要访问的远程节点的地址，例如 http://example.com/_geecache/。
使用 http.Get() 方式获取返回值，并转换为 []bytes 类型。

*/

/*
为HTTPPool添加节点的选择功能
*/

const (
	defaultBasePath = "/_geecache/"
	defaultReplicas = 50
)

type HTTPPool struct {
	self     string
	basePath string
	mu       sync.Mutex
	peers    *consistenthash.Map
	//string表示服务端标识，HttpGetter表示是对应的客户端
	httpGetter map[string]*httpGetter
}

/*
新增成员变量 peers，类型是一致性哈希算法的 Map，用来根据具体的 key 选择节点。
新增成员变量 httpGetters，映射远程节点与对应的 httpGetter。每一个远程节点对应一个 httpGetter，因为 httpGetter 与远程节点的地址 baseURL 有关。
*/

/*
第三步，实现PeerPicker接口
*/

func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetter = make(map[string]*httpGetter, len(peers))

	for _, peer := range peers {
		p.httpGetter[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetter[peer], true
	}
	return nil, false
}

/*
确保这个类型实现了这个接口 如果没有实现会报错的
*/
var _ PeerPicker = (*HTTPPool)(nil)

/*

Set() 方法实例化了一致性哈希算法，并且添加了传入的节点。
并为每一个节点创建了一个 HTTP 客户端 httpGetter。
PickerPeer() 包装了一致性哈希算法的 Get() 方法，根据具体的 key，选择节点，返回节点对应的 HTTP 客户端。
至此，HTTPPool 既具备了提供 HTTP 服务的能力，也具备了根据具体的 key，创建 HTTP 客户端从远程节点获取缓存值的能力。

*/

/*

下面就是在geecache中实现主流程

*/
