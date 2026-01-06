import { renderHook, waitFor, act } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { TodoProvider, useTodos } from './TodoContext';
import * as api from './../api/todo';

// api/todo.ts の関数をすべてモック化
vi.mock('./../api/todo');

describe('TodoContext', () => {
  const mockTodos = [
    { id: 1, title: 'Test Todo', is_completed: false }
  ];

  beforeEach(() => {
    vi.clearAllMocks();
  });

  const wrapper = ({ children }: { children: React.ReactNode }) => (
    <TodoProvider>{children}</TodoProvider>
  );

  it('初期化時に fetchTodos が呼ばれ、todos がセットされること', async () => {
    vi.mocked(api.fetchTodos).mockResolvedValue(mockTodos);

    const { result } = renderHook(() => useTodos(), { wrapper });

    // 最初に loading が true であることを確認
    expect(result.current.loading).toBe(true);

    // 非同期のデータ取得が終わるのを待つ
    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.todos).toEqual(mockTodos);
    expect(api.fetchTodos).toHaveBeenCalledTimes(1);
  });

  it('addTodo を呼ぶと、createTodo が実行されデータが再取得されること', async () => {
    vi.mocked(api.fetchTodos).mockResolvedValue(mockTodos);
    vi.mocked(api.createTodo).mockResolvedValue(undefined);

    const { result } = renderHook(() => useTodos(), { wrapper });
    await waitFor(() => expect(result.current.loading).toBe(false));

    // addTodo を実行
    await act(async () => {
      await result.current.addTodo('New Todo');
    });

    expect(api.createTodo).toHaveBeenCalledWith('New Todo');
    // refresh() が呼ばれるため、fetchTodos は初期化と合わせて計2回呼ばれる
    expect(api.fetchTodos).toHaveBeenCalledTimes(2);
  });

  it('toggleTodo を呼ぶと、ステートが即座に更新されること', async () => {
    vi.mocked(api.fetchTodos).mockResolvedValue(mockTodos);
    const { result } = renderHook(() => useTodos(), { wrapper });
    await waitFor(() => expect(result.current.loading).toBe(false));

    await act(async () => {
      await result.current.toggleTodo(1, false);
    });

    expect(api.updateTodoStatus).toHaveBeenCalledWith(1, true);
    // toggle は再取得せずステートを直接書き換える仕様なので、todos[0] が true になっているはず
    expect(result.current.todos[0].is_completed).toBe(true);
  });

  it('deleteTodo を呼ぶと、ステートから対象が削除されること', async () => {
    vi.mocked(api.fetchTodos).mockResolvedValue(mockTodos);
    const { result } = renderHook(() => useTodos(), { wrapper });
    await waitFor(() => expect(result.current.loading).toBe(false));

    await act(async () => {
      await result.current.deleteTodo(1);
    });

    expect(api.deleteTodo).toHaveBeenCalledWith(1);
    expect(result.current.todos).toHaveLength(0);
  });
});