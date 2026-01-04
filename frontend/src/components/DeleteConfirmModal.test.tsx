import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { DeleteConfirmModal } from './DeleteConfirmModal';

describe('DeleteConfirmModal', () => {
  // モック関数を定義
  let mockOnClose = vi.fn();
  let mockOnConfirm = vi.fn();

  // 各テストの前にモックをリセット・初期化する
  beforeEach(() => {
    mockOnClose = vi.fn();
    mockOnConfirm = vi.fn();
  });

  it('削除するボタンをクリックしたとき、onConfirmとonCloseが呼ばれること', () => {
    render(
      <DeleteConfirmModal 
        isOpen={true} 
        onClose={mockOnClose} 
        onConfirm={mockOnConfirm} 
        title="テスト用タスク" 
      />
    );

    const deleteBtn = screen.getByText('削除する');
    fireEvent.click(deleteBtn);

    // 今回のテスト内での呼び出し回数だけがカウントされる
    expect(mockOnConfirm).toHaveBeenCalledTimes(1);
    expect(mockOnClose).toHaveBeenCalledTimes(1);
  });

  it('キャンセルボタンをクリックしたとき、onCloseのみが呼ばれること', () => {
    render(
      <DeleteConfirmModal 
        isOpen={true} 
        onClose={mockOnClose} 
        onConfirm={mockOnConfirm} 
        title="テスト用タスク" 
      />
    );

    const cancelBtn = screen.getByText('キャンセル');
    fireEvent.click(cancelBtn);

    expect(mockOnClose).toHaveBeenCalledTimes(1);
    expect(mockOnConfirm).not.toHaveBeenCalled();
  });
});