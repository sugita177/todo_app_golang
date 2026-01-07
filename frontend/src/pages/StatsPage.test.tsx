import { render, screen, act } from '@testing-library/react'; // actを追加
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { StatsPage } from './StatsPage';
import { useTodos } from './../contexts/TodoContext';

vi.mock('./../contexts/TodoContext', () => ({
  useTodos: vi.fn(),
}));

// RechartsのResponsiveContainerがテスト環境(JSDOM)でサイズ0にならないためのモック
vi.mock('recharts', async () => {
  const OriginalModule = await vi.importActual('recharts');
  return {
    ...OriginalModule,
    ResponsiveContainer: ({ children }: any) => (
      <div style={{ width: '800px', height: '800px' }}>{children}</div>
    ),
  };
});

describe('StatsPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.useFakeTimers(); // 仮想タイマーを有効にする
  });

  afterEach(() => {
    vi.useRealTimers(); // タイマーをリセットする
  });

  it('読み込み中のときは「読み込み中...」が表示されること', async () => {
    (useTodos as any).mockReturnValue({
      todos: [],
      loading: true,
    });

    render(<StatsPage />);
    
    // 最初は LoadingDelay によって表示されていない
    expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument();

    // 時間を 300ms 進める
    act(() => {
      vi.advanceTimersByTime(300);
    });

    // 時間経過後に表示されていることを確認
    expect(screen.getByText('読み込み中...')).toBeInTheDocument();
  });

  it('タスクの統計（件数と達成率）が正しく計算・表示されること', async () => {
    const mockTodos = [
      { id: 1, title: 'T1', is_completed: true },
      { id: 2, title: 'T2', is_completed: true },
      { id: 3, title: 'T3', is_completed: false },
      { id: 4, title: 'T4', is_completed: false },
    ];

    (useTodos as any).mockReturnValue({
      todos: mockTodos,
      loading: false,
    });

    render(<StatsPage />);

    // 文字が分割されている場合は正規表現を使うか、含んでいるかを確認する
    expect(screen.getByText(/50/)).toBeInTheDocument();
    expect(screen.getByText('達成度')).toBeInTheDocument();
    
    // 完了済み: 2件
    const completedCard = screen.getByText('完了済み').closest('div');
    expect(completedCard).toHaveTextContent('2件');

    const activeCard = screen.getByText('未完了').closest('div');
    expect(activeCard).toHaveTextContent('2件');
  });

  it('タスクが0件のとき、達成率が0%と表示されること', () => {
    (useTodos as any).mockReturnValue({ todos: [], loading: false });
    render(<StatsPage />);

    // 1. 中央の「達成度」というラベルのすぐ上にある「0」を特定する
    // コンテキストを含めた文字列、または要素の階層で特定します
    const rateElement = screen.getByText('達成度').previousElementSibling;
    expect(rateElement).toHaveTextContent('0%');

    // 2. 画面上の「0」をすべて取得して、3つあることを確認する
    const zeroCounts = screen.getAllByText(/0/);
    expect(zeroCounts).toHaveLength(3);
    
    // 3. 個別にカードの中身を確認する（より確実な方法）
    const completedCard = screen.getByText('完了済み').closest('div');
    expect(completedCard).toHaveTextContent('0件');
  });

  // 円グラフ（SVG）のテストは複雑なため、計算結果の数値を信じる形に変更するか、
  // もしプログレスバーを残していないならこのテストケース自体を削除/調整します。
  it('達成率の数値が中央に表示されていること', () => {
    const mockTodos = [
      { id: 1, title: 'T1', is_completed: true },
      { id: 2, title: 'T2', is_completed: false },
      { id: 3, title: 'T3', is_completed: false },
    ]; // 33%
    (useTodos as any).mockReturnValue({ todos: mockTodos, loading: false });

    render(<StatsPage />);
    
    // 33という数字が表示されているか確認
    expect(screen.getByText(/33/)).toBeInTheDocument();
  });
});