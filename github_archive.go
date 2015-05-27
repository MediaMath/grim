package grim

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

func getRepoArchive(token, location, owner, repo, ref string) (string, error) {
	downloadedName := filepath.Join(location, fmt.Sprintf("%s-%s-%s", owner, repo, ref))
	dst, err := os.Create(downloadedName)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if err := downloadRepo(token, owner, repo, ref, dst); err != nil {
		return "", err
	}

	return downloadedName, nil
}

func downloadRepo(token, owner, repo, ref string, dst io.Writer) error {
	client, err := getClientForToken(token)
	if err != nil {
		return err
	}

	u := fmt.Sprintf("repos/%s/%s/tarball/%s", owner, repo, ref)
	req, err := client.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}

	_, err = client.Do(req, dst)
	if err != nil {
		return err
	}

	return nil
}
