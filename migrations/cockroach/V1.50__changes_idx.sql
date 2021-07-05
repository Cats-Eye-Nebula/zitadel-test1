BEGIN;

SET experimental_enable_hash_sharded_indexes=on;
CREATE INDEX changes_idx ON eventstore.events (aggregate_type, aggregate_id, creation_date) USING HASH WITH BUCKET_COUNT = 10;

COMMIT;
