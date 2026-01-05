export const StatsPage = () => {
  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-slate-900">統計データ</h2>
      <div className="grid grid-cols-2 gap-4">
        <div className="bg-white p-4 rounded-xl shadow-sm border border-slate-100">
          <p className="text-sm text-slate-500">完了したタスク</p>
          <p className="text-3xl font-bold text-green-600">--</p>
        </div>
        <div className="bg-white p-4 rounded-xl shadow-sm border border-slate-100">
          <p className="text-sm text-slate-500">達成率</p>
          <p className="text-3xl font-bold text-blue-600">--%</p>
        </div>
      </div>
      <p className="text-center text-sm text-slate-400">
        ※TODOリストからデータを集計するロジックを次に追加します
      </p>
    </div>
  );
};