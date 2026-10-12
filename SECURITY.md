# Security and private data

This preview captures a game-history connection through a local HTTPS proxy. Treat connection files, game logs, databases, and generated certificate keys as private.

Do not post secrets in public issues. If GitHub private vulnerability reporting is enabled, use the repository Security tab to report a vulnerability privately. Otherwise, open a public issue requesting a private contact without exploit details or private data.

Only the latest community preview receives fixes. There is no automatic binary updater. Use the repository’s Releases page and check the published SHA-256 values. Checksums detect mismatches but are not a code-signing identity.

Certificate setup requires a separate review before the native Windows confirmation. New CAs are named Elmo Archive Local Capture CA. The SHA-1 thumbprint is an identifier matching Windows, not a signing algorithm or publisher-verification claim. A root certificate has broad trust; protect its private key. Matching the fingerprint does not establish that an executable is trustworthy. Existing valid certificate files are reused when setup is declined or repeated.

## Capture scope and removal

Capture temporarily uses a local HTTPS proxy and a certificate trusted for the current Windows user. Before installation, compare the name and thumbprint shown in the app with the Windows confirmation. Windows cannot independently validate this local CA; the warning remains visible. A matching thumbprint identifies the certificate, not the publisher. The app decrypts only explicitly listed game record hosts, plus the Darkwinter webshop API in webshop mode. Other HTTPS traffic is tunneled without decryption. Upstream TLS certificates are verified. These host limits are enforced by the app; Windows root trust itself is broader. Protect ca.key.

**Remove after a successful import** is on by default. You can also open **Import records → Capture & certificate → Remove certificate trust** at any time. Removal stops capture, restores the previous proxy settings, and removes only the matching certificate from the current user’s root store. It keeps your records and local certificate files. A valid saved connection can sync without this certificate; a new capture requires installation again. If import fails or is cancelled, use the manual removal button. Other app copies or certificate stores may require separate removal through certmgr.msc. New certificates use **Elmo Archive Local Capture CA**; valid older **ElmoBeacon Custom CA** files retain their name and thumbprint.

The exact host allowlist is in `service/capture_scope.go`. Only the Darkwinter US record host has been confirmed from a local client; Haoplay entries need live regional testing. Unknown hosts are not decrypted.
