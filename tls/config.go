// Kiebitz - Privacy-Friendly Appointment Scheduling
// Copyright (C) 2021-2021 The Kiebitz Authors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version. Additional terms
// as defined in section 7 of the license (e.g. regarding attribution)
// are specified at https://kiebitz.eu/en/docs/open-source/additional-terms.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/kiebitz-oss/services"
	"io/ioutil"
)

func TLSConfig(settings *services.TLSSettings) (*tls.Config, error) {

	bs, err := ioutil.ReadFile(settings.CACertificateFile)

	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()

	if ok := certPool.AppendCertsFromPEM(bs); !ok {
		return nil, fmt.Errorf("cannot import CA certificate")
	}

	cert, err := tls.LoadX509KeyPair(settings.CertificateFile, settings.KeyFile)

	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
		PreferServerCipherSuites: true,
		Certificates:             []tls.Certificate{cert},
		ClientCAs:                certPool,
		RootCAs:                  certPool,
	}

	return tlsConfig, nil
}

func TLSClientConfig(settings *services.TLSSettings, serverName string) (*tls.Config, error) {

	if config, err := TLSConfig(settings); err != nil {
		return nil, err
	} else {
		config.ServerName = serverName
		return config, nil
	}
}

func TLSServerConfig(settings *services.TLSSettings) (*tls.Config, error) {

	if config, err := TLSConfig(settings); err != nil {
		return nil, err
	} else {
		config.ClientAuth = tls.RequireAndVerifyClientCert
		return config, nil
	}
}
