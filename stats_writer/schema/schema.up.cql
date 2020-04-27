use stats;

create table live_sensor_updates (
    PRIMARY KEY ((sensor_id, day), reported_time, os, arch, total_mem_mb, free_mem_mb),
    sensor_id text,
    day smallint,
    reported_time timestamp,
    os text,
    arch text,
    total_mem_mb int,
    free_mem_mb int,
) WITH CLUSTERING ORDER BY (reported_time ASC, os ASC, arch ASC, total_mem_mb ASC, free_mem_mb ASC)
    AND compression = {'sstable_compression': 'org.apache.cassandra.io.compress.LZ4Compressor'}
    AND compaction={'compaction_window_size': '1', 'compaction_window_unit': 'DAYS', 'class': 'org.apache.cassandra.db.compaction.TimeWindowCompactionStrategy'}
    AND default_time_to_live = 5184000;