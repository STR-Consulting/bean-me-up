package beans

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Client executes beans CLI commands and parses output.
type Client struct {
	beansPath string // Path to the beans directory
}

// NewClient creates a new beans CLI client.
func NewClient(beansPath string) *Client {
	return &Client{
		beansPath: beansPath,
	}
}

// List returns all beans from the beans CLI.
func (c *Client) List() ([]Bean, error) {
	args := []string{"list", "--json", "--full"}
	if c.beansPath != "" {
		args = append(args, "--beans-path", c.beansPath)
	}

	out, err := c.exec(args...)
	if err != nil {
		return nil, err
	}

	var beans []Bean
	if err := json.Unmarshal(out, &beans); err != nil {
		return nil, fmt.Errorf("parsing beans JSON: %w", err)
	}

	return beans, nil
}

// Get returns a specific bean by ID.
func (c *Client) Get(id string) (*Bean, error) {
	args := []string{"show", "--json", id}
	if c.beansPath != "" {
		args = append(args, "--beans-path", c.beansPath)
	}

	out, err := c.exec(args...)
	if err != nil {
		return nil, err
	}

	// beans show --json with a single ID returns a single object (not an array)
	var bean Bean
	if err := json.Unmarshal(out, &bean); err != nil {
		return nil, fmt.Errorf("parsing bean JSON: %w", err)
	}

	if bean.ID == "" {
		return nil, fmt.Errorf("bean not found: %s", id)
	}

	return &bean, nil
}

// GetMultiple returns multiple beans by ID.
func (c *Client) GetMultiple(ids []string) ([]Bean, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	// For a single ID, use Get (beans show returns object, not array)
	if len(ids) == 1 {
		bean, err := c.Get(ids[0])
		if err != nil {
			return nil, err
		}
		return []Bean{*bean}, nil
	}

	args := []string{"show", "--json"}
	args = append(args, ids...)
	if c.beansPath != "" {
		args = append(args, "--beans-path", c.beansPath)
	}

	out, err := c.exec(args...)
	if err != nil {
		return nil, err
	}

	var beans []Bean
	if err := json.Unmarshal(out, &beans); err != nil {
		return nil, fmt.Errorf("parsing beans JSON: %w", err)
	}

	return beans, nil
}

// GraphQL executes a GraphQL query via the beans CLI, passing the query via stdin.
func (c *Client) GraphQL(query string) ([]byte, error) {
	args := []string{"query", "--json"}
	if c.beansPath != "" {
		args = append(args, "--beans-path", c.beansPath)
	}

	cmd := exec.Command("beans", args...)
	cmd.Stdin = strings.NewReader(query)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("beans query: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("beans query: %w", err)
	}
	return out, nil
}

// ExternalDataOp describes a single setExternalData operation for batching.
type ExternalDataOp struct {
	BeanID string
	Plugin string
	Data   map[string]any
}

// SetExternalData sets external plugin data on a single bean.
func (c *Client) SetExternalData(id, plugin string, data map[string]any) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshaling external data: %w", err)
	}

	query := fmt.Sprintf(
		`mutation { setExternalData(id: %q, plugin: %q, data: %s) { id } }`,
		id, plugin, string(dataJSON),
	)

	_, err = c.GraphQL(query)
	return err
}

// RemoveExternalData removes external plugin data from a single bean.
func (c *Client) RemoveExternalData(id, plugin string) error {
	query := fmt.Sprintf(
		`mutation { removeExternalData(id: %q, plugin: %q) { id } }`,
		id, plugin,
	)

	_, err := c.GraphQL(query)
	return err
}

// SetExternalDataBatch sets external plugin data on multiple beans in a single
// GraphQL call using aliased mutations.
func (c *Client) SetExternalDataBatch(ops []ExternalDataOp) error {
	if len(ops) == 0 {
		return nil
	}

	var b strings.Builder
	b.WriteString("mutation {\n")
	for i, op := range ops {
		dataJSON, err := json.Marshal(op.Data)
		if err != nil {
			return fmt.Errorf("marshaling external data for %s: %w", op.BeanID, err)
		}
		fmt.Fprintf(&b, "  op%d: setExternalData(id: %q, plugin: %q, data: %s) { id }\n",
			i, op.BeanID, op.Plugin, string(dataJSON))
	}
	b.WriteString("}\n")

	_, err := c.GraphQL(b.String())
	return err
}

// exec runs a beans command and returns the output.
func (c *Client) exec(args ...string) ([]byte, error) {
	cmd := exec.Command("beans", args...)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("beans %s: %s", strings.Join(args, " "), string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("beans %s: %w", strings.Join(args, " "), err)
	}
	return out, nil
}
