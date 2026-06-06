package bootstrap

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/comfucios/vw/internal/paths"
)

const releasesURL = "https://api.github.com/repos/bitwarden/clients/releases?per_page=30"

type release struct {
	TagName string  `json:"tag_name"`
	Assets  []asset `json:"assets"`
}

type asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type Options struct {
	Version string
	Force   bool
}

func InstallBW(opts Options) (string, error) {
	dest := paths.ManagedBWPath()
	if !opts.Force {
		if st, err := os.Stat(dest); err == nil && !st.IsDir() && st.Mode()&0o111 != 0 {
			return dest, nil
		}
	}
	assetName, url, err := resolveAsset(opts.Version)
	if err != nil {
		return "", err
	}
	tmp, err := os.MkdirTemp("", "vw-bw-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmp)

	archivePath := filepath.Join(tmp, assetName)
	if err := download(url, archivePath); err != nil {
		return "", err
	}
	binaryName := "bw"
	if runtime.GOOS == "windows" {
		binaryName = "bw.exe"
	}
	if err := extractZipFile(archivePath, binaryName, dest); err != nil {
		return "", err
	}
	if runtime.GOOS != "windows" {
		if err := os.Chmod(dest, 0o755); err != nil {
			return "", err
		}
	}
	return dest, nil
}

func resolveAsset(version string) (string, string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(releasesURL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", "", fmt.Errorf("GitHub API returned %s", resp.Status)
	}
	var releases []release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return "", "", err
	}
	wantTag := ""
	if version != "" {
		wantTag = "cli-v" + strings.TrimPrefix(version, "cli-v")
	}
	for _, r := range releases {
		if !strings.HasPrefix(r.TagName, "cli-v") {
			continue
		}
		if wantTag != "" && r.TagName != wantTag {
			continue
		}
		versionPart := strings.TrimPrefix(r.TagName, "cli-v")
		name := assetNameFor(versionPart)
		for _, a := range r.Assets {
			if a.Name == name {
				return a.Name, a.BrowserDownloadURL, nil
			}
		}
		return "", "", fmt.Errorf("release %s found, but asset %s was not present", r.TagName, name)
	}
	if wantTag != "" {
		return "", "", fmt.Errorf("release %s not found", wantTag)
	}
	return "", "", errors.New("no cli-v Bitwarden release found")
}

func assetNameFor(version string) string {
	switch runtime.GOOS {
	case "linux":
		return "bw-linux-" + version + ".zip"
	case "darwin":
		return "bw-macos-" + version + ".zip"
	case "windows":
		return "bw-windows-" + version + ".zip"
	default:
		return ""
	}
}

func download(url, dest string) error {
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("download failed: %s", resp.Status)
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func extractZipFile(zipPath, wanted, dest string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		if filepath.Base(f.Name) != wanted {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}
		src, err := f.Open()
		if err != nil {
			return err
		}
		defer src.Close()
		out, err := os.Create(dest)
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, src)
		return err
	}
	return fmt.Errorf("%s not found in archive", wanted)
}
