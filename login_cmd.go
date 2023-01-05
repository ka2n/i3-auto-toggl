package i3autotoggl

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/jason0x43/go-toggl"
	"github.com/manifoldco/promptui"
	"r00t2.io/gosecret"
)

const (
	SECRET_LABEL              = "i3-auto-toggl: Toggl API Token"
	SECRET_ATTR_SERVICE       = "service"
	SECRET_ATTR_SERVICE_VALUE = "service"
	SECRET_ATTR_USERNAME      = "account"
)

func LoginCmd(ctx context.Context, errOut io.Writer) error {
	toggl.DisableLog()

	svc, err := gosecret.NewService()
	if err != nil {
		return fmt.Errorf("failed to init secret service, please check if libsecret is installed: %w", err)
	}
	defer svc.Close()

	// Get collection by "default" alias
	collection, err := svc.GetCollection("default")
	if err != nil {
		return fmt.Errorf("failed to get login collection: %w", err)
	}

	if err := collection.Unlock(); err != nil {
		return fmt.Errorf("failed to unlock login collection: %w", err)
	}

	var token string
	for token == "" {
		if token, err = findStoredToken(svc); err != nil {
			return fmt.Errorf("failed to find stored token: %w", err)
		} else if token == "" {
			user, token, err := askCredentials(errOut)
			if err != nil {
				return fmt.Errorf("failed to ask credentials: %w", err)
			}
			if err := storeToken(svc, collection, user, token); err != nil {
				return fmt.Errorf("failed to store token: %w", err)
			}
		}
	}

	if err := checkToken(token, errOut); err != nil {
		return fmt.Errorf("failed to check token: %w", err)
	}

	errOut.Write([]byte("Login successful!\n"))
	return nil
}

func storeToken(svc *gosecret.Service, collection *gosecret.Collection, user string, token string) error {
	sec := gosecret.NewSecret(svc.Session, []byte{}, []byte(token), "text/plain")

	if _, err := collection.CreateItem(SECRET_LABEL, map[string]string{
		SECRET_ATTR_SERVICE:  SECRET_ATTR_SERVICE_VALUE,
		SECRET_ATTR_USERNAME: user,
	}, sec, true); err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}
	return nil
}

func findStoredToken(secret *gosecret.Service) (string, error) {
	if unlocked, _, err := secret.SearchItems(map[string]string{
		SECRET_ATTR_SERVICE: SECRET_ATTR_SERVICE_VALUE,
	}); err != nil {
		return "", fmt.Errorf("failed to search for items: %w", err)
	} else if len(unlocked) == 0 {
		return "", nil
	} else {
		return string(unlocked[0].Secret.Value), nil
	}
}

func askCredentials(errOut io.Writer) (string, string, error) {
	user_input := promptui.Prompt{
		Label: "Email",
		Validate: func(input string) error {
			if len(input) < 1 {
				return fmt.Errorf("Email is required")
			}
			return nil
		},
	}

	pass_input := promptui.Prompt{
		Label: "Password",
		Mask:  '*',
		Validate: func(input string) error {
			if len(input) < 1 {
				return fmt.Errorf("Password is required")
			}
			return nil
		},
	}
	user, err := user_input.Run()
	if err != nil {
		return "", "", err
	}

	password, err := pass_input.Run()
	if err != nil {
		return "", "", err
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(errOut))
	s.Start()
	defer s.Stop()
	sess, err := toggl.NewSession(user, password)
	if err != nil {
		return "", "", err
	}

	return user, sess.APIToken, nil
}
func checkToken(token string, errOut io.Writer) error {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(errOut))
	s.Suffix = " Checking token..."
	s.Start()
	defer s.Stop()
	sess := toggl.OpenSession(token)
	account, err := sess.GetAccount()
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}
	s.FinalMSG = fmt.Sprintf("Logged in as user id: %s\n", strconv.Itoa(account.Data.ID))

	return nil
}
