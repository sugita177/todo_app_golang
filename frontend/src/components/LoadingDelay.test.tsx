import { render, screen, act } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { LoadingDelay } from './LoadingDelay';

describe('LoadingDelay', () => {
  beforeEach(() => {
    // 仮想タイマーを有効にする
    vi.useFakeTimers();
  });

  afterEach(() => {
    // タイマーを元に戻す
    vi.useRealTimers();
  });

  it('指定した時間が経過するまで children を表示しないこと', () => {
    render(
      <LoadingDelay delay={300}>
        <div data-testid="loading-content">読み込み中...</div>
      </LoadingDelay>
    );

    // 最初（0ms）は表示されていないはず
    expect(screen.queryByTestId('loading-content')).not.toBeInTheDocument();

    // 299ms 進めてもまだ表示されないはず
    act(() => {
      vi.advanceTimersByTime(299);
    });
    expect(screen.queryByTestId('loading-content')).not.toBeInTheDocument();

    // 300ms ちょうど、あるいは経過した後は表示されるはず
    act(() => {
      vi.advanceTimersByTime(1);
    });
    expect(screen.getByTestId('loading-content')).toBeInTheDocument();
  });

  it('アンマウント時にタイマーがクリーンアップされること', () => {
    // spy を使って clearTimeout が呼ばれるか監視する
    const clearTimeoutSpy = vi.spyOn(window, 'clearTimeout');
    
    const { unmount } = render(
      <LoadingDelay delay={300}>
        <div>Content</div>
      </LoadingDelay>
    );

    unmount();

    expect(clearTimeoutSpy).toHaveBeenCalled();
    clearTimeoutSpy.mockRestore();
  });
});