import { Routes, Route } from 'react-router-dom';
import { Layout } from './components/Layout';
import { TodoListPage } from './pages/TodoListPage';
import { StatsPage } from './pages/StatsPage';


function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        {/* Layoutの中の Outlet に表示される子ルート */}
        <Route index element={<TodoListPage />} />
        <Route path="stats" element={<StatsPage />} />
      </Route>
    </Routes>
  );
}

export default App;