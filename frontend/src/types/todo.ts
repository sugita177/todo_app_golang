export interface Todo {
  id: number;
  title: string;
  is_completed: boolean;
  created_at: string;
}

export type FilterType = 'all' | 'active' | 'completed';