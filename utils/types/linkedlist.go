package types

// Data Linked List.
type LinkedList struct {
	Value   interface{}
	NextKey string
	Next    *LinkedList
}

// Append LinkedList.
func AddLinkedList(p *LinkedList, c *LinkedList) {
	if p.Next == nil {
		p.Next = c
	} else {
		AddLinkedList(p.Next, c)
	}
}

// Reassignment value.
func (ll *LinkedList) Reassignment() interface{} {
	if ll.Next != nil {
		switch GetType(ll.Value) {
		case TAnyMap:
			ll.Value.(map[string]interface{})[ll.NextKey] = ll.Next.Reassignment()
		case TStringMap:
			ll.Value = ToType(ll.Value, TAnyMap)
			ll.Value.(map[string]interface{})[ll.NextKey] = ll.Next.Reassignment()
		case TAnySlice:
			intKey := ToType(ll.NextKey, TInt)
			if intKey.(int) < len(ll.Value.([]interface{})) {
				ll.Value.([]interface{})[intKey.(int)] = ll.Next.Reassignment()
			} else {
				ll.Value = append(ll.Value.([]interface{}), ll.Next.Reassignment())
			}
		case TStringSlice:
			ll.Value = ToType(ll.Value, TAnySlice)
			intKey := ToType(ll.NextKey, TInt)
			if intKey.(int) < len(ll.Value.([]interface{})) {
				ll.Value.([]interface{})[intKey.(int)] = ll.Next.Reassignment()
			} else {
				ll.Value = append(ll.Value.([]interface{}), ll.Next.Reassignment())
			}
		}
	}
	return ll.Value
}
