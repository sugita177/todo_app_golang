import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { TodoInput } from './TodoInput'

describe('TodoInput', () => {
  it('入力してボタンを押すと、onAddが呼ばれ、入力欄が空になること', async () => {
    // ユーザーのセットアップ
    const user = userEvent.setup()
    const onAddMock = vi.fn().mockResolvedValue(undefined)
    
    render(<TodoInput onAdd={onAddMock} />)

    const input = screen.getByLabelText('todo-input')
    const button = screen.getByText('追加')

    // 人間が操作するように、タイピングしてクリック
    await user.type(input, '新しいタスク')
    await user.click(button)

    // 検証
    expect(onAddMock).toHaveBeenCalledWith('新しいタスク')
    expect(input).toHaveValue('') // jest-domの便利なマッチャーを使用
  })
})