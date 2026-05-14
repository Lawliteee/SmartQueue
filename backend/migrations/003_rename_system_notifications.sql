DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'queues' AND column_name = 'system_notifications'
  ) THEN
    ALTER TABLE queues RENAME COLUMN system_notifications TO im_free_feature;
  END IF;
END $$;