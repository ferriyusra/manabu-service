@echo off
REM Migration runner for Windows
REM Usage: migrate.bat <migration_file>
REM Example: migrate.bat 002_rename_users_uuid_constraint.sql

if "%1"=="" (
    echo Usage: migrate.bat ^<migration_file^>
    echo Example: migrate.bat 002_rename_users_uuid_constraint.sql
    exit /b 1
)

go run tools/migrate.go %1
