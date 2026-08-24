// Package sfnutil provides thin, idiomatic helpers over the AWS Step Functions
// (SFN) SDK v2: a Client wrapping *sfn.Client with convenience methods for
// provisioning state machines and activities, starting and inspecting
// executions, and the activity task-token lifecycle (poll, succeed, fail,
// heartbeat). It follows the awsgo *util convention.
package sfnutil

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
	"github.com/grokify/mogo/pointer"
)

// Client wraps an AWS Step Functions SDK v2 client.
type Client struct {
	sfnClient *sfn.Client
}

// NewClient builds a Client from an already-loaded aws.Config.
func NewClient(cfg aws.Config, optFns ...func(*sfn.Options)) *Client {
	return &Client{sfnClient: sfn.NewFromConfig(cfg, optFns...)}
}

// NewClientDefault loads the default AWS config for the given region and
// returns a Client. Pass an empty region to use the SDK's default resolution.
func NewClientDefault(ctx context.Context, region string) (*Client, error) {
	var opts []func(*awsconfig.LoadOptions) error
	if region != "" {
		opts = append(opts, awsconfig.WithRegion(region))
	}
	cfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return NewClient(cfg), nil
}

// SFN returns the underlying SDK client for calls not covered by a helper.
func (c *Client) SFN() *sfn.Client { return c.sfnClient }

// CreateStandardStateMachine creates a STANDARD-type state machine and returns
// its ARN. STANDARD (not EXPRESS) is required for long-running, human-in-the-
// loop workflows.
func (c *Client) CreateStandardStateMachine(ctx context.Context, name, definitionJSON, roleARN string) (arn string, err error) {
	out, err := c.sfnClient.CreateStateMachine(ctx, &sfn.CreateStateMachineInput{
		Name:       pointer.Pointer(name),
		Definition: pointer.Pointer(definitionJSON),
		RoleArn:    pointer.Pointer(roleARN),
		Type:       types.StateMachineTypeStandard,
	})
	if err != nil {
		return "", err
	}
	return pointer.Dereference(out.StateMachineArn), nil
}

// UpdateStateMachineDefinition updates an existing state machine's definition
// and role.
func (c *Client) UpdateStateMachineDefinition(ctx context.Context, stateMachineARN, definitionJSON, roleARN string) error {
	_, err := c.sfnClient.UpdateStateMachine(ctx, &sfn.UpdateStateMachineInput{
		StateMachineArn: pointer.Pointer(stateMachineARN),
		Definition:      pointer.Pointer(definitionJSON),
		RoleArn:         pointer.Pointer(roleARN),
	})
	return err
}

// FindStateMachineARN returns the ARN of the state machine with the given name,
// or ok=false if none exists. It pages through ListStateMachines.
func (c *Client) FindStateMachineARN(ctx context.Context, name string) (arn string, ok bool, err error) {
	p := sfn.NewListStateMachinesPaginator(c.sfnClient, &sfn.ListStateMachinesInput{})
	for p.HasMorePages() {
		page, perr := p.NextPage(ctx)
		if perr != nil {
			return "", false, perr
		}
		for _, sm := range page.StateMachines {
			if pointer.Dereference(sm.Name) == name {
				return pointer.Dereference(sm.StateMachineArn), true, nil
			}
		}
	}
	return "", false, nil
}

// CreateActivity creates an activity (idempotent in SFN: creating one that
// already exists returns the existing ARN) and returns its ARN.
func (c *Client) CreateActivity(ctx context.Context, name string) (arn string, err error) {
	out, err := c.sfnClient.CreateActivity(ctx, &sfn.CreateActivityInput{
		Name: pointer.Pointer(name),
	})
	if err != nil {
		return "", err
	}
	return pointer.Dereference(out.ActivityArn), nil
}

// StartExecution starts a named execution and returns its ARN. The name is the
// idempotency key: SFN rejects a duplicate name with different input and is a
// no-op for an identical re-start, which is how run idempotency is enforced.
func (c *Client) StartExecution(ctx context.Context, stateMachineARN, name, inputJSON string) (executionARN string, err error) {
	in := &sfn.StartExecutionInput{
		StateMachineArn: pointer.Pointer(stateMachineARN),
	}
	if name != "" {
		in.Name = pointer.Pointer(name)
	}
	if inputJSON != "" {
		in.Input = pointer.Pointer(inputJSON)
	}
	out, err := c.sfnClient.StartExecution(ctx, in)
	if err != nil {
		return "", err
	}
	return pointer.Dereference(out.ExecutionArn), nil
}

// Execution is a normalized view of DescribeExecution.
type Execution struct {
	ARN    string
	Status types.ExecutionStatus
	Input  string
	Output string // empty unless the execution succeeded
}

// DescribeExecution returns the current state of an execution.
func (c *Client) DescribeExecution(ctx context.Context, executionARN string) (*Execution, error) {
	out, err := c.sfnClient.DescribeExecution(ctx, &sfn.DescribeExecutionInput{
		ExecutionArn: pointer.Pointer(executionARN),
	})
	if err != nil {
		return nil, err
	}
	return &Execution{
		ARN:    pointer.Dereference(out.ExecutionArn),
		Status: out.Status,
		Input:  pointer.Dereference(out.Input),
		Output: pointer.Dereference(out.Output),
	}, nil
}

// StopExecution aborts a running execution.
func (c *Client) StopExecution(ctx context.Context, executionARN, errCode, cause string) error {
	in := &sfn.StopExecutionInput{ExecutionArn: pointer.Pointer(executionARN)}
	if errCode != "" {
		in.Error = pointer.Pointer(errCode)
	}
	if cause != "" {
		in.Cause = pointer.Pointer(cause)
	}
	_, err := c.sfnClient.StopExecution(ctx, in)
	return err
}

// ActivityTask is a task polled from an activity queue. An empty TaskToken
// means the long poll timed out with no work — the caller should poll again.
type ActivityTask struct {
	TaskToken string
	Input     string
}

// PollActivityTask long-polls for one task on the given activity. workerName is
// an optional identifier for observability. A returned task with an empty
// TaskToken indicates the poll timed out with no work available.
func (c *Client) PollActivityTask(ctx context.Context, activityARN, workerName string) (*ActivityTask, error) {
	in := &sfn.GetActivityTaskInput{ActivityArn: pointer.Pointer(activityARN)}
	if workerName != "" {
		in.WorkerName = pointer.Pointer(workerName)
	}
	out, err := c.sfnClient.GetActivityTask(ctx, in)
	if err != nil {
		return nil, err
	}
	return &ActivityTask{
		TaskToken: pointer.Dereference(out.TaskToken),
		Input:     pointer.Dereference(out.Input),
	}, nil
}

// SendTaskSuccess completes a task token with the given JSON output.
func (c *Client) SendTaskSuccess(ctx context.Context, taskToken, outputJSON string) error {
	if outputJSON == "" {
		outputJSON = "{}"
	}
	_, err := c.sfnClient.SendTaskSuccess(ctx, &sfn.SendTaskSuccessInput{
		TaskToken: pointer.Pointer(taskToken),
		Output:    pointer.Pointer(outputJSON),
	})
	return err
}

// SendTaskFailure fails a task token with an error code and cause.
func (c *Client) SendTaskFailure(ctx context.Context, taskToken, errCode, cause string) error {
	in := &sfn.SendTaskFailureInput{TaskToken: pointer.Pointer(taskToken)}
	if errCode != "" {
		in.Error = pointer.Pointer(errCode)
	}
	if cause != "" {
		in.Cause = pointer.Pointer(cause)
	}
	_, err := c.sfnClient.SendTaskFailure(ctx, in)
	return err
}

// SendTaskHeartbeat reports liveness for a long-running task token, resetting
// its heartbeat timeout.
func (c *Client) SendTaskHeartbeat(ctx context.Context, taskToken string) error {
	_, err := c.sfnClient.SendTaskHeartbeat(ctx, &sfn.SendTaskHeartbeatInput{
		TaskToken: pointer.Pointer(taskToken),
	})
	return err
}

// ErrNoStateMachineName is returned by helpers requiring a non-empty name.
var ErrNoStateMachineName = errors.New("sfnutil: state machine name is required")
