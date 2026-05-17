DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'queues' AND column_name = 'anonymous_chat'
  ) THEN
    ALTER TABLE queues RENAME COLUMN anonymous_chat TO skip_feature;
  END IF;
END $$;

ALTER TABLE queues ADD COLUMN IF NOT EXISTS skip_duration INT NOT NULL DEFAULT 10;