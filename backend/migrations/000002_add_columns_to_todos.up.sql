-- カラムの追加
ALTER TABLE todos ADD COLUMN description TEXT;
ALTER TABLE todos ADD COLUMN priority VARCHAR(10) DEFAULT 'medium';
ALTER TABLE todos ADD COLUMN due_date TIMESTAMP WITH TIME ZONE;
ALTER TABLE todos ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

-- 優先度の値に制約を設ける（任意ですが推奨）
ALTER TABLE todos ADD CONSTRAINT chk_priority CHECK (priority IN ('low', 'medium', 'high'));