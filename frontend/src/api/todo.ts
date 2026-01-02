const API_URL = 'http://localhost:8080/todos';

// 一覧取得
export const fetchTodos = async () => {
  const response = await fetch(API_URL);
  return await response.json();
};

// 新規作成
export const createTodo = async (title: string) => {
  const response = await fetch('http://localhost:8080/todos', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ title }),
  })
  if (!response.ok) throw new Error('Create failed')
}

// 完了状態の更新
export const updateTodoStatus = async (id: number, is_completed: boolean): Promise<void> => {
  const response = await fetch(`${API_URL}/${id}`, {
    method: 'PATCH', // 部分更新なのでPATCH
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ is_completed }),
  });

  if (!response.ok) {
    throw new Error('ステータスの更新に失敗しました');
  }
};

// 削除
export const deleteTodo = async (id: number): Promise<void> => {
  const response = await fetch(`${API_URL}/${id}`, {
    method: 'DELETE',
  });
  if (!response.ok) {
    throw new Error('削除に失敗しました');
  }
};