import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter, useLocation } from 'react-router-dom';
import { describe, it, expect } from 'vitest';
import { TodoFilter } from './TodoFilter';

// URLの状態を画面に出力するテスト用ヘルパー
const LocationTracker = () => {
  const location = useLocation();
  return <div data-testid="location">{location.search}</div>;
};

describe('TodoFilter (URL連動テスト)', () => {
  const renderWithRouter = (initialEntry = '/') => {
    return render(
      <MemoryRouter initialEntries={[initialEntry]}>
        <TodoFilter />
        <LocationTracker />
      </MemoryRouter>
    );
  };

  it('「未完了」をクリックしたとき、URLが ?filter=active に更新されること', () => {
    renderWithRouter('/');

    const activeBtn = screen.getByText('未完了');
    fireEvent.click(activeBtn);

    // useLocation の結果を保持している要素の中身をチェック
    const locationDisplay = screen.getByTestId('location');
    expect(locationDisplay.textContent).toBe('?filter=active');
  });

  it('「完了済」をクリックしたとき、URLが ?filter=completed に更新されること', () => {
    renderWithRouter('/?filter=active'); // 途中の状態から開始

    const completedBtn = screen.getByText('完了済');
    fireEvent.click(completedBtn);

    const locationDisplay = screen.getByTestId('location');
    expect(locationDisplay.textContent).toBe('?filter=completed');
  });
});