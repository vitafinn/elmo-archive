package handler

import (
	"context"
	"sync"

	"ElmoBeacon/service"
	"ElmoBeacon/util"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type App struct {
	ctx           context.Context
	importMu      sync.Mutex
	pendingImport *preparedImport
	importBusy    bool
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) SyncRecordsOS() (*service.SyncResult, error) {
	gameUserInfo, err := service.GetUserInfoOS()
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, errors.Errorf("failed to extract game user information")
	}

	syncResult, err := service.SyncRecords(gameUserInfo)
	if err != nil {
		return nil, err
	}

	return syncResult, nil
}

func (a *App) StartCaptureOS(uidMode string) error {
	return service.StartCapture(a.ctx, uidMode)
}

func (a *App) StopCaptureOS() error {
	return service.StopCapture()
}

func (a *App) PrepareCaptureCertificate() (util.CaptureCertificate, error) {
	return util.PrepareCaptureCertificate()
}

func (a *App) InstallCaptureCertificate(thumbprint string) error {
	return util.InstallReviewedCertificate(thumbprint)
}

func (a *App) FocusCertificatePrompt() bool { return util.FocusCertificatePrompt() }

func (a *App) GetCaptureSecurityStatus() (service.CaptureSecurityStatus, error) {
	return service.GetCaptureSecurityStatus()
}

func (a *App) RemoveCaptureCertificate(thumbprint string) error {
	return service.RemoveCaptureCertificate(thumbprint)
}
