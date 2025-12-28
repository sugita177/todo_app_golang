import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { TodoItem } from './TodoItem'

describe('TodoItem', () => {
  const mockTodo = {
    id: 1,
    title: 'テストタスク',
    is_completed: false
  }

  it('TODOのタイトルが正しく表示されること', () => {
    render(<TodoItem todo={mockTodo} />)
    expect(screen.getByText('テストタスク')).toBeInTheDocument()
  })

  it('未完了のTODOの場合、チェックボックスがオフであること', () => {
    render(<TodoItem todo={mockTodo} />)
    const checkbox = screen.getByLabelText(`todo-status-${mockTodo.id}`)
    expect(checkbox).not.toBeChecked()
  })

  it('完了済みのTODOの場合、チェックボックスがオンであること', () => {
    const completedTodo = { ...mockTodo, is_completed: true }
    render(<TodoItem todo={completedTodo} />)
    const checkbox = screen.getByLabelText(`todo-status-${completedTodo.id}`)
    expect(checkbox).toBeChecked()
  })
})