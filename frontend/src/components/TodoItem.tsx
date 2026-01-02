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
    <li className="flex items-center justify-between p-4 bg-white border border-gray-100 rounded-xl shadow-sm hover:shadow-md transition-shadow group">
      <div className="flex items-center gap-3">
        <input 
          type="checkbox" 
          checked={todo.is_completed} 
          readOnly 
          aria-label={`todo-status-${todo.id}`} // テストで見つけやすくする
          onChange={() => onToggle(todo.id, todo.is_completed)}
          className="w-5 h-5 cursor-pointer accent-blue-600"
        />
        <span className={`text-lg ${todo.is_completed ? 'line-through text-gray-400' : 'text-gray-700'}`}>
          {todo.title}
        </span>
      </div>
      <button
        onClick={() => onDelete(todo.id)}
        className="opacity-0 group-hover:opacity-100 p-2 text-red-500 hover:bg-red-50 rounded-lg transition-all cursor-pointer"
      >
        削除
      </button>
    </li>
  )
}