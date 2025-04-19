-- Drop all River Queue tables in reverse order
DROP TABLE IF EXISTS river_periodic_schedules;
DROP TABLE IF EXISTS river_unique_job_reservations;
DROP TABLE IF EXISTS river_periodic_job_executions;
DROP TABLE IF EXISTS river_leader_elections;
DROP TABLE IF EXISTS river_job_archival_jobs;
DROP TABLE IF EXISTS river_jobs;