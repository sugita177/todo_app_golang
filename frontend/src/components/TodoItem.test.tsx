import { render, screen } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import { TodoItem } from './TodoItem'
import userEvent from '@testing-library/user-event'

describe('TodoItem', () => {
  const mockTodo = {
    id: 1,
    title: 'テストタスク',
    is_completed: false
  }

  // テスト全体で使えるダミーの関数を定義しておく
  const dummyOnDelete = vi.fn().mockResolvedValue(undefined)
  const dummyOnToggle = vi.fn().mockResolvedValue(undefined)

  it('TODOのタイトルが正しく表示されること', () => {
    render(<TodoItem todo={mockTodo} onDelete={dummyOnDelete} onToggle={dummyOnToggle} />)
    expect(screen.getByText('テストタスク')).toBeInTheDocument()
  })

  it('未完了のTODOの場合、チェックボックスがオフであること', () => {
    render(<TodoItem todo={mockTodo} onDelete={dummyOnDelete} onToggle={dummyOnToggle} />)
    const checkbox = screen.getByLabelText(`todo-status-${mockTodo.id}`)
    expect(checkbox).not.toBeChecked()
  })

  it('完了済みのTODOの場合、チェックボックスがオンであること', () => {
    const completedTodo = { ...mockTodo, is_completed: true }
    render(<TodoItem todo={completedTodo} onDelete={dummyOnDelete} onToggle={dummyOnToggle} />)
    const checkbox = screen.getByLabelText(`todo-status-${completedTodo.id}`)
    expect(checkbox).toBeChecked()
  })
})

describe('TodoItem 削除機能', () => {
  const mockTodo = { id: 123, title: '削除されるタスク', is_completed: false }

  it('削除ボタンをクリックした後、確認モーダルで「削除する」をクリックするとonDeleteが呼ばれること', async () => {
    const user = userEvent.setup()
    const onDeleteMock = vi.fn().mockResolvedValue(undefined)
    const dummyOnToggle = vi.fn().mockResolvedValue(undefined)

    render(<TodoItem todo={mockTodo} onDelete={onDeleteMock} onToggle={dummyOnToggle} />)

    // 1. リストにある「削除」ボタンをクリック
    const deleteButton = screen.getByRole('button', { name: '削除' })
    await user.click(deleteButton)

    // 2. モーダルが表示されていることを確認（「削除する」というボタンを探す）
    const confirmButton = screen.getByRole('button', { name: '削除する' })
    
    // 3. モーダル内の確定ボタンをクリック
    await user.click(confirmButton)

    // IDが正しく渡されているか検証
    expect(onDeleteMock).toHaveBeenCalledWith(123)
  });

  it('削除ボタンをクリックした後、キャンセルを押すとonDeleteが呼ばれないこと', async () => {
    const user = userEvent.setup()
    const onDeleteMock = vi.fn()
    
    render(<TodoItem todo={mockTodo} onDelete={onDeleteMock} onToggle={vi.fn()} />)

    // 削除ボタンクリック
    await user.click(screen.getByText('削除'))

    // キャンセルボタンクリック
    const cancelButton = screen.getByText('キャンセル')
    await user.click(cancelButton)

    // onDeleteが呼ばれていないことを確認
    expect(onDeleteMock).not.toHaveBeenCalled()
  });
})