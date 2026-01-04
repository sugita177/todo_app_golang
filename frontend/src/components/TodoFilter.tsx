import { type FilterType } from "../types/todo";

interface Props {
  currentFilter: FilterType;
  onFilterChange: (filter: FilterType) => void;
}

export const TodoFilter = ({ currentFilter, onFilterChange }: Props) => {
  const filters: { label: string; value: FilterType }[] = [
    { label: 'すべて', value: 'all' },
    { label: '未完了', value: 'active' },
    { label: '完了済', value: 'completed' },
  ];

  return (
    <div className="flex justify-center space-x-2 mb-8 p-1 bg-slate-100 rounded-xl w-fit mx-auto">
      {filters.map((f) => (
        <button
          key={f.value}
          onClick={() => onFilterChange(f.value)}
          className={`px-4 py-1.5 rounded-lg text-sm font-medium transition-all ${
            currentFilter === f.value
              ? 'bg-white text-blue-600 shadow-sm'
              : 'text-slate-500 hover:text-slate-700 hover:bg-slate-200/50'
          }`}
        >
          {f.label}
        </button>
      ))}
    </div>
  );
};