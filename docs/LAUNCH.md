# Gentle launch checklist

## Scope

Maintenance-focused: the current features are considered complete. No feature roadmap or scheduled updates are promised. The preview checks installation and existing behavior; future changes are for defects or game compatibility.

## Repository choice

Recommended: publish this as an independent continuation while preserving upstream commit history and GPL attribution. Do not delete the local `.git` folder for SEO. A new local history does not change GitHub’s fork relationship or make Google index a page.

The intended repository name is `elmo-archive`; the chosen owner is `vitafinn`. The release workflow sets `VITE_REPOSITORY_URL` from the actual GitHub repository, so no personal handle is hardcoded in the app. Local builds without that setting show the version without a release link.

Options:

- Keep the current repository and use GitHub Settings → General → Leave fork network if eligible. Read GitHub’s warnings first: detachment is permanent and repository metadata may be lost. See [GitHub’s detachment guide](https://docs.github.com/en/pull-requests/how-tos/work-with-forks/detaching-a-fork).
- Create an empty standalone repository using **New repository**, not **Fork**, then push a reviewed branch with its history. Keep this checkout and upstream remote for provenance. No need to delete the existing fork or local history.
- If you explicitly prefer fresh history, initialize a separate copy of the vetted source snapshot from `build/release/<version>/source`, preserving LICENSE and credits. Do not copy the entire working folder: it contains personal data. This is a source snapshot, not a history-preserving migration.

No remote repository has been created, detached, deleted, or pushed by the preparation work.

## Discoverability

GitHub repository search and Google Search are different systems. [GitHub repository search](https://docs.github.com/en/search-github/searching-on-github/searching-for-repositories) normally excludes forks unless the query includes `fork:true` or `fork:only`. This can explain GitHub discovery problems; it does not prove why Google omitted a result. Google explicitly says [indexing is not guaranteed](https://developers.google.com/search/docs/essentials/technical).

Suggested GitHub About description:

> Windows pull tracker for Girls’ Frontline 2: Exilium (GFL2/GF2). Pity, Custom/Reunion banners, outfit records, copy costs, and luck estimates. English & Chinese.

Suggested topics (GitHub Settings/About):

`girls-frontline-2`, `gfl2`, `gf2-exilium`, `gacha-tracker`, `pull-tracker`, `pity-tracker`, `darkwinter`, `haoplay`, `windows`, `golang`, `wails`, `vue`

The README now includes the game’s full name, common abbreviations, features, install steps, FAQ-style limitations, and upstream attribution. Avoid copying the old README unchanged or claiming official game affiliation. Add a current, UID-redacted screenshot after testing; the old `preview.png` is not used because it does not show this edition accurately.

After testing, publish one small pre-release with an understandable title and download links. Add an accurate profile/pinned-repository link and share the launch in relevant communities where permitted. Genuine links and useful documentation help people find it; keyword repetition and a fresh `.git` do not guarantee search ranking. A project page on a domain you control is optional later; it would let you use Google Search Console to monitor that page, not take ownership of github.com.

## Local checks already automated

- Frontend typecheck, build, localization coverage, and calculation tests.
- Go tests, including account-confirmation checks.
- A limited source preflight for private files and common credential patterns.
- Binary packaging from an explicit allowlist: executable, README, LICENSE, and documentation.
- Separate corresponding-source ZIP and SHA-256 checksums.
- Manual release workflow creates a **draft pre-release**, with no upstream cloud upload or automatic public release.

These checks do not prove that every secret pattern is absent or that every region works. Review changed source and generated archives before publishing. `.gitignore` prevents accidental new additions but does not remove anything already committed in older history.

## Checks for you before publication

1. Close other builds, back up your database, and run the candidate from its extracted ZIP in a writable folder. Do not run two copies against the same database.
2. On a fresh installation, review the certificate name/thumbprint, check Windows prompt visibility, cancel once, and confirm the same thumbprint is offered on retry. Confirm no capture starts before trust is installed. Existing valid certificates should not be replaced. Test manual UID capture and Darkwinter webshop capture. Confirm the captured-UID signal, final account confirmation, and successful sync.
3. Open an outfit Access Record page and confirm capture works; check outfit seasons after sync.
4. Check existing history, Custom/Reunion separate counters, a known lost-roll cost, and the copy locator. Check English/Chinese and large text.
5. Cancel capture and close the app; confirm normal networking and your previous proxy settings. Test the recovery directions if capture is interrupted.
6. Ask a Haoplay tester to verify their path before marking those regions tested. Share a candidate ZIP and matching source, never your app folder or database.
7. Review ZIP contents, license/source files, and checksums. Add a redacted current screenshot. Mark untested paths honestly in release notes.

## Publishing

Commit the reviewed source, including `frontend/package-lock.json`, new components, tests, and scripts. Do not use `git add .` without reviewing the resulting paths. Keep local session files and databases out of all commits.

Push to the chosen repository. In Actions, run **Preview release** with a version such as `0.2.0-preview.12`. Review the draft’s Windows ZIP, source ZIP, checksums, and notes, then publish it as a pre-release when the live checks are satisfactory. The workflow has not yet been run on GitHub.

To package a local build:

```powershell
# Set VITE_REPOSITORY_URL to the chosen GitHub repository URL before building.
wails build -ldflags '-X ElmoBeacon/handler.Version=0.2.0-preview.12' -o ElmoArchive.exe
./scripts/package-release.ps1 -Version 0.2.0-preview.12
```

Each version folder is new; the packaging script refuses to overwrite an existing candidate. The source archive contains current source files, including uncommitted changes, rather than blindly archiving an older HEAD. Use a clean committed checkout for the public release so its tag and source match.


Certificate cleanup checks: verify default removal after successful import, optional retention, manual removal during and after capture, cleanup failure messaging, and reconnect installation. Check that only the displayed certificate is removed, proxy settings restore, and an existing valid connection still syncs. Test unrelated HTTPS remains encrypted and Haoplay record hosts match the allowlist. Automated removal tests use an in-memory certificate store, not Windows root trust.

## Prepared publishing checkout

The publishing checkout is a separate clone under `build/publish/elmo-archive`. It retains the original reachable commit history and layers this edition’s changes on top. The original working folder remains unchanged by publishing operations. The earlier fresh-history snapshot is superseded; do not use it to initialize the public repository.

The clean source ZIP is still useful as corresponding source for a release. Source ZIPs intentionally omit Git metadata; that does not mean the GitHub repository should omit its history.

Create an empty GitHub repository named `vitafinn/elmo-archive` (without an initial README, license, or gitignore). After reviewing and committing the prepared changes in the publishing checkout:

```powershell
git status
git log --oneline -5
git push -u origin main
```

Do not force-push over an existing repository. Run the Preview release workflow on the committed main branch, review its draft, and publish when live checks are satisfactory. The workflow sets the app’s Releases link to its actual repository. Locally built candidates without that setting have no repository link.

Original author names, dates, messages, and commits remain intact. Use the maintainer’s chosen Git author identity only for the new changes. No remote has been created or pushed by this preparation.
