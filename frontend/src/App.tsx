import { useState, useEffect, useCallback } from 'react'
import { TodoInput } from './components/TodoInput'
import { TodoItem } from './components/TodoItem'
import { fetchTodos as apiFetchTodos, updateTodoStatus, createTodo, deleteTodo } from './api/todo'

// 型定義（別のファイルに切り出してもOK）
interface Todo {
  id: number
  title: string
  is_completed: boolean
  created_at: string
}

function App() {
  const [todos, setTodos] = useState<Todo[]>([])

  // 一覧取得：apiFetchTodos を使用
  const loadTodos = useCallback(async () => {
    try {
      const data = await apiFetchTodos()
      setTodos(data || [])
    } catch (error) {
      console.error('取得失敗:', error)
    }
  }, [])

  useEffect(() => {
    loadTodos()
  }, [loadTodos])

  // 作成：createTodo を使用
  const handleCreateTodo = async (title: string) => {
    try {
      await createTodo(title)
      await loadTodos() // リストを再取得
    } catch (error) {
      console.error('作成失敗:', error)
      alert('作成に失敗しました')
    }
  }

  // 削除：deleteTodo を使用
  const handleDeleteTodo = async (id: number) => {
    try {
      await deleteTodo(id)
      // 再取得せず、フロントのステートから直接消すと動作が軽快になります
      setTodos(prev => prev.filter(t => t.id !== id))
    } catch (error) {
      console.error('削除失敗:', error)
    }
  }

  // 更新：updateTodoStatus を使用
  const handleToggleTodo = async (id: number, currentStatus: boolean) => {
    try {
      const nextStatus = !currentStatus
      await updateTodoStatus(id, nextStatus)

      // 画面上の表示を更新
      setTodos(prev =>
        prev.map(t => (t.id === id ? { ...t, is_completed: nextStatus } : t))
      )
    } catch (error) {
      alert('更新に失敗しました')
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 py-10 px-4">
      <div className="max-w-md mx-auto">
        <h1 className="text-3xl font-bold text-gray-800 mb-8 text-center">
          My TODO List
        </h1>
        
        <div className="bg-white p-6 rounded-2xl shadow-xl">
          {/* 入力コンポーネントに作成関数を渡す */}
          <TodoInput onAdd={handleCreateTodo} />
    
          {/* リスト表示（TodoItemをループ回す） */}
          <ul className="space-y-3">
            {todos.map((todo) => (
              <TodoItem key={todo.id} todo={todo} onDelete={handleDeleteTodo} onToggle={handleToggleTodo}/>
            ))}
          </ul>
          
          {todos.length === 0 && (
            <p className="text-center text-gray-500 mt-6">タスクがありません。追加してみましょう！</p>
          )}
        </div>
      </div>
    </div>
  )
}

export default App