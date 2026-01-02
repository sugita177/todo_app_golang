import { useState, useEffect, useCallback } from 'react'
import { TodoInput } from './components/TodoInput'
import { TodoItem } from './components/TodoItem'
import { updateTodoStatus } from './api/todo'

// 型定義（別のファイルに切り出してもOK）
interface Todo {
  id: number
  title: string
  is_completed: boolean
  created_at: string
}

function App() {
  const [todos, setTodos] = useState<Todo[]>([])

  // APIからTODO一覧を取得する関数
  const fetchTodos = useCallback(async () => {
    try {
      const response = await fetch('http://localhost:8080/todos')
      if (!response.ok) throw new Error('Network response was not ok')
      const data = await response.json()
      setTodos(data || [])
    } catch (error) {
      console.error('取得失敗:', error)
    }
  }, [])

  // 初回レンダリング時に実行
  useEffect(() => {
    fetchTodos()
  }, [fetchTodos])

  // 新規TODOを作成する関数（TodoInputに渡す）
  const handleCreateTodo = async (title: string) => {
    try {
      const response = await fetch('http://localhost:8080/todos', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title }),
      })

      if (response.ok) {
        await fetchTodos() // リストを更新
      } else {
        alert('作成に失敗しました')
      }
    } catch (error) {
      console.error('作成失敗:', error)
    }
  }

  const handleDeleteTodo = async (id: number) => {
    try {
      const response = await fetch(`http://localhost:8080/todos/${id}`, {
        method: 'DELETE',
      })
      if (response.ok) {
        await fetchTodos() // 再取得してリストを更新
      }
    } catch (error) {
      console.error('削除失敗:', error)
    }
  }

  const handleToggleTodo = async (id: number, currentStatus: boolean) => {
    try {
      // 1. APIを叩いてDBを更新（現在の状態を反転させて送る）
      await updateTodoStatus(id, !currentStatus);

      // 2. 画面上の表示を更新（ステートの更新）
      setTodos(prev =>
        prev.map(t => (t.id === id ? { ...t, is_completed: !currentStatus } : t))
      );
    } catch (error) {
      alert('更新に失敗しました');
    }
  };

  return (
    <div style={{ padding: '20px', maxWidth: '400px', margin: '0 auto' }}>
      <h1>My TODO List</h1>
      
      {/* 入力コンポーネントに作成関数を渡す */}
      <TodoInput onAdd={handleCreateTodo} />

      {/* リスト表示（TodoItemをループ回す） */}
      <ul style={{ listStyle: 'none', padding: 0 }}>
        {todos.map((todo) => (
          <TodoItem key={todo.id} todo={todo} onDelete={handleDeleteTodo} onToggle={handleToggleTodo}/>
        ))}
      </ul>
      
      {todos.length === 0 && <p>タスクがありません。追加してみましょう！</p>}
    </div>
  )
}

export default App