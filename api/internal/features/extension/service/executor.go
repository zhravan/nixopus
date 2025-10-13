package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/raghavyuva/nixopus-api/internal/features/logger"
	"github.com/raghavyuva/nixopus-api/internal/features/ssh"
	"github.com/raghavyuva/nixopus-api/internal/types"
)

func (s *ExtensionService) executeRun(exec *types.ExtensionExecution, spec types.ExtensionSpec, vars map[string]interface{}) {
	exec.StartedAt = time.Now()
	_ = s.storage.UpdateExecution(exec)

	s.logger.Log(logger.Info, fmt.Sprintf("Starting extension execution: %s", exec.ID.String()), "")

	replacer := buildReplacer(vars)
	sshClient := ssh.NewSSH()
	steps, _ := s.storage.ListExecutionSteps(exec.ID.String())

	if stop := s.processPhase(exec, &steps, "run", spec.Execution.Run, 0, sshClient, replacer); stop {
		s.logger.Log(logger.Error, fmt.Sprintf("Extension execution failed during run phase: %s", exec.ID.String()), "")
		return
	}
	if stop := s.processPhase(exec, &steps, "validate", spec.Execution.Validate, len(spec.Execution.Run), sshClient, replacer); stop {
		s.logger.Log(logger.Error, fmt.Sprintf("Extension execution failed during validate phase: %s", exec.ID.String()), "")
		return
	}

	exec.Status = types.ExecutionStatusCompleted
	finished := time.Now()
	exec.CompletedAt = &finished
	_ = s.storage.UpdateExecution(exec)

	s.logger.Log(logger.Info, fmt.Sprintf("Extension execution completed successfully: %s", exec.ID.String()), "")
}

func buildReplacer(vars map[string]interface{}) func(string) string {
	return func(in string) string {
		out := in
		for k, v := range vars {
			token := fmt.Sprintf("{{ %s }}", k)
			out = strings.ReplaceAll(out, token, fmt.Sprint(v))
		}
		return out
	}
}

func (s *ExtensionService) processPhase(
	exec *types.ExtensionExecution,
	steps *[]types.ExecutionStep,
	phase string,
	specSteps []types.SpecStep,
	offset int,
	sshClient *ssh.SSH,
	replacer func(string) string,
) bool {
	s.logger.Log(logger.Info, fmt.Sprintf("Processing phase: %s with %d spec steps, offset: %d", phase, len(specSteps), offset), "")
	s.logger.Log(logger.Info, fmt.Sprintf("Available steps: %d", len(*steps)), "")

	for idx := range specSteps {
		if s.shouldCancel(exec) {
			s.markCancelled(exec)
			return true
		}

		stepOrder := offset + idx + 1
		step := s.getStepByPhaseAndOrder(steps, phase, stepOrder)
		if step == nil {
			s.logger.Log(logger.Error, fmt.Sprintf("Step not found for phase: %s, order: %d", phase, stepOrder), "")
			continue
		}
		s.beginStep(step)

		out, runErr := s.executeSpecStep(specSteps[idx], sshClient, replacer)
		step.Output = out
		if runErr != nil {
			if s.failStep(exec, step, out, runErr, specSteps[idx].IgnoreErrors) {
				return true
			}
			continue
		}

		s.completeStep(exec, step, out)
	}
	return false
}

func (s *ExtensionService) shouldCancel(exec *types.ExtensionExecution) bool {
	cur, err := s.storage.GetExecutionByID(exec.ID.String())
	return err == nil && cur.Status == types.ExecutionStatusCancelled
}

func (s *ExtensionService) markCancelled(exec *types.ExtensionExecution) {
	exec.Status = types.ExecutionStatusCancelled
	finished := time.Now()
	exec.CompletedAt = &finished
	_ = s.storage.UpdateExecution(exec)
}

func (s *ExtensionService) beginStep(step *types.ExecutionStep) {
	step.Status = types.ExecutionStatusRunning
	step.StartedAt = time.Now()

	// Log to console
	s.logger.Log(logger.Info, fmt.Sprintf("STEP STARTED: %s", step.StepName), "")

	_ = s.storage.UpdateExecutionStep(step)
}

func (s *ExtensionService) completeStep(exec *types.ExtensionExecution, step *types.ExecutionStep, out string) {
	step.Status = types.ExecutionStatusCompleted
	completed := time.Now()
	step.CompletedAt = &completed

	// Enhanced logging with timestamp and step details
	successLog := fmt.Sprintf("[%s] STEP COMPLETED: %s\nOutput: %s",
		completed.Format("2006-01-02 15:04:05"),
		step.StepName,
		out)

	// Log to console
	s.logger.Log(logger.Info, fmt.Sprintf("STEP COMPLETED: %s - %s", step.StepName, out), "")

	exec.ExecutionLog = exec.ExecutionLog + "\n" + successLog
	_ = s.storage.UpdateExecutionStep(step)
	_ = s.storage.UpdateExecution(exec)
}

func (s *ExtensionService) failStep(exec *types.ExtensionExecution, step *types.ExecutionStep, out string, runErr error, ignore bool) bool {
	step.Status = types.ExecutionStatusFailed
	completed := time.Now()
	step.CompletedAt = &completed

	// Enhanced logging with timestamp and step details
	errorLog := fmt.Sprintf("[%s] STEP FAILED: %s\nError: %v\nOutput: %s",
		completed.Format("2006-01-02 15:04:05"),
		step.StepName,
		runErr,
		out)

	// Log to console
	s.logger.Log(logger.Error, fmt.Sprintf("STEP FAILED: %s - Error: %v, Output: %s", step.StepName, runErr, out), "")

	exec.ExecutionLog = exec.ExecutionLog + "\n" + errorLog
	_ = s.storage.UpdateExecutionStep(step)
	_ = s.storage.UpdateExecution(exec)

	if ignore {
		return false
	}

	exec.Status = types.ExecutionStatusFailed
	exec.ErrorMessage = runErr.Error()
	finished := time.Now()
	exec.CompletedAt = &finished
	_ = s.storage.UpdateExecution(exec)
	return true
}

func (s *ExtensionService) getStepByPhaseAndOrder(steps *[]types.ExecutionStep, phase string, order int) *types.ExecutionStep {
	s.logger.Log(logger.Info, fmt.Sprintf("Looking for step: phase=%s, order=%d", phase, order), "")
	for i := range *steps {
		st := &(*steps)[i]
		s.logger.Log(logger.Info, fmt.Sprintf("Checking step: phase=%s, order=%d, name=%s", st.Phase, st.StepOrder, st.StepName), "")
		if st.Phase == phase && st.StepOrder == order {
			s.logger.Log(logger.Info, fmt.Sprintf("Found step: %s", st.StepName), "")
			return st
		}
	}
	s.logger.Log(logger.Error, fmt.Sprintf("Step not found: phase=%s, order=%d", phase, order), "")
	return nil
}

func (s *ExtensionService) executeSpecStep(spec types.SpecStep, sshClient *ssh.SSH, replacer func(string) string) (string, error) {
	switch spec.Type {
	case "command":
		cmd, _ := spec.Properties["cmd"].(string)
		if cmd == "" {
			return "", fmt.Errorf("command is required for command step")
		}
		cmd = replacer(cmd)
		prefix := s.timeoutPrefix(sshClient, spec.Timeout)
		output, err := sshClient.RunCommand(prefix + cmd)
		if err != nil {
			return "", fmt.Errorf("failed to execute command '%s': %w (output: %s)", cmd, err, output)
		}
		return output, nil
	case "file":
		output, err := executeFileStep(sshClient, spec.Properties, replacer)
		if err != nil {
			return "", fmt.Errorf("file operation failed: %w", err)
		}
		return output, nil
	case "service":
		name, _ := spec.Properties["name"].(string)
		action, _ := spec.Properties["action"].(string)
		if name == "" {
			return "", fmt.Errorf("service name is required for service step")
		}
		if action == "" {
			return "", fmt.Errorf("service action is required for service step")
		}
		name = replacer(name)
		// Prefer systemctl; if not available, fallback to service
		cmd := s.serviceCmd(sshClient, action, name, spec.Timeout)
		output, err := sshClient.RunCommand(cmd)
		if err != nil {
			return "", fmt.Errorf("failed to execute service command '%s': %w (output: %s)", cmd, err, output)
		}
		return output, nil
	case "user":
		output, err := s.executeUserStep(sshClient, spec.Properties, replacer, spec.Timeout)
		if err != nil {
			return "", fmt.Errorf("user operation failed: %w", err)
		}
		return output, nil
	default:
		return "", fmt.Errorf("unsupported step type: %s", spec.Type)
	}
}
