# How to Manage `.ctx/docs/`

1. **Read previous sync logs** — Check `.ctx/logs/` for previous `sync_docs_*.md` files. Read the most recent ones (up to 3). If logs exist, extract the last commit hash and run `git diff <that_hash>..HEAD --name-only` to see what files changed since the last sync. Focus the rest of your work on those files and the docs that cover them. You don't need to review the entire codebase every time — only what actually changed.

2. **Review integrity** — Read each doc (or only the affected ones if step 1 narrowed the scope) and verify that what it describes matches the current codebase. Fix or remove anything outdated.

3. **Restructure if needed** — Consider if any doc should be renamed, split into two, merged, or reorganized to better reflect the current project structure.

4. **Create/update docs** — Each doc should group ~8-16 important files, following `.ctx/templates/doc_template.md`.

5. **Update project overview** — Once all docs are accurate, create or update `project_overview.md` following `.ctx/templates/project_overview_template.md`.

6. **Write a sync log** — Create a new log file in `.ctx/logs/` named `sync_docs_<YYYY-MM-DD_HH-MM-SS>.md`. This log gives future runs context about what you did. Follow `.ctx/templates/sync_log_template.md`.
