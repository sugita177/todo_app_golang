interface Todo {
  id: number
  title: string
  is_completed: boolean
}

interface Props {
  todo: Todo
  onDelete: (id: number) => Promise<void>
}

export const TodoItem = ({ todo, onDelete }: Props) => {
  return (
    <li style={{ padding: '10px', borderBottom: '1px solid #eee' }}>
      <input 
        type="checkbox" 
        checked={todo.is_completed} 
        readOnly 
        aria-label={`todo-status-${todo.id}`} // テストで見つけやすくする
      />
      <span style={{ marginLeft: '10px' }}>{todo.title}</span>
      <button onClick={() => onDelete(todo.id)} style={{ color: 'red' }}>削除</button>
    </li>
  )
}