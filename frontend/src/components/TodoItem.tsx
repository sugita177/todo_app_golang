interface Todo {
  id: number
  title: string
  is_completed: boolean
}

interface Props {
  todo: Todo
  onDelete: (id: number) => Promise<void>
  onToggle: (id: number, currentStatus: boolean) => Promise<void>
}

export const TodoItem = ({ todo, onDelete, onToggle }: Props) => {
  return (
    <li style={{ padding: '10px', borderBottom: '1px solid #eee' }}>
      <input 
        type="checkbox" 
        checked={todo.is_completed} 
        readOnly 
        aria-label={`todo-status-${todo.id}`} // テストで見つけやすくする
        onChange={() => onToggle(todo.id, todo.is_completed)}
      />
      <span style={{ marginLeft: '10px' }}>{todo.title}</span>
      <button onClick={() => onDelete(todo.id)} style={{ color: 'red' }}>削除</button>
    </li>
  )
}