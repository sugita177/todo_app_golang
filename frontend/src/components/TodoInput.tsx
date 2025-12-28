import { useState } from 'react'

interface Props {
  onAdd: (title: string) => Promise<void>
}

export const TodoInput = ({ onAdd }: Props) => {
  const [title, setTitle] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!title.trim()) return
    await onAdd(title)
    setTitle('')
  }

  return (
    <form onSubmit={handleSubmit} style={{ marginBottom: '20px' }}>
      <input
        type="text"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="TODOを入力..."
        aria-label="todo-input"
        style={{ padding: '8px', width: '70%' }}
      />
      <button type="submit">追加</button>
    </form>
  )
}