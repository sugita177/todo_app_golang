import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { StatsPage } from './StatsPage';
import { useTodos } from './../contexts/TodoContext';

// useTodos をモック化
vi.mock('./../contexts/TodoContext', () => ({
  useTodos: vi.fn(),
}));

describe('StatsPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('読み込み中のときは「読み込み中...」が表示されること', () => {
    (useTodos as any).mockReturnValue({
      todos: [],
      loading: true,
    });

    render(<StatsPage />);
    expect(screen.getByText('読み込み中...')).toBeInTheDocument();
  });

  it('タスクの統計（件数と達成率）が正しく計算・表示されること', () => {
    const mockTodos = [
      { id: 1, title: 'Task 1', is_completed: true },
      { id: 2, title: 'Task 2', is_completed: true },
      { id: 3, title: 'Task 3', is_completed: false },
      { id: 4, title: 'Task 4', is_completed: false },
    ];

    (useTodos as any).mockReturnValue({
      todos: mockTodos,
      loading: false,
    });

    render(<StatsPage />);

    // 達成率の計算: (2 / 4) * 100 = 50%
    expect(screen.getByText('50')).toBeInTheDocument();
    
    // 完了済み: 2件
    const completedCard = screen.getByText('完了済み').closest('div');
    expect(completedCard).toHaveTextContent('2件');

    // 未完了: 2件
    const activeCard = screen.getByText('未完了').closest('div');
    expect(activeCard).toHaveTextContent('2件');
  });

  it('タスクが0件のとき、達成率が0%と表示されること', () => {
    (useTodos as any).mockReturnValue({
      todos: [],
      loading: false,
    });

    render(<StatsPage />);

    // 「全体の達成率」というテキストの近くにある「0」を特定する
    // 達成率の数値部分を特定（％の隣の要素など、より厳密に指定）
    expect(screen.getByText('全体の達成率')).toBeInTheDocument();
    
    // 「完了済み」や「未完了」のカード内の数値を特定
    // getAllByText を使い、合計3つの「0」が存在することを確認する
    const zeroElements = screen.getAllByText('0');
    expect(zeroElements).toHaveLength(3); // 達成率、完了数、未完了数

    // StatCard内の値を個別にチェックする
    const completedCard = screen.getByText('完了済み').closest('div');
    expect(completedCard).toHaveTextContent('0件');

    const activeCard = screen.getByText('未完了').closest('div');
    expect(activeCard).toHaveTextContent('0件');
  });

  it('プログレスバーに正しい幅（style）が適用されていること', () => {
    const mockTodos = [
      { id: 1, title: 'T1', is_completed: true },
      { id: 2, title: 'T2', is_completed: false },
      { id: 3, title: 'T3', is_completed: false },
    ]; // 33.333...% -> 33%

    (useTodos as any).mockReturnValue({
      todos: mockTodos,
      loading: false,
    });

    const { container } = render(<StatsPage />);
    
    // プログレスバーのdivを探して、styleを確認
    // 達成率33%のバーを探す
    const progressBar = container.querySelector('[style*="width: 33%"]');
    expect(progressBar).toBeInTheDocument();
  });
});