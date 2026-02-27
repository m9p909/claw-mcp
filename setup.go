package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

// getInput prompts user for input and returns the trimmed response
func getInput(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// getPassword prompts user for password without echoing to terminal
func getPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // Print newline after password input
	if err != nil {
		return "", err
	}
	return string(passwordBytes), nil
}

// isValidDomain checks if domain is non-empty and has basic DNS format
func isValidDomain(domain string) bool {
	if domain == "" {
		return false
	}
	// Basic DNS format check: must contain at least one dot and valid characters
	dnsRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	return dnsRegex.MatchString(domain)
}

// promptDomain prompts user for domain with validation
func promptDomain() (string, error) {
	for {
		domain, err := getInput("Enter domain name (e.g., claw.example.com): ")
		if err != nil {
			return "", err
		}
		if isValidDomain(domain) {
			return domain, nil
		}
		fmt.Println("Invalid domain format. Please enter a valid domain name.")
	}
}

// promptPassword prompts user for password with confirmation
func promptPassword() (string, error) {
	for {
		pwd1, err := getPassword("Enter password: ")
		if err != nil {
			return "", err
		}
		if pwd1 == "" {
			fmt.Println("Password cannot be empty.")
			continue
		}

		pwd2, err := getPassword("Confirm password: ")
		if err != nil {
			return "", err
		}

		if pwd1 != pwd2 {
			fmt.Println("Passwords do not match. Please try again.")
			continue
		}

		return pwd1, nil
	}
}

// hashPassword generates bcrypt hash of password
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate password hash: %w", err)
	}
	return string(hash), nil
}

// checkSetupDone checks if setup has already been run
func checkSetupDone() (bool, error) {
	markerPath := filepath.Join("/var/lib/claw", ".setup-done")
	_, err := os.Stat(markerPath)
	if err == nil {
		return true, nil // Marker exists, setup already done
	}
	if os.IsNotExist(err) {
		return false, nil // Marker doesn't exist, setup not done
	}
	return false, err // Other error
}

// createSetupMarker creates the setup completion marker file
func createSetupMarker() error {
	markerDir := "/var/lib/claw"
	if err := os.MkdirAll(markerDir, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %w", markerDir, err)
	}

	markerPath := filepath.Join(markerDir, ".setup-done")
	f, err := os.Create(markerPath)
	if err != nil {
		return fmt.Errorf("failed to create marker file: %w", err)
	}
	f.Close()
	return nil
}

// extractCaddyfile reads the embedded Caddyfile
func extractCaddyfile() (string, error) {
	data, err := caddyfileFS.ReadFile("Caddyfile")
	if err != nil {
		return "", fmt.Errorf("failed to read embedded Caddyfile: %w", err)
	}
	return string(data), nil
}

// customizeCaddyfile replaces placeholders in Caddyfile template
func customizeCaddyfile(content string, domain string, passwordHash string) string {
	// Replace {$DOMAIN} placeholder with user domain
	content = strings.ReplaceAll(content, "{$DOMAIN:localhost}", domain)

	// Add basicauth directive for /mcp endpoint (Caddy v2.6.2 uses 'basicauth', not 'basic_auth')
	// Insert after the reverse_proxy block
	basicAuthDirective := fmt.Sprintf(`
    # Basic authentication for /mcp endpoint
    basicauth /mcp/* {
        admin %s
    }
`, passwordHash)

	// Find reverse_proxy closing brace and insert after it
	reverseProxyEnd := strings.Index(content, "    }\n\n    # TLS Configuration")
	if reverseProxyEnd != -1 {
		insertPos := reverseProxyEnd + len("    }")
		content = content[:insertPos] + basicAuthDirective + content[insertPos:]
	}

	return content
}

// writeCaddyfileConfig writes customized Caddyfile to /etc/caddy/Caddyfile
func writeCaddyfileConfig(content string) error {
	caddyfilePath := "/etc/caddy/Caddyfile"

	// Check if file already exists
	_, err := os.Stat(caddyfilePath)
	if err == nil {
		return fmt.Errorf("%s already exists; remove it manually and run setup again", caddyfilePath)
	}
	if !os.IsNotExist(err) {
		return err
	}

	// Ensure /etc/caddy directory exists
	if err := os.MkdirAll("/etc/caddy", 0755); err != nil {
		return fmt.Errorf("failed to create /etc/caddy directory: %w", err)
	}

	// Write file with appropriate permissions
	if err := os.WriteFile(caddyfilePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write Caddyfile: %w", err)
	}

	return nil
}

// getClawServiceTemplate returns the systemd service file content for Claw
func getClawServiceTemplate() string {
	return `[Unit]
Description=Claw MCP Server
After=network.target
Wants=caddy.service

[Service]
Type=simple
ExecStart=/usr/local/bin/mcpclaw -port 8080
Restart=on-failure
RestartSec=5s
StandardOutput=journal
StandardError=journal
SyslogIdentifier=claw

[Install]
WantedBy=multi-user.target
`
}

// getCaddyServiceTemplate returns the systemd service file content for Caddy
func getCaddyServiceTemplate() string {
	return `[Unit]
Description=Caddy Web Server
Documentation=https://caddyserver.com/docs/
After=network.target claw.service

[Service]
Type=notify
ExecStart=/usr/bin/caddy run --config /etc/caddy/Caddyfile
ExecReload=/usr/bin/caddy reload --config /etc/caddy/Caddyfile
TimeoutStopSec=5s
LimitNOFILE=1048576
Restart=on-failure
RestartSec=5s
StandardOutput=journal
StandardError=journal
SyslogIdentifier=caddy

[Install]
WantedBy=multi-user.target
`
}

// writeServiceFile writes a systemd service file to the specified path
func writeServiceFile(path string, content string) error {
	// Ensure /etc/systemd/system directory exists
	if err := os.MkdirAll("/etc/systemd/system", 0755); err != nil {
		return fmt.Errorf("failed to create /etc/systemd/system directory: %w", err)
	}

	// Write file with appropriate permissions
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write service file %s: %w", path, err)
	}

	return nil
}

// setupCommand runs the interactive VM setup
func setupCommand() error {
	fmt.Println("Claw VM Setup")
	fmt.Println("=============")
	fmt.Println()

	// Check if setup has already been run
	done, err := checkSetupDone()
	if err != nil {
		return fmt.Errorf("failed to check setup status: %w", err)
	}
	if done {
		return fmt.Errorf("Setup already completed on this system. To re-run setup, remove /var/lib/claw/.setup-done and try again")
	}

	// Prompt for domain
	fmt.Println("Step 1: Domain Configuration")
	domain, err := promptDomain()
	if err != nil {
		return fmt.Errorf("failed to read domain: %w", err)
	}

	// Prompt for password
	fmt.Println("\nStep 2: Set Admin Password")
	password, err := promptPassword()
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}

	// Generate password hash
	fmt.Println("\nGenerating password hash...")
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}

	// Extract and customize Caddyfile
	fmt.Println("Extracting Caddyfile template...")
	caddyfileContent, err := extractCaddyfile()
	if err != nil {
		return err
	}

	customized := customizeCaddyfile(caddyfileContent, domain, hash)

	// Write Caddyfile
	fmt.Println("Writing Caddyfile configuration...")
	if err := writeCaddyfileConfig(customized); err != nil {
		return err
	}

	// Create systemd service files
	fmt.Println("Creating systemd service files...")
	if err := writeServiceFile("/etc/systemd/system/claw.service", getClawServiceTemplate()); err != nil {
		return err
	}

	if err := writeServiceFile("/etc/systemd/system/caddy.service", getCaddyServiceTemplate()); err != nil {
		return err
	}

	// Create marker file
	fmt.Println("Finalizing setup...")
	if err := createSetupMarker(); err != nil {
		return err
	}

	fmt.Println("\n✓ Setup completed successfully!")
	fmt.Printf("\nConfiguration summary:\n")
	fmt.Printf("  Domain: %s\n", domain)
	fmt.Printf("  Caddyfile: /etc/caddy/Caddyfile\n")
	fmt.Printf("  Admin user: admin\n")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Ensure Caddy is installed: sudo apt-get install -y caddy")
	fmt.Println("  2. Run: sudo systemctl daemon-reload")
	fmt.Println("  3. Run: sudo systemctl enable --now claw caddy")
	fmt.Println("  4. Verify with: systemctl status claw caddy")
	fmt.Printf("  5. Access at: https://%s/mcp\n", domain)

	return nil
}
