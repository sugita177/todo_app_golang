import { useTodos } from '../contexts/TodoContext';
import { LoadingDelay } from '../components/LoadingDelay';
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from 'recharts';

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
  const completedCount = todos.filter((t) => t.is_completed).length;
  const activeCount = totalCount - completedCount;
  const completionRate = totalCount > 0 ? Math.round((completedCount / totalCount) * 100) : 0;

  // Recharts用のデータ構造
  const data = [
    { name: '完了済み', value: completedCount, color: '#2563eb' }, // blue-600
    { name: '未完了', value: activeCount, color: '#e2e8f0' },    // slate-200
  ];

  // 1. カスタムコンポーネントを定義
  const CustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      return (
        <div className="bg-white/90 backdrop-blur-sm p-2 shadow-lg rounded-lg border border-slate-100 text-xs text-slate-600">
          <p className="font-bold text-slate-800">{payload[0].name}</p>
          <p>{payload[0].value}件</p>
        </div>
      );
    }
    return null;
  };

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-slate-900 font-sans">統計データ</h2>

      <div className="bg-white p-6 rounded-2xl shadow-xl border border-slate-100">
        <div className="relative h-64 w-full">
          {/* グラフコンテナ */}
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={data}
                cx="50%"
                cy="50%"
                innerRadius={60}  // ドーナツの穴のサイズ
                outerRadius={80}
                paddingAngle={5}
                dataKey="value"
                stroke="none"
              >
                {data.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.color} />
                ))}
              </Pie>
              <Tooltip 
              content={<CustomTooltip />}
              // 座標を固定（親要素の右上付近）
              position={{ x: 260, y:20 }}
              // アニメーションをOFF
              isAnimationActive={false}
            />
            </PieChart>
          </ResponsiveContainer>
          
          {/* 中央に表示する達成率テキスト */}
          <div className="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
            <span className="text-3xl font-bold text-slate-800">{completionRate}%</span>
            <span className="text-xs text-slate-500 font-medium">達成度</span>
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4 mt-6">
          <StatCard label="完了済み" value={completedCount} color="text-blue-600" />
          <StatCard label="未完了" value={activeCount} color="text-slate-400" />
        </div>
      </div>
    </div>
  );
};

const StatCard = ({ label, value, color }: { label: string; value: number; color: string }) => (
  <div className="bg-slate-50 p-4 rounded-xl text-center">
    <p className="text-xs text-slate-500 font-medium mb-1">{label}</p>
    <p className={`text-2xl font-bold ${color}`}>
      {value}
      <span className="text-sm ml-1 text-slate-400 font-normal">件</span>
    </p>
  </div>
);