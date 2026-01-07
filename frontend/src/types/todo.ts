export interface Todo {
  id: number;
  title: string;
  description: string;
  is_completed: boolean;
  priority: 'low' | 'medium' | 'high';
  due_date: string | null;
  created_at: string;
}

export type FilterType = 'all' | 'active' | 'completed';