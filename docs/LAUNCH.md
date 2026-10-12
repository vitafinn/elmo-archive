# Release maintenance

Elmo Archive is a personally maintained fork of ElmoBeacon. Keep upstream history, copyright notices, and GPL v3 licensing. No feature roadmap or regular release schedule is planned.

## Publish an update

1. Review the source diff. Never commit databases, tokens, capture URLs, logs, or certificate private keys.
2. Run the checks and commit the reviewed changes. Push to main.
3. In GitHub Actions, run **Release** with a new version, such as `0.2.1`.
4. Review the draft's Windows ZIP, matching source ZIP, checksums, and notes. Publish when ready.

The workflow builds the app with the selected version and this repository's Releases link. It does not publish the draft automatically.

## Checks before release

- Back up the database with the app closed. Test the ZIP in a writable folder.
- Check capture, account confirmation, sync completion, and certificate removal. Confirm proxy settings are restored after capture and cancellation.
- Check copy costs, Custom/Reunion banner pity, outfit records, English, and Chinese.
- List untested regions honestly. Do not mark Haoplay paths tested without a live test.
- Check package contents and screenshots for private information.

Use the separate publishing checkout under `build/publish/elmo-archive`. It retains the upstream history. Do not use the earlier fresh-history snapshot or force-push over published commits.
