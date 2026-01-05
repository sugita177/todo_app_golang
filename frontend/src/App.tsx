import { Routes, Route } from 'react-router-dom';
import { Layout } from './components/Layout';
import { TodoListPage } from './pages/TodoListPage';
import { StatsPage } from './pages/StatsPage';
import { TodoProvider } from './contexts/TodoContext';


function App() {
  return (
    <TodoProvider>
      <Routes>
        <Route path="/" element={<Layout />}>
          {/* Layoutの中の Outlet に表示される子ルート */}
          <Route index element={<TodoListPage />} />
          <Route path="stats" element={<StatsPage />} />
        </Route>
      </Routes>
    </TodoProvider>
  );
}

export default App;