package lru

import . "MyProject/Short_Url/models"

type Node struct {
	key        string
	value      string
	prev, next *Node
}

type LRUCache struct {
	cache map[string]*Node
	head  *Node
}

func Constructors(capacity int) *LRUCache { //初始化LRUCache
	head := &Node{"", "", nil, nil}
	node := head
	for i := 0; i < capacity-1; i++ {
		node.next = &Node{"", "", node, nil} //注意prev节点应该指向前一个node
		node = node.next
	}
	//首尾相连，构造环
	node.next = head
	head.prev = node

	return &LRUCache{
		head:  head,
		cache: make(map[string]*Node, capacity),
	}
}

func (this *LRUCache) MoveToFront(cur *Node) {
	if cur == this.head { //如果带移动的节点已经是头结点，则将头结点指向前一个节点，即指向最长时间没有更新过的那个节点
		this.head = this.head.prev
		return
	}

	//从链中取下待更新的节点
	cur.prev.next = cur.next
	cur.next.prev = cur.prev

	//将待更新节点放到头结点前面
	cur.next = this.head.next
	cur.next.prev = cur
	this.head.next = cur
	cur.prev = this.head
}

func (this *LRUCache) Get(key string) (string, error) {
	if node, ok := this.cache[key]; ok { //如果当前节点存在，则取出节点，并将节点位置更新至头结点的前面
		this.MoveToFront(node)
		return this.head.next.value, nil //返回的是head.next.value，因为最新的节点在头结点前面
	}

	//如果当前节点不存在，则到数据库中进行查询
	var urlCode UrlCode

	res, err := urlCode.GetByUrl(key)
	if err != nil { //如果数据库查询出错或者数据库中也没有存储，则直接返回
		return "", err
	}

	//查出的结果加入缓存中
	this.Put(res.Url, res.Code)
	return res.Code, nil
}

func (this *LRUCache) Put(key string, value string) {
	if node, ok := this.cache[key]; ok { //如果当前节点存在，则更新节点的value和位置
		node.value = value
		this.MoveToFront(node)
	} else { //如果节点不存在，则将节点插入到头结点的位置
		if this.head.value != "" { //如果头结点的值不等于 -1，说明该节点上已经有值，需先删除节点
			delete(this.cache, this.head.key)
		}

		//设置头结点的值为插入节点的值
		this.head.key = key
		this.head.value = value
		//更新cache
		this.cache[key] = this.head
		//头结点前移指向最久没有更新的那个节点
		this.head = this.head.prev
	}
}
