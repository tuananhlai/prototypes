**The example is currently broken.**

- The replica has started replicating, but it can not keep up with the primary due to some inexplicable errors.

Run on the primary and take not of the source log file position.

```sql
SHOW MASTER STATUS;
```

Run on the replica.

```sql
CHANGE REPLICATION SOURCE TO
  SOURCE_HOST='mysql-primary',
  SOURCE_USER='replicator',
  SOURCE_PASSWORD='replica_password',
  SOURCE_LOG_FILE='mysql-bin.000001',
  SOURCE_LOG_POS=157;

-- Check if the replica succeeded in connecting to the primary and is replicating its data.
SHOW REPLICA STATUS;\G
```
