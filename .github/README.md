# GitHub Actions Workflows

This directory contains GitHub Actions workflows for the Dora Block Explorer project.

## Available Workflows

### Update TOGAF Documentation

**File**: `.github/workflows/update-togaf-docs.yml`

This workflow automatically updates the TOGAF documentation in the wiki whenever the source file (`TOGAF/architecture_documentation.md`) is modified.

**Triggered on**:
- Push to `TOGAF/architecture_documentation.md`
- Manual trigger via workflow_dispatch

**Actions performed**:
1. Updates the wiki page with the latest content from the source file
2. Updates the last modified timestamp
3. Commits and pushes the changes to both the GitHub wiki and the local `docs/wiki` directory

## Manual Trigger

To manually trigger the TOGAF documentation update:

1. Go to the "Actions" tab in the GitHub repository
2. Select "Update TOGAF Documentation" workflow
3. Click "Run workflow"

## Requirements

For the wiki update to work, the repository must have the wiki feature enabled, and the GitHub Actions bot must have write access to the wiki.
