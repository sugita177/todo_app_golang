import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { TodoListPage } from './TodoListPage';
import { useTodos } from './../contexts/TodoContext';

// useTodos カスタムフックをモック化する
vi.mock('./../contexts/TodoContext', () => ({
  useTodos: vi.fn(),
}));

describe('TodoListPage', () => {
  const mockTodos = [
    { id: 1, title: 'テストタスク1', is_completed: false, created_at: '' },
    { id: 2, title: 'テストタスク2', is_completed: true, created_at: '' },
  ];

  const mockContextValue = {
    todos: mockTodos,
    loading: false,
    addTodo: vi.fn(),
    toggleTodo: vi.fn(),
    deleteTodo: vi.fn(),
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('読み込み中のときは「読み込み中...」が表示されること', () => {
    (useTodos as any).mockReturnValue({
      ...mockContextValue,
      loading: true,
      todos: [],
    });

    render(
      <MemoryRouter>
        <TodoListPage />
      </MemoryRouter>
    );

    expect(screen.getByText('読み込み中...')).toBeInTheDocument();
  });

  it('タスクが存在する場合、リストが表示されること', () => {
    (useTodos as any).mockReturnValue(mockContextValue);

    render(
      <MemoryRouter>
        <TodoListPage />
      </MemoryRouter>
    );

    expect(screen.getByText('テストタスク1')).toBeInTheDocument();
    expect(screen.getByText('テストタスク2')).toBeInTheDocument();
  });

  it('タスクが空のとき、専用のメッセージが表示されること', () => {
    (useTodos as any).mockReturnValue({
      ...mockContextValue,
      todos: [],
    });

    render(
      <MemoryRouter>
        <TodoListPage />
      </MemoryRouter>
    );

    expect(screen.getByText('タスクがありません。')).toBeInTheDocument();
  });

  it('URLパラメータ（?filter=completed）に応じてフィルタリングされること', () => {
    (useTodos as any).mockReturnValue(mockContextValue);

    // 完了済みのみを表示するURLをシミュレート
    render(
      <MemoryRouter initialEntries={['/?filter=completed']}>
        <TodoListPage />
      </MemoryRouter>
    );

    // 完了済みのタスクは表示される
    expect(screen.getByText('テストタスク2')).toBeInTheDocument();
    // 未完了のタスクは表示されない
    expect(screen.queryByText('テストタスク1')).not.toBeInTheDocument();
  });
});