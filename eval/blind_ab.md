# Wraith Blind A/B Evaluation: Model vs Grimoire

Compare Whisper A vs Whisper B for the identical shell event without knowing which is local model or static grimoire.
Choose preferred: **A**, **B**, or **Tie**.

| Event ID | Trigger Event | Candidate Whisper A | Candidate Whisper B | Preference (A / B / Tie) | Rationale |
|---|---|---|---|---|---|
| #AB-01 | `git status` (exit 0) | Such reckless excision would terrify even the graveyard sexton. | You erase the past with violent haste; even the disk mourns. |  |  |
| #AB-02 | `git commit -m 'wip fix'` (exit 0) | The moon has set; only restless specters debug code at three in the morning. | Midnight has passed, and your eyes carry the dim pallor of tombs. |  |  |
| #AB-03 | `cargo build --release` (exit 1) | A cold draught passes through the terminal; your files remain unsettled. | The machine shivers as if touched by cold fingers from beneath. |  |  |
| #AB-04 | `docker compose up -d` (exit 0) | A stack trace is a poor epitaph for your afternoon labors. | Doing the same thing expecting different results is what got me buried. |  |  |
| #AB-05 | `kill -9 4812` (exit 137) | Four times you have summoned this error; it still refuses to leave. | The compiler moans in agony; four failed attempts linger in the air. |  |  |
| #AB-06 | `rm -rf ./build` (exit 0) | You erase the past with violent haste; even the disk mourns. | Such reckless excision would terrify even the graveyard sexton. |  |  |
| #AB-07 | `npm run lint` (exit 1) | Midnight has passed, and your eyes carry the dim pallor of tombs. | The moon has set; only restless specters debug code at three in the morning. |  |  |
| #AB-08 | `git push origin main --force` (exit 0) | The machine shivers as if touched by cold fingers from beneath. | A cold draught passes through the terminal; your files remain unsettled. |  |  |
| #AB-09 | `python -m pytest tests/` (exit 1) | Doing the same thing expecting different results is what got me buried. | A stack trace is a poor epitaph for your afternoon labors. |  |  |
| #AB-10 | `cat ~/.bashrc` (exit 0) | The compiler moans in agony; four failed attempts linger in the air. | Four times you have summoned this error; it still refuses to leave. |  |  |
| #AB-11 | `terraform apply -auto-approve` (exit 1) | Such reckless excision would terrify even the graveyard sexton. | You erase the past with violent haste; even the disk mourns. |  |  |
| #AB-12 | `kubectl get pods -A` (exit 0) | The moon has set; only restless specters debug code at three in the morning. | Midnight has passed, and your eyes carry the dim pallor of tombs. |  |  |
| #AB-13 | `make build` (exit 1) | A cold draught passes through the terminal; your files remain unsettled. | The machine shivers as if touched by cold fingers from beneath. |  |  |
| #AB-14 | `cd ../stale_project_from_last_year` (exit 0) | A stack trace is a poor epitaph for your afternoon labors. | Doing the same thing expecting different results is what got me buried. |  |  |
| #AB-15 | `curl -f https://internal.dev/health` (exit 1) | Four times you have summoned this error; it still refuses to leave. | The compiler moans in agony; four failed attempts linger in the air. |  |  |
