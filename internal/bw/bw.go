package bw

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/comfucios/vw/internal/model"
	"github.com/comfucios/vw/internal/paths"
)

type Client struct {
	Path    string
	Session string
}

type LoginOptions struct {
	Email    string
	Password string
	Method   *int
	Code     string
}

func Resolve(configured string) (string, error) {
	if configured != "" {
		return configured, nil
	}
	if p := paths.ManagedBWPath(); isExecutable(p) {
		return p, nil
	}
	if p, err := exec.LookPath("bw"); err == nil {
		return p, nil
	}
	return "", errors.New("bw not found; run 'vw bootstrap-bw' or install the Bitwarden CLI")
}

func New(configuredPath, session string) (*Client, error) {
	p, err := Resolve(configuredPath)
	if err != nil {
		return nil, err
	}
	return &Client{Path: p, Session: session}, nil
}

func isExecutable(path string) bool {
	st, err := os.Stat(path)
	if err != nil || st.IsDir() {
		return false
	}
	return st.Mode()&0o111 != 0
}

func (c *Client) Run(args ...string) (string, error) {
	cmd := exec.Command(c.Path, args...)
	cmd.Env = os.Environ()
	if c.Session != "" {
		cmd.Env = append(cmd.Env, "BW_SESSION="+c.Session)
	}

	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(errb.String())
		if msg == "" {
			msg = strings.TrimSpace(out.String())
		}
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("bw %s: %s", strings.Join(args, " "), msg)
	}

	return strings.TrimSpace(out.String()), nil
}

func (c *Client) ConfigServer(server string) error {
	_, err := c.Run("config", "server", server)
	return err
}

// Login keeps the previous public method name for compatibility.
// It delegates fully to bw so bw can handle its native prompts.
func (c *Client) Login() error {
	return c.LoginInteractive()
}

// LoginInteractive runs `bw login` directly with stdin/stdout/stderr attached.
// This is the most compatible flow because bw handles email/password/OTP prompts.
func (c *Client) LoginInteractive() error {
	cmd := exec.Command(c.Path, "login")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}

// LoginWithPassword runs `bw login EMAIL PASSWORD` and optionally adds
// Bitwarden two-step login flags.
func (c *Client) LoginWithPassword(opts LoginOptions) error {
	email := strings.TrimSpace(opts.Email)
	code := strings.TrimSpace(opts.Code)

	if email == "" {
		return fmt.Errorf("email address is required")
	}
	if opts.Password == "" {
		return fmt.Errorf("master password is required")
	}

	args := []string{"login", email, opts.Password}

	if opts.Method != nil {
		args = append(args, "--method", strconv.Itoa(*opts.Method))
	}
	if code != "" {
		args = append(args, "--code", code)
	}

	_, err := c.Run(args...)
	return err
}

// LoginAPIKey runs `bw login --apikey` interactively.
// bw will prompt for client_id and client_secret unless BW_CLIENTID and
// BW_CLIENTSECRET are already present in the environment.
func (c *Client) LoginAPIKey() error {
	cmd := exec.Command(c.Path, "login", "--apikey")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}

// LoginSSO runs `bw login --sso` interactively.
func (c *Client) LoginSSO() error {
	cmd := exec.Command(c.Path, "login", "--sso")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}

func (c *Client) UnlockRaw() (string, error) {
	cmd := exec.Command(c.Path, "unlock", "--raw")
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()

	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(errb.String())
		if msg == "" {
			msg = strings.TrimSpace(out.String())
		}
		if msg == "" {
			msg = err.Error()
		}
		return "", errors.New(msg)
	}

	return strings.TrimSpace(out.String()), nil
}

func (c *Client) Lock() error {
	_, err := c.Run("lock")
	return err
}

func (c *Client) Sync() error {
	_, err := c.Run("sync")
	return err
}

func (c *Client) ListItems(query string) ([]model.Item, error) {
	args := []string{"list", "items"}
	if query != "" {
		args = append(args, "--search", query)
	}

	out, err := c.Run(args...)
	if err != nil {
		return nil, err
	}

	var items []model.Item
	if err := json.Unmarshal([]byte(out), &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (c *Client) FindOne(query string) (*model.Item, error) {
	items, err := c.ListItems(query)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no items matched %q", query)
	}

	if len(items) > 1 {
		return nil, fmt.Errorf("%d items matched %q; use 'vw list %s' to disambiguate", len(items), query, query)
	}

	return &items[0], nil
}

func (c *Client) GetPassword(id string) (string, error) {
	return c.Run("get", "password", id)
}

func (c *Client) GetTOTP(id string) (string, error) {
	return c.Run("get", "totp", id)
}

func (c *Client) CreateLogin(name, username, password, url string) error {
	out, err := c.Run("get", "template", "item")
	if err != nil {
		return err
	}

	var item map[string]any
	if err := json.Unmarshal([]byte(out), &item); err != nil {
		return err
	}

	item["type"] = 1
	item["name"] = name
	item["login"] = map[string]any{
		"username": username,
		"password": password,
		"uris":     []map[string]string{{"uri": url}},
	}

	b, err := json.Marshal(item)
	if err != nil {
		return err
	}

	encoded, err := c.Encode(string(b))
	if err != nil {
		return err
	}

	_, err = c.Run("create", "item", encoded)
	return err
}

func (c *Client) Encode(input string) (string, error) {
	cmd := exec.Command(c.Path, "encode")
	cmd.Stdin = strings.NewReader(input)
	cmd.Env = os.Environ()
	if c.Session != "" {
		cmd.Env = append(cmd.Env, "BW_SESSION="+c.Session)
	}

	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(errb.String())
		if msg == "" {
			msg = strings.TrimSpace(out.String())
		}
		if msg == "" {
			msg = err.Error()
		}
		return "", errors.New(msg)
	}

	return strings.TrimSpace(out.String()), nil
}

func (c *Client) PathInfo() string {
	abs, err := filepath.Abs(c.Path)
	if err != nil {
		return c.Path
	}
	return abs
}
