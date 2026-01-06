import { useTodos } from '../contexts/TodoContext';
import { LoadingDelay } from '../components/LoadingDelay';

export const StatsPage = () => {
  const { todos, loading } = useTodos();

  if (loading) {
    return (
      <LoadingDelay delay={300}>
        <div className="text-center py-10 text-slate-400">読み込み中...</div>
      </LoadingDelay>
    );
  }

  const totalCount = todos.length;
  const completedCount = todos.filter(t => t.is_completed).length;
  const activeCount = totalCount - completedCount;
  const completionRate = totalCount > 0 ? Math.round((completedCount / totalCount) * 100) : 0;

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-slate-900 font-sans">統計データ</h2>
      
      <div className="grid grid-cols-1 gap-4">
        {/* 大きな進捗カード */}
        <div className="bg-linear-to-br from-blue-600 to-indigo-700 p-6 rounded-2xl text-white shadow-lg">
          <p className="text-blue-100 text-sm font-medium">全体の達成率</p>
          <div className="flex items-end space-x-2 mt-1">
            <span className="text-5xl font-bold">{completionRate}</span>
            <span className="text-xl mb-1">%</span>
          </div>
          {/* プログレスバー */}
          <div className="w-full bg-white/20 h-2 rounded-full mt-4 overflow-hidden">
            <div 
              className="bg-white h-full transition-all duration-1000" 
              style={{ width: `${completionRate}%` }}
            />
          </div>
        </div>

        {/* 詳細グリッド */}
        <div className="grid grid-cols-2 gap-4">
          <StatCard label="完了済み" value={completedCount} color="text-green-600" />
          <StatCard label="未完了" value={activeCount} color="text-orange-500" />
        </div>
      </div>
    </div>
  );
};

const StatCard = ({ label, value, color }: { label: string; value: number, color: string }) => (
  <div className="bg-white p-4 rounded-xl border border-slate-100 shadow-sm">
    <p className="text-xs text-slate-500 font-medium mb-1">{label}</p>
    <p className={`text-2xl font-bold ${color}`}>{value}<span className="text-sm ml-1 text-slate-400 font-normal">件</span></p>
  </div>
);