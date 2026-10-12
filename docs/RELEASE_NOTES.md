# Elmo Archive preview

A Windows pull tracker for Girls’ Frontline 2: Exilium, forked from [MatchaCabin/ElmoBeacon](https://github.com/MatchaCabin/ElmoBeacon).

The feature set is complete for this edition. Future maintenance is limited to bug fixes and compatibility changes when needed; no feature roadmap or regular update schedule is planned.

## Certificate controls

Certificate trust is removed after a successful import by default. Manual removal is available in Import records, and stops active capture first. Existing valid connections can sync without reinstalling. HTTPS capture is limited to explicit hosts and verifies upstream TLS. Removal errors are shown separately from import results.

## Included

- Custom and Reunion procurement with separate banner pity.
- Outfit access records by season.
- Copy costs with linked losses, signature weapons, and V/R milestones.
- Approximate chart-based luck comparisons with readable explanations.
- English and Chinese UI, clearer import steps, and account confirmation.
- Manual game UID and Darkwinter-only webshop identification.
- Visible sync completion and compact target-result editing that preserves your position.

## Install

Download the **windows-x64.zip**, extract it to a writable folder, and run `ElmoArchive.exe`. Requires the Windows game client and WebView2 Runtime. Binaries are unsigned. This preview does not replace itself through the upstream updater.

Back up `ElmoBeacon.db` while the app is closed before upgrading. Do not include your database, capture files, certificates, or game logs in public reports.

## Preview status

Live capture needs region-by-region confirmation. Unknown target outcomes and missing history can affect costs; probability comparisons remain estimates. Imported history is not live inventory. See README for capture/proxy behavior and limitations.

The source ZIP corresponds to this build. SHA256SUMS.txt covers both ZIPs. This is an independently maintained fork under the original GPL v3 license, not an official game product.

Live testing on every publisher/region is not yet complete. Please include your app version and region when reporting a problem.
