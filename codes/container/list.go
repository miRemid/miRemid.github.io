package container

type LinkedList[T any] struct {
	root   Node[T]
	length int
}

type Node[T any] struct {
	next    *Node[T]
	prev    *Node[T]
	Element T
}

func NewLinkedList[T any]() *LinkedList[T] {
	return new(LinkedList[T]).Init()
}

func (list *LinkedList[T]) Init() *LinkedList[T] {
	list.root.next = &list.root
	list.root.prev = &list.root
	list.length = 0
	return list
}

func (list *LinkedList[T]) lazyInit() {
	if list.root.next == nil {
		list.Init()
	}
}

func (list *LinkedList[T]) Len() int {
	return list.length
}

func (list *LinkedList[T]) insert(at, e *Node[T]) {
	e.prev = at
	e.next = at.next
	e.next.prev = e
	at.next = e
	list.length++
}

// OtherMoveToBackList 将另一个Link插入到目前的链表后面
func (list *LinkedList[T]) OtherMoveToBackList(other *LinkedList[T]) *LinkedList[T] {
	list.lazyInit()

	if list == other || other.length == 0 {
		return list
	}

	list.length += other.length

	tail := list.root.prev
	otherHead := other.root.next
	otherTail := other.root.prev

	tail.next = otherHead
	otherHead.prev = tail

	otherTail.next = &list.root
	list.root.prev = otherTail

	other.Init()
	return list
}

func (list *LinkedList[T]) OtherMoveToFrontList(other *LinkedList[T]) *LinkedList[T] {
	list.lazyInit()
	if list == other || other.length == 0 {
		return list
	}
	list.length += other.length
	head := list.root.next
	otherHead, otherTail := other.root.next, other.root.prev
	otherTail.next = head
	head.prev = otherTail
	list.root.next = otherHead
	otherHead.prev = &list.root
	other.Init()
	return list
}

func (list *LinkedList[T]) MoveToFront(node *Node[T]) {
	list.lazyInit()

}

type DLinkNode struct {
	Key, Val int
	Prev     *DLinkNode
	Next     *DLinkNode
}

type DLinkList struct {
	length, cap int
	head, tail  *DLinkNode
	index       map[int]*DLinkNode
}

func NewDLinkList(cap int) *DLinkList {
	head := new(DLinkNode)
	tail := new(DLinkNode)
	head.Next = tail
	tail.Prev = head
	return &DLinkList{
		length: 0,
		cap:    cap,
		head:   head,
		tail:   tail,
		index:  make(map[int]*DLinkNode, 0),
	}
}

func (d *DLinkList) Insert(key, val int) {
	if _, ok := d.index[key]; !ok {
		node := &DLinkNode{
			Key: key,
			Val: val,
		}
		d.index[key] = node
		d.length++
		d.addToHead(node)
		if d.length > d.cap {
			removed := d.removeTail()
			delete(d.index, removed.Key)
			d.length--
		}
	} else {
		node := d.index[key]
		node.Val = val
		d.moveToFront(node)
	}
}

func (d *DLinkList) Get(key int) int {
	node, ok := d.index[key]
	if !ok {
		return -1
	}
	d.moveToFront(node)
	return node.Val
}

func (d *DLinkList) addToHead(node *DLinkNode) {
	node.Next = d.head.Next
	d.head.Next.Prev = node
	node.Prev = d.head
	d.head.Next = node
}

func (d *DLinkList) removeNode(node *DLinkNode) {
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
}

func (d *DLinkList) moveToFront(node *DLinkNode) {
	d.removeNode(node)
	d.addToHead(node)
}

func (d *DLinkList) removeTail() *DLinkNode {
	node := d.tail.Prev
	d.removeNode(node)
	return node
}
