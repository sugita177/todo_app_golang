import { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from 'react';
import { type Todo } from '../types/todo';
import { fetchTodos as apiFetchTodos, updateTodoStatus, createTodo, deleteTodo as apiDeleteTodo } from '../api/todo';

interface TodoContextType {
  todos: Todo[];
  loading: boolean;
  addTodo: (title: string) => Promise<void>;
  toggleTodo: (id: number, currentStatus: boolean) => Promise<void>;
  deleteTodo: (id: number) => Promise<void>;
  refresh: () => Promise<void>;
}

const TodoContext = createContext<TodoContextType | undefined>(undefined);

export const TodoProvider = ({ children }: { children: ReactNode }) => {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [loading, setLoading] = useState(true);

  // 一覧取得
  const refresh = useCallback(async () => {
    try {
      const data = await apiFetchTodos();
      setTodos(data || []);
    } catch (error) {
      console.error('取得失敗:', error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  // 作成
  const addTodo = async (title: string) => {
    await createTodo(title);
    await refresh(); // 最新状態を再取得
  };

  // 更新
  const toggleTodo = async (id: number, currentStatus: boolean) => {
    const nextStatus = !currentStatus;
    await updateTodoStatus(id, nextStatus);
    // 状態を直接更新することで、再取得なしで高速に反映
    setTodos(prev =>
      prev.map(t => (t.id === id ? { ...t, is_completed: nextStatus } : t))
    );
  };

  // 削除
  const deleteTodo = async (id: number) => {
    await apiDeleteTodo(id);
    setTodos(prev => prev.filter(t => t.id !== id));
  };

  return (
    <TodoContext.Provider value={{ todos, loading, addTodo, toggleTodo, deleteTodo, refresh }}>
      {children}
    </TodoContext.Provider>
  );
};

export const useTodos = () => {
  const context = useContext(TodoContext);
  if (!context) throw new Error('useTodos must be used within a TodoProvider');
  return context;
};