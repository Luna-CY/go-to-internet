package socket

// Stack Connection的栈结构
type Stack []*Connection

// Push 推入
func (s *Stack) Push(connection *Connection) {
    if !connection.IsClosed {
        *s = append(*s, connection)
    }
}

// Pop 弹出
func (s *Stack) Pop() *Connection {
    if s.IsEmpty() {
        return nil
    }

    res := (*s)[len(*s)-1]
    *s = (*s)[:len(*s)-1]

    return res
}

// IsEmpty 检查栈是否为空
func (s *Stack) IsEmpty() bool {
    return 0 == len(*s)
}
