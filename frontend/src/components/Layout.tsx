import { Link, Outlet, useLocation } from 'react-router-dom';
import { ChartBarIcon, ListBulletIcon } from '@heroicons/react/24/outline';

export const Layout = () => {
  const location = useLocation();

  const navItems = [
    { name: 'タスク一覧', path: '/', icon: ListBulletIcon },
    { name: '統計データ', path: '/stats', icon: ChartBarIcon },
  ];

  return (
    <div className="min-h-screen bg-slate-50">
      {/* ナビゲーションバー */}
      <nav className="bg-white border-b border-slate-200 sticky top-0 z-10">
        <div className="max-w-md mx-auto px-4">
          <div className="flex justify-around">
            {navItems.map((item) => {
              const isActive = location.pathname === item.path;
              return (
                <Link
                  key={item.path}
                  to={item.path}
                  className={`flex flex-col items-center py-3 px-4 border-b-2 transition-colors ${
                    isActive
                      ? 'border-blue-600 text-blue-600'
                      : 'border-transparent text-slate-500 hover:text-slate-700'
                  }`}
                >
                  <item.icon className="h-6 w-6" />
                  <span className="text-xs mt-1 font-medium">{item.name}</span>
                </Link>
              );
            })}
          </div>
        </div>
      </nav>

      {/* 各ページの中身がここに表示される */}
      <main className="max-w-md mx-auto py-8 px-4">
        <Outlet />
      </main>
    </div>
  );
};