-- Creation of block table

CREATE TABLE IF NOT EXISTS blk (
  block_num BIGINT NOT NULL PRIMARY KEY,
  block_hash bytea NOT NULL,
  parent_hash bytea,
  block_time BIGINT
);

CREATE INDEX idx_blk_hash ON blk(block_hash);


-- Creation of tx table
CREATE TABLE IF NOT EXISTS tx (
  tx_hash bytea NOT NULL PRIMARY KEY,
  block_num BIGINT NOT NULL,
  tx_from bytea NOT NULL,
  tx_to bytea,
  tx_nonce BIGINT,
  tx_data bytea,
  tx_value bytea,
  CONSTRAINT fk_bk
      FOREIGN KEY(block_num) 
	  REFERENCES blk(block_num)
);

CREATE INDEX idx_tx_blk ON tx(block_num);


CREATE TABLE IF NOT EXISTS job_summary (
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  lastest_blk_num BIGINT NOT NULL,
  job_info JSON,
  success BOOLEAN, 
  create_at BIGINT,
  update_at BIGINT
);

CREATE VIEW max_block as 
    SELECT
        MAX(block_num)
    FROM blk;


--for dev, we need to avoid scrape whole chain data
INSERT INTO blk (block_num, block_hash, parent_hash, block_time) VALUES (1810001,'0x666666666666666','0x777777777777',14561676122);