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

  it('TODOのタイトルが正しく表示されること', () => {
    render(<TodoItem todo={mockTodo} onDelete={dummyOnDelete} />)
    expect(screen.getByText('テストタスク')).toBeInTheDocument()
  })

  it('未完了のTODOの場合、チェックボックスがオフであること', () => {
    render(<TodoItem todo={mockTodo} onDelete={dummyOnDelete} />)
    const checkbox = screen.getByLabelText(`todo-status-${mockTodo.id}`)
    expect(checkbox).not.toBeChecked()
  })

  it('完了済みのTODOの場合、チェックボックスがオンであること', () => {
    const completedTodo = { ...mockTodo, is_completed: true }
    render(<TodoItem todo={completedTodo} onDelete={dummyOnDelete} />)
    const checkbox = screen.getByLabelText(`todo-status-${completedTodo.id}`)
    expect(checkbox).toBeChecked()
  })
})

describe('TodoItem 削除機能', () => {
  const mockTodo = { id: 123, title: '削除されるタスク', is_completed: false }

  it('削除ボタンをクリックすると、正しいIDでonDeleteが呼ばれること', async () => {
    const user = userEvent.setup()
    const onDeleteMock = vi.fn().mockResolvedValue(undefined)

    render(<TodoItem todo={mockTodo} onDelete={onDeleteMock} />)

    const deleteButton = screen.getByText('削除')
    await user.click(deleteButton)

    // IDが正しく渡されているか検証
    expect(onDeleteMock).toHaveBeenCalledWith(123)
  })
})