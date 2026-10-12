package service

import "ElmoBeacon/util"

type CaptureSecurityStatus struct {
	Active      bool
	Certificate *util.CaptureCertificate
	Hosts       []string
}

func GetCaptureSecurityStatus() (CaptureSecurityStatus, error) {
	captureMutex.Lock()
	defer captureMutex.Unlock()
	cert, err := util.CaptureCertificateStatus()
	return CaptureSecurityStatus{Active: mitmServer != nil || systemProxyActive, Certificate: cert, Hosts: append(append([]string{}, captureRecordHosts...), captureWebshopHost)}, err
}

func RemoveCaptureCertificate(thumbprint string) error {
	captureMutex.Lock()
	defer captureMutex.Unlock()
	if err := stopCaptureLocked(); err != nil {
		return err
	}
	return util.RemoveReviewedCaptureCertificate(thumbprint)
}
