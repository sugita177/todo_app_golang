import { useSearchParams } from 'react-router-dom';
import { TodoInput } from '../components/TodoInput';
import { TodoItem } from '../components/TodoItem';
import { TodoFilter } from '../components/TodoFilter';
import { LoadingDelay } from '../components/LoadingDelay';
import { type FilterType } from '../types/todo';
import { useTodos } from '../contexts/TodoContext';

export const TodoListPage = () => {
  const { todos, loading, addTodo, toggleTodo, deleteTodo } = useTodos();
  const [searchParams] = useSearchParams();

  const filter = (searchParams.get('filter') as FilterType) || 'all';

  const filteredTodos = todos.filter((todo) => {
    if (filter === 'active') return !todo.is_completed;
    if (filter === 'completed') return todo.is_completed;
    return true;
  });

  if (loading) {
    return (
      <LoadingDelay delay={300}>
        <div className="text-center py-10 text-slate-400">読み込み中...</div>
      </LoadingDelay>
    );
  }

  return (
    <div className="max-w-md mx-auto">
      <h1 className="text-3xl font-bold text-gray-800 mb-8 text-center">
        My TODO List
      </h1>
      
      <div className="bg-white p-6 rounded-2xl shadow-xl">
        {/* ContextのaddTodoを直接渡す */}
        <TodoInput onAdd={addTodo} />

        <TodoFilter />
  
        <ul className="space-y-2">
          {filteredTodos.map(todo => (
            <TodoItem 
              key={todo.id} 
              todo={todo} 
              onDelete={deleteTodo} 
              onToggle={toggleTodo}
            />
          ))}
        </ul>
        
        {todos.length === 0 && (
          <p className="text-center text-gray-500 mt-6">タスクがありません。</p>
        )}
      </div>
    </div>
  );
}