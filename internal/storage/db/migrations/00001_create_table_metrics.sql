CREATE TABLE IF NOT EXISTS metrics (
                                       id           TEXT,
                                       mtype 	  TEXT,
                                       delta		   BIGINT,
                                       val        DOUBLE PRECISION,
                                       PRIMARY KEY (id, mtype)
    );